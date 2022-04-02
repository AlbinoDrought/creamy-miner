package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"go.snowblossom/internal/extra"
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
	currentPool         string
	currentWorkID       int32
	currentWorkUnit     *mining_pool.WorkUnit
	currentWorkUnitLock sync.RWMutex

	sharesPerformedLastDrain time.Time
	sharesPerformed          int64

	baseLogger *logrus.Logger

	poolStates map[string]*poolState
)

type poolState struct {
	shareSubmitter chan *mining_pool.WorkSubmitRequest
	shares         uint64
	client         mining_pool.MiningPoolServiceClient
	lock           sync.RWMutex
}

func (ps *poolState) Submit(work *mining_pool.WorkSubmitRequest) (*snowblossom.SubmitReply, error) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	attempts := 0
	for ps.client == nil {
		attempts++
		if attempts > 10 {
			return nil, errors.New("no client available after 10s")
		}
		ps.lock.RUnlock()
		time.Sleep(5 * time.Second)
		ps.lock.RLock()
	}

	reply, err := ps.client.SubmitWork(context.Background(), work)
	if err != nil {
		return nil, err
	}

	atomic.AddUint64(&ps.shares, 1)
	return reply, nil
}

type poolWork struct {
	pool string
	work *mining_pool.WorkUnit
}

func main() {
	var err error

	baseLogger = logrus.New()
	baseLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	envFilePath := os.Getenv("CREAMY_MINER_ENV_FILE")
	if envFilePath != "" {
		log := baseLogger.WithField("env-file-path", envFilePath)

		envFile, err := ioutil.ReadFile(envFilePath)
		if err != nil {
			log.WithError(err).Fatal("failed opening env file, does it exist?")
		}

		loaded := 0
		for line, envFileLine := range strings.Split(string(envFile), "\n") {
			if envFileLine == "" || envFileLine[0] == '#' {
				continue
			}
			envFilePair := strings.SplitN(envFileLine, "=", 2)
			if len(envFilePair) != 2 {
				log.WithField("line", line).WithField("contents", envFileLine).Fatal("line must be in format of KEY=VALUE")
			}

			if err := os.Setenv(envFilePair[0], envFilePair[1]); err != nil {
				log.WithField("line", line).WithField("contents", envFileLine).WithError(err).Fatal("failed to set env")
			}
			loaded++
		}
		log.WithField("loaded-vars", loaded).Info("loaded env file!")
	}

	poolsStr := os.Getenv("CREAMY_MINER_POOLS")
	if poolsStr == "" {
		poolsStr = os.Getenv("CREAMY_MINER_POOL")
		if poolsStr == "" {
			poolsStr = "localhost:23380"
		}
	}
	pools := strings.Split(poolsStr, ",")

	poolStates = map[string]*poolState{}
	for _, pool := range pools {
		poolStates[pool] = &poolState{
			shareSubmitter: make(chan *mining_pool.WorkSubmitRequest),
		}
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

	fieldsStr := os.Getenv("CREAMY_MINER_FIELDS_DIRS")
	if fieldsStr == "" {
		fieldsStr = os.Getenv("CREAMY_MINER_FIELDS_DIR")
		if fieldsStr == "" {
			fieldsStr = "fields"
		}
	}

	fieldsDirs := strings.Split(fieldsStr, ",")

	numThreads := runtime.NumCPU() * threadMultiplier

	baseLogger.
		WithField("pools", pools).
		WithField("address", address).
		WithField("thread-multiplier", threadMultiplier).
		WithField("total-num-threads", numThreads).
		WithField("fields-dirs", fieldsDirs).
		Info("loaded config")

	totalFields := 0
	allFields := make([]map[int]sblib.SnowMerkleProof, len(fieldsDirs))
	for i, fieldsDir := range fieldsDirs {
		allFields[i], err = sblib.LoadFields(fieldsDir)
		if err != nil {
			baseLogger.WithField("fields-dir", fieldsDir).WithError(err).Fatal("failed to load fields dir")
		}
		baseLogger.
			WithField("fields-dir", fieldsDir).
			WithField("field-count", len(allFields[i])).
			Info("loaded fields")
		totalFields += len(allFields[i])
	}

	fields = extra.MergeFields(allFields)

	baseLogger.
		WithField("field-count-pre-merge", totalFields).
		WithField("field-count-post-merge", len(fields)).
		Info("finished loading and merging all fields")

	if len(fields) == 0 {
		baseLogger.Fatal("no fields loaded")
	}

	workUnitChan := make(chan poolWork, len(poolStates)*5)

	// pool recvWork threads
	for pool := range poolStates {
		go func(pool string) {
			log := baseLogger.
				WithField("thread", "pool-conn").
				WithField("pool", pool)

			sleepDuration := time.Second

			for {
				conn, err := grpc.Dial(pool, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.WithError(err).
						WithField("sleep-duration", sleepDuration).
						Warn("failed to dial pool!")
					time.Sleep(sleepDuration)
					if sleepDuration < 30*time.Minute {
						sleepDuration *= 2
					}
					continue
				}
				sleepDuration = time.Second

				client := mining_pool.NewMiningPoolServiceClient(conn)

				poolStates[pool].lock.Lock()
				poolStates[pool].client = client
				poolStates[pool].lock.Unlock()

				log.Info("opening GetWork connection")
				stream, err := client.GetWork(context.Background(), &mining_pool.GetWorkRequest{
					PayToAddress: address,
				})
				if err != nil {
					log.WithError(err).Warn("failed to open GetWork stream!")
					time.Sleep(10 * time.Second)
					continue // retry connection
				}

				for {
					log.Debug("waiting for work")
					work, err := stream.Recv()
					if err != nil {
						log.WithError(err).Warn("failed to receive work!")
						time.Sleep(10 * time.Second)
						break // retry connection
					}

					log.Debug("received work: %+v", work)
					workUnitChan <- poolWork{pool, work}
					log.Debug("assigned work")
				}

			}
		}(pool)
	}

	// miner threads
	for i := 0; i < numThreads; i++ {
		go func(i int) {
			minerThread(i)
		}(i)
	}

	// print stats every 15s
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

	// one-time print after we've computed a single hash
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

	// publish work to workers
	log := baseLogger.WithField("thread", "main")
	for workUnit := range workUnitChan {
		currentShareCount := uint64(0)
		currentBlockHeight := uint64(0)
		if currentPool != "" {
			currentShareCount = poolStates[currentPool].shares
			if currentWorkUnit != nil {
				currentBlockHeight = uint64(currentWorkUnit.Header.BlockHeight)
			}
		}

		newShareCount := poolStates[workUnit.pool].shares
		newBlockHeight := uint64(workUnit.work.Header.BlockHeight)

		if newShareCount > currentShareCount && newBlockHeight <= currentBlockHeight {
			log.
				WithField("pool", workUnit.pool).
				WithField("work-id", workUnit.work.WorkId).
				Debug("ignoring work unit")
			continue
		}

		log.
			WithField("pool", workUnit.pool).
			WithField("work-id", workUnit.work.WorkId).
			Debug("accepting work unit")
		currentWorkUnitLock.Lock()
		currentPool = workUnit.pool
		currentWorkUnit = workUnit.work
		currentWorkID = workUnit.work.WorkId
		currentWorkUnitLock.Unlock()

		if newBlockHeight != currentBlockHeight {
			log.WithField("pool", workUnit.pool).WithField("block-height", newBlockHeight).Info("work for new blockheight")
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

		workPool := currentPool
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
				log.
					WithField("pool", workPool).
					WithField("context", fmt.Sprintf("%X", snowContext)).
					Info("passable solution")

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

				reply, err := poolStates[workPool].Submit(&mining_pool.WorkSubmitRequest{
					WorkId: workID,
					Header: &workHeader,
				})
				if err != nil {
					log.WithError(err).Warn("share submission failed")
					continue
				}

				log.
					WithField("pool", workPool).
					WithField("reply", reply).
					Info("submitted!")
			}
		}
	}
}
