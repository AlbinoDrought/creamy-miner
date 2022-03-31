package sblib

import (
	"hash"

	"github.com/aead/skein/skein256"
)

func NewMessageDigest() hash.Hash {
	return skein256.New(32, nil)
}

func NewSnowMerkleMessageDigest() hash.Hash {
	return skein256.New(16, nil)
}
