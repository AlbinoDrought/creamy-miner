package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"go.snowblossom/internal/sbconst"
	"go.snowblossom/internal/sblib"
	"go.snowblossom/internal/sbproto/mining_pool"
	"go.snowblossom/internal/sbproto/snowblossom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	fields              map[int]sblib.SnowMerkleProof
	client              mining_pool.MiningPoolServiceClient
	currentWorkID       int32
	currentWorkUnit     *mining_pool.WorkUnit
	currentWorkUnitLock sync.RWMutex

	sharesPerformedLastDrain time.Time
	sharesPerformed          int64

	baseLogger *logrus.Logger
)

func main() {
	var err error

	pool := os.Getenv("CREAMY_MINER_POOL")
	if pool == "" {
		pool = "localhost:23380"
	}

	address := os.Getenv("CREAMY_MINER_ADDRESS")
	if address == "" {
		address = "c04rt84spfjc9xy88snx5r256qv0tmy664zcdrnc"
	}

	threadMultiplierString := os.Getenv("CREAMY_MINER_THREAD_MULTIPLIER")
	if threadMultiplierString == "" {
		threadMultiplierString = "10"
	}
	threadMultiplier, err := strconv.Atoi(threadMultiplierString)
	if err != nil {
		panic(err)
	}

	numThreads := runtime.NumCPU() * threadMultiplier

	baseLogger = logrus.New()
	baseLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	baseLogger.
		WithField("pool", pool).
		WithField("address", address).
		WithField("thread-multiplier", threadMultiplier).
		WithField("total-num-threads", numThreads).
		Info("loaded config")

	fields, err = sblib.LoadFields("fields")
	if err != nil {
		panic(err)
	}
	baseLogger.
		WithField("field-count", len(fields)).
		Info("loaded fields")

	// conn, err := grpc.Dial("snowypool.com:23380", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(pool, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = mining_pool.NewMiningPoolServiceClient(conn)

	for i := 0; i < numThreads; i++ {
		go func(i int) {
			minerThread(i)
		}(i)
	}

	sharesPerformedLastDrain = time.Now()
	go func() {
		log := baseLogger.WithField("thread", "benchmark")
		for {
			time.Sleep(15 * time.Second)
			secondsSince := time.Since(sharesPerformedLastDrain)
			sharesPerformedLastDrain = time.Now()
			sharesPerformedCopy := atomic.AddInt64(&sharesPerformed, 0)
			atomic.AddInt64(&sharesPerformed, -sharesPerformedCopy)

			log.WithField("H/s", fmt.Sprintf("%.4f", float64(sharesPerformedCopy)/secondsSince.Seconds())).Printf("15 second average")
		}
	}()

	go func() {
		log := baseLogger.WithField("thread", "start-watcher")
		for {
			sharesPerformedCopy := atomic.AddInt64(&sharesPerformed, 0)
			if sharesPerformedCopy == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			log.Info("started mining!")
			break
		}
	}()

	log := baseLogger.WithField("thread", "main")
	for {
		log.Printf("opening GetWork connection")
		stream, err := client.GetWork(context.Background(), &mining_pool.GetWorkRequest{
			PayToAddress: address,
		})
		if err != nil {
			log.WithError(err).Warn("failed to open GetWork stream!")
			time.Sleep(10 * time.Second)
			continue // retry connection
		}

		lastBlockHeight := int32(0)
		for {
			log.Debug("waiting for work")
			work, err := stream.Recv()
			if err != nil {
				log.WithError(err).Warn("failed to receive work!")
				time.Sleep(10 * time.Second)
				break // retry connection
			}

			if work.Header.BlockHeight != lastBlockHeight {
				log.WithField("block-height", work.Header.BlockHeight).Info("workunit for new blockheight")
				lastBlockHeight = work.Header.BlockHeight
			}

			log.Debug("received work: %+v", work)
			currentWorkUnitLock.Lock()
			currentWorkUnit = work
			currentWorkID = work.WorkId
			currentWorkUnitLock.Unlock()
			log.Debug("assigned work")
		}
	}
}

func minerThread(threadID int) {
	// once per thread
	log := baseLogger.WithField("thread", threadID)

	md := sblib.NewMessageDigest()

	wordBuffer := make([]byte, sbconst.SnowMerkleHashLen)
	wordIndexHistory := make([]int64, sbconst.POWLookPasses)

	nonce := make([]byte, sbconst.NonceLength)

	log.Debug("thread booted")

	// once per unique work unit
	for {
		// copy header
		currentWorkUnitLock.RLock()
		if currentWorkUnit == nil {
			currentWorkUnitLock.RUnlock()
			log.Debug("no work unit found on thread, sleeping")
			time.Sleep(time.Second)
			continue
		}

		workID := currentWorkUnit.WorkId
		reportTarget := currentWorkUnit.ReportTarget
		workHeader := snowblossom.BlockHeader{
			Version:             currentWorkUnit.Header.Version,
			BlockHeight:         currentWorkUnit.Header.BlockHeight,
			PrevBlockHash:       currentWorkUnit.Header.PrevBlockHash,
			MerkleRootHash:      currentWorkUnit.Header.MerkleRootHash,
			UtxoRootHash:        currentWorkUnit.Header.UtxoRootHash,
			Nonce:               currentWorkUnit.Header.Nonce,
			Timestamp:           currentWorkUnit.Header.Timestamp,
			Target:              currentWorkUnit.Header.Target,
			SnowField:           currentWorkUnit.Header.SnowField,
			SnowHash:            currentWorkUnit.Header.SnowHash,
			PowProof:            currentWorkUnit.Header.PowProof,
			ShardId:             currentWorkUnit.Header.ShardId,
			ShardExportRootHash: currentWorkUnit.Header.ShardExportRootHash,
			ShardImport:         currentWorkUnit.Header.ShardImport,
			TxDataSizeSum:       currentWorkUnit.Header.TxDataSizeSum,
			TxCount:             currentWorkUnit.Header.TxCount,
		}
		currentWorkUnitLock.RUnlock()
		merkleProof := fields[int(workHeader.SnowField)]
		// workingSince := time.Now()
		copy(nonce, workHeader.Nonce[:4])
		rand.Read(nonce[4:])
		log.Debug("activated work unit #%v", uint64(workID))

		// repeatable
		quickies := 0
		for {
			quickies++
			if quickies > 10 { // only check every ten loops
				quickies = 0
				// if time.Since(workingSince) > 75*time.Second {
				activeWorkID := atomic.AddInt32(&currentWorkID, 0)
				if workID != activeWorkID {
					break // update work unit
				}
				// }
			}

			// if threadID%2 == 0 {
			// read from rand
			if _, err := rand.Read(nonce[4:]); err != nil {
				panic(err)
			}
			/*
					} else {
					// just add 1
					binary.BigEndian.PutUint64(
						nonce[4:],
						binary.BigEndian.Uint64(nonce[4:])+1,
					)
				}
			*/

			snowContext, err := sblib.HashHeaderBits(&workHeader, nonce, md) // todo: this could be generated without I/O waits
			if err != nil {
				panic(err)
			}

			for pass := 0; pass < sbconst.POWLookPasses; pass++ {
				wordIndex := sblib.GetNextSnowFieldIndex(snowContext, merkleProof.TotalWords(), md)
				wordIndexHistory[pass] = wordIndex
				err = merkleProof.ReadWord(wordIndex, wordBuffer, pass)
				if err != nil {
					panic(err)
				}
				snowContext = sblib.GetNextContext(snowContext, wordBuffer, md)
			}

			atomic.AddInt64(&sharesPerformed, 1)
			if sblib.LessThanTarget(snowContext, reportTarget) {
				log.Printf("Found passable solution: %X", snowContext)

				workHeader.Nonce = nonce

				powProofs := make([]*snowblossom.SnowPowProof, len(wordIndexHistory))
				for pass, wordIndex := range wordIndexHistory {
					powProofs[pass], err = merkleProof.GetProof(wordIndex)
					if err != nil {
						panic(err)
					}
				}
				workHeader.PowProof = powProofs
				workHeader.SnowHash = snowContext

				// todo: client threadsafe?
				resp, err := client.SubmitWork(context.Background(), &mining_pool.WorkSubmitRequest{
					WorkId: workID,
					Header: &workHeader,
				})
				log.Printf("%+v", resp)
				log.Printf("%+v", err)
			} /* else {
				log.Println("Solution not passable")
			} */
		}
	}
}
