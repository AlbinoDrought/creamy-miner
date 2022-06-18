package sblib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash"
	"sort"

	"go.snowblossom/internal/sbproto/snowblossom"
)

func HashHeaderBits(header *snowblossom.BlockHeader, nonce []byte, md hash.Hash) ([]byte, error) {
	if header.Version != 1 && header.Version != 2 {
		return nil, errors.New("only versions 1 and 2 are supported")
	}

	md.Reset()

	intData := []byte{
		byte(header.Version >> 24),
		byte(header.Version >> 16),
		byte(header.Version >> 8),
		byte(header.Version & 0xFF),

		byte(header.BlockHeight >> 24),
		byte(header.BlockHeight >> 16),
		byte(header.BlockHeight >> 8),
		byte(header.BlockHeight & 0xFF),

		byte(header.Timestamp >> 56),
		byte(header.Timestamp >> 48),
		byte(header.Timestamp >> 40),
		byte(header.Timestamp >> 32),
		byte(header.Timestamp >> 24),
		byte(header.Timestamp >> 16),
		byte(header.Timestamp >> 8),
		byte(header.Timestamp & 0xFF),

		byte(header.SnowField >> 24),
		byte(header.SnowField >> 16),
		byte(header.SnowField >> 8),
		byte(header.SnowField & 0xFF),
	}

	if _, err := md.Write(nonce); err != nil {
		return nil, err
	}
	if _, err := md.Write(intData); err != nil {
		return nil, err
	}
	if _, err := md.Write(header.PrevBlockHash); err != nil {
		return nil, err
	}
	if _, err := md.Write(header.MerkleRootHash); err != nil {
		return nil, err
	}
	if _, err := md.Write(header.UtxoRootHash); err != nil {
		return nil, err
	}
	if _, err := md.Write(header.Target); err != nil {
		return nil, err
	}

	if header.Version == 2 {
		if _, err := md.Write([]byte{
			byte(header.ShardId >> 24),
			byte(header.ShardId >> 16),
			byte(header.ShardId >> 8),
			byte(header.ShardId & 0xFF),

			byte(header.TxDataSizeSum >> 24),
			byte(header.TxDataSizeSum >> 16),
			byte(header.TxDataSizeSum >> 8),
			byte(header.TxDataSizeSum & 0xFF),

			byte(header.TxCount >> 24),
			byte(header.TxCount >> 16),
			byte(header.TxCount >> 8),
			byte(header.TxCount & 0xFF),
		}); err != nil {
			return nil, err
		}

		// todo: maps are randomly ordered in golang
		for id, value := range header.GetShardExportRootHash() {
			if _, err := md.Write([]byte{
				byte(id >> 24),
				byte(id >> 16),
				byte(id >> 8),
				byte(id & 0xFF),
			}); err != nil {
				return nil, err
			}

			if _, err := md.Write(value); err != nil {
				return nil, err
			}
		}

		i := 0
		importShardIDs := make([]int32, len(header.ShardImport))
		for importShardID := range header.ShardImport {
			importShardIDs[i] = importShardID
			i++
		}
		sort.Slice(importShardIDs, func(i, j int) bool {
			return importShardIDs[i] < importShardIDs[j]
		})

		for importShardID, bil := range header.GetShardImport() {
			i = 0
			heightMapKeys := make([]int32, len(bil.HeightMap))
			for heightMapKey := range bil.HeightMap {
				heightMapKeys[i] = heightMapKey
				i++
			}
			sort.Slice(heightMapKeys, func(i, j int) bool {
				return heightMapKeys[i] < heightMapKeys[j]
			})

			for importHeight, value := range heightMapKeys {
				if _, err := md.Write([]byte{
					byte(importShardID >> 24),
					byte(importShardID >> 16),
					byte(importShardID >> 8),
					byte(importShardID & 0xFF),

					byte(importHeight >> 24),
					byte(importHeight >> 16),
					byte(importHeight >> 8),
					byte(importHeight & 0xFF),
				}); err != nil {
					return nil, err
				}

				// java uses toByteArray, not sure if that is the same
				if _, err := md.Write([]byte{
					byte(value >> 24),
					byte(value >> 16),
					byte(value >> 8),
					byte(value & 0xFF),
				}); err != nil {
					return nil, err
				}
			}
		}
	}

	return md.Sum(nil), nil
}

func GetNextSnowFieldIndex(context []byte, wordCount int64, md hash.Hash) int64 {
	md.Reset()
	md.Write(context)
	// todo: for some reason Globals.BLOCKCHAIN_HASH_LEN is 32, uses a 32-len buffer here?
	buff := md.Sum(nil)
	buff[0] = 0
	return int64(binary.BigEndian.Uint64(buff)) % wordCount
}

func GetNextContext(prevContext []byte, foundData []byte, md hash.Hash) []byte {
	md.Reset()
	md.Write(prevContext)
	md.Write(foundData)
	return md.Sum(nil)
}

func LessThanTarget(foundHash []byte, target []byte) bool {
	// Quickly reject nonsense
	for i := 0; i < 8; i++ {
		// If the target is not zero, done with loop
		if target[i] != 0 {
			break
		}
		// If the target is zero, but found is not, fail
		if foundHash[i] != 0 {
			return false
		}
	}

	if len(foundHash) != len(target) {
		return false
	}
	return bytes.Compare(foundHash, target) == -1
}
