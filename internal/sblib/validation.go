package sblib

import (
	"errors"
	"fmt"

	"go.snowblossom/internal/sbconst"
	"go.snowblossom/internal/sbproto/snowblossom"
)

func CheckProof(proof *snowblossom.SnowPowProof, expectedMerkleRoot []byte, snowFieldSize int64) error {
	targetIndex := proof.WordIdx
	wordCount := snowFieldSize / sbconst.HashLenLong
	if targetIndex < 0 {
		return errors.New("bad target index: less than zero")
	}
	if targetIndex >= wordCount {
		return errors.New("bad target index: greater than word count")
	}

	// todo: reuse
	md := NewSnowMerkleMessageDigest()
	stack := proof.MerkleComponent
	tip := 0

	if len(stack) <= 0 {
		return errors.New("empty merkle components")
	}

	currentHash := stack[tip]
	tip++
	if currentHash == nil {
		return errors.New("top merkle component is nil")
	}

	start := targetIndex
	end := targetIndex
	dist := int64(1)

	for (len(stack)-tip) > 0 && end <= wordCount {
		dist *= 2
		start = start - (start % dist)
		end = start + dist

		mid := (start + end) / 2

		var left []byte
		var right []byte
		if targetIndex < mid {
			left = currentHash
			right = stack[tip]
		} else {
			left = stack[tip]
			right = currentHash
		}
		tip++

		md.Reset()

		md.Write(left)
		md.Write(right)

		currentHash = md.Sum(nil)
	}

	if start != 0 {
		return errors.New("expected start to be 0")
	}

	if end != wordCount {
		return errors.New("expected end to be word count")
	}

	for i := range expectedMerkleRoot {
		if currentHash[i] != expectedMerkleRoot[i] {
			return fmt.Errorf("does not match expected merkle root: got %X want %X", currentHash, expectedMerkleRoot)
		}
	}

	return nil
}
