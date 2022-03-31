package main

import (
	"context"
	"crypto/rand"
	"log"

	"go.snowblossom/internal/sbconst"
	"go.snowblossom/internal/sblib"
	"go.snowblossom/internal/sbproto/mining_pool"
	"go.snowblossom/internal/sbproto/snowblossom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fields, err := sblib.LoadFields("fields")
	if err != nil {
		panic(err)
	}

	// conn, err := grpc.Dial("snowypool.com:23380", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:23380", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := mining_pool.NewMiningPoolServiceClient(conn)

	stream, err := client.GetWork(context.Background(), &mining_pool.GetWorkRequest{
		PayToAddress: "c04rt84spfjc9xy88snx5r256qv0tmy664zcdrnc",
	})
	if err != nil {
		panic(err)
	}

	work, err := stream.Recv()
	if err != nil {
		panic(err)
	}

	log.Printf("received work: %+v", work)

	/*
		start := time.Now()

		count := int64(0)

		for i := 0; i < runtime.NumCPU()*10; i++ {
			go func() {
				md := skein.New256(nil)

				for {
					_, err := hashHeaderBits(work.Header, work.Header.Nonce, md)
					if err != nil {
						panic(err)
					}

					atomic.AddInt64(&count, 1) // todo: not sure if the count given by arktika is /6
				}
			}()
		}

		// targetHashes := int64(10 * (1 << 20))
		targetHashes := int64(0)
		for {
			time.Sleep(15 * time.Second)
			currentCount := atomic.AddInt64(&count, 0)
			if currentCount >= targetHashes {
				end := time.Since(start)
				perSecond := float64(currentCount) / end.Seconds()
				log.Printf("%vM hashes in %v", currentCount/(1<<20), end)
				log.Printf("that is %.2f/sec", perSecond)
				log.Printf("depending how Arktika calculates it, this could be %.2f/sec or %.2f/sec", perSecond/6, perSecond/7)
				os.Exit(0)
			}
		}
	*/
	md := sblib.NewMessageDigest()

	nonce := make([]byte, sbconst.NonceLength)
	if _, err = rand.Read(nonce); err != nil {
		panic(err)
	}
	copy(nonce, work.Header.Nonce[:4])
	/*
		nonce[0] = work.Header.Nonce[0]
		nonce[1] = work.Header.Nonce[1]
		nonce[2] = work.Header.Nonce[2]
		nonce[3] = work.Header.Nonce[3]
	*/

	snowContext, err := sblib.HashHeaderBits(work.Header, nonce, md)
	if err != nil {
		panic(err)
	}

	merkleProof := fields[int(work.Header.SnowField)]
	wordBuffer := make([]byte, sbconst.SnowMerkleHashLen)
	wordIndexHistory := make([]int64, sbconst.POWLookPasses)

	for pass := 0; pass < sbconst.POWLookPasses; pass++ {
		wordIndex := sblib.GetNextSnowFieldIndex(snowContext, merkleProof.TotalWords(), md)
		wordIndexHistory[pass] = wordIndex
		err = merkleProof.ReadWord(wordIndex, wordBuffer, pass)
		if err != nil {
			panic(err)
		}
		snowContext = sblib.GetNextContext(snowContext, wordBuffer, md)
	}

	if sblib.LessThanTarget(snowContext, work.ReportTarget) {
		log.Printf("Found passable solution: %X", snowContext)

		work.Header.Nonce = nonce

		powProofs := make([]*snowblossom.SnowPowProof, len(wordIndexHistory))
		for pass, wordIndex := range wordIndexHistory {
			powProofs[pass], err = merkleProof.GetProof(wordIndex)
			if err != nil {
				panic(err)
			}
		}
		work.Header.PowProof = powProofs
		work.Header.SnowHash = snowContext

		resp, err := client.SubmitWork(context.Background(), &mining_pool.WorkSubmitRequest{
			WorkId: work.WorkId,
			Header: work.Header,
		})
		log.Printf("%+v", resp)
		log.Printf("%+v", err)
	} else {
		log.Println("Solution not passable")
	}
}
