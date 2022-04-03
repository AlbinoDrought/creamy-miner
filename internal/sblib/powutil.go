package sblib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash"

	"go.snowblossom/internal/sbproto/snowblossom"
)

func HashHeaderBits(header *snowblossom.BlockHeader, nonce []byte, md hash.Hash) ([]byte, error) {
	if header.Version != 1 {
		return nil, errors.New("only version 1 is supported")
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

	/*
		if header.Version == 2 {
			if _, err := md.Write([]byte{
				byte(header.GetShardId() >> 24),
				byte(header.GetShardId() >> 16),
				byte(header.GetShardId() >> 8),
				byte(header.GetShardId() & 0xFF),

				byte(header.GetTxDataSizeSum() >> 24),
				byte(header.GetTxDataSizeSum() >> 16),
				byte(header.GetTxDataSizeSum() >> 8),
				byte(header.GetTxDataSizeSum() & 0xFF),

				byte(header.GetTxCount() >> 24),
				byte(header.GetTxCount() >> 16),
				byte(header.GetTxCount() >> 8),
				byte(header.GetTxCount() & 0xFF),
			}); err != nil {
				return nil, err
			}

			header.GetShardExportRootHash()
		}
	*/

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
