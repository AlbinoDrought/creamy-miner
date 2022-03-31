package sblib

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"

	"go.snowblossom/internal/sbconst"
	"go.snowblossom/internal/sbproto/snowblossom"
)

func getNumberOfDecks(totalWords int64) int {
	count := 0
	for totalWords > sbconst.DeckEntries {
		totalWords /= sbconst.DeckEntries
		count++
	}
	return count
}

type ReadAtCloser interface {
	io.ReaderAt
	io.Closer
}

type snowMerkleProof struct {
	snowFile       ReadAtCloser
	snowFileLength int64
	deckMap        map[int64]ReadAtCloser
	totalWords     int64

	closeLock sync.Mutex
}

func (smp *snowMerkleProof) TotalWords() int64 {
	return smp.totalWords
}

func (smp *snowMerkleProof) ReadWord(wordIndex int64, buffer []byte, currentDepth int) error {
	wordPos := wordIndex * sbconst.HashLenLong
	_, err := smp.snowFile.ReadAt(buffer, wordPos)
	if err != nil {
		return err
	}
	return nil
}

func (smp *snowMerkleProof) GetProof(wordIndex int64) (*snowblossom.SnowPowProof, error) {
	// todo: use more efficient datastructure, array will need to be resized+copy each addition
	partners := &partners{[][]byte{}}

	// todo: reuse this
	md := NewSnowMerkleMessageDigest()

	if _, err := smp.getInnerProof(md, partners, wordIndex, 0, smp.totalWords); err != nil {
		return nil, errors.Wrapf(err, "GetProof -> getInnerProof failed for wordIndex %v", wordIndex)
	}

	proof := &snowblossom.SnowPowProof{}
	proof.WordIdx = wordIndex
	proof.MerkleComponent = partners.data

	return proof, nil
}

type partners struct {
	data [][]byte
}

func (smp *snowMerkleProof) getInnerProof(md hash.Hash, partners *partners, targetWordIndex int64, start int64, end int64) ([]byte, error) {
	inside := (start <= targetWordIndex) && (targetWordIndex < end)

	dist := end - start
	if !inside {
		if deckFile, ok := smp.deckMap[dist]; ok {
			deckPos := sbconst.HashLenLong * (start / dist)
			buff := make([]byte, sbconst.SnowMerkleHashLen)
			if _, err := deckFile.ReadAt(buff, deckPos); err != nil {
				return nil, errors.Wrapf(err, "getInnerProof !inside deckFile %v ReadAt %v failed", dist, deckPos)
			}
			return buff, nil
		}

		if dist == 1 {
			wordPos := start * sbconst.HashLenLong
			buff := make([]byte, sbconst.SnowMerkleHashLen)
			if _, err := smp.snowFile.ReadAt(buff, wordPos); err != nil {
				return nil, errors.Wrapf(err, "getInnerProof !inside snowFile ReadAt %v failed", wordPos)
			}
			return buff, nil
		}
		mid := (start + end) / 2

		left, err := smp.getInnerProof(md, partners, targetWordIndex, start, mid)
		if err != nil {
			return nil, errors.Wrap(err, "getInnerProof !inside left read failed")
		}
		right, err := smp.getInnerProof(md, partners, targetWordIndex, mid, end)
		if err != nil {
			return nil, errors.Wrap(err, "getInnerProof !inside right read failed")
		}

		md.Reset()

		md.Write(left)
		md.Write(right)

		return md.Sum(nil), nil
	} else {
		// we are inside
		if dist == 1 {
			wordPos := start * sbconst.HashLenLong
			buff := make([]byte, sbconst.SnowMerkleHashLen)

			if _, err := smp.snowFile.ReadAt(buff, wordPos); err != nil {
				return nil, errors.Wrapf(err, "getInnerProof inside snowFile ReadAt %v failed", wordPos)
			}

			partners.data = append(partners.data, buff)
			return nil, nil
		}

		mid := (start + end) / 2

		left, err := smp.getInnerProof(md, partners, targetWordIndex, start, mid)
		if err != nil {
			return nil, errors.Wrap(err, "getInnerProof inside left read failed")
		}
		right, err := smp.getInnerProof(md, partners, targetWordIndex, mid, end)
		if err != nil {
			return nil, errors.Wrap(err, "getInnerProof inside right read failed")
		}

		if targetWordIndex < mid {
			partners.data = append(partners.data, right)

			// replacement for Assert.assertNull:
			if left != nil {
				return nil, errors.New("left unexpectedly non-nil")
			}
		} else {
			partners.data = append(partners.data, left)

			// replacement for Assert.assertNull:
			if right != nil {
				return nil, errors.New("right unexpectedly non-nil")
			}
		}

		return nil, nil
	}
}

func (smp *snowMerkleProof) Close() error {
	var lastErr error

	smp.closeLock.Lock()
	defer smp.closeLock.Unlock()

	if smp.snowFile != nil {
		if err := smp.snowFile.Close(); err != nil {
			lastErr = err
		}
		smp.snowFile = nil
	}

	for _, f := range smp.deckMap {
		if err := f.Close(); err != nil {
			lastErr = err
		}
	}
	smp.deckMap = map[int64]ReadAtCloser{}

	if lastErr != nil {
		return lastErr
	}

	return nil
}

type SnowMerkleProof interface {
	io.Closer

	TotalWords() int64
	ReadWord(wordIndex int64, buffer []byte, currentDepth int) error
	GetProof(wordIndex int64) (*snowblossom.SnowPowProof, error)
}

func NewSnowMerkleProof(snowPath string, snowBase string) (SnowMerkleProof, error) {
	snowFilePath := path.Join(snowPath, snowBase+".snow")

	snowFileStat, err := os.Stat(snowFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat snowFile %v", snowFilePath)
	}
	snowFileLength := snowFileStat.Size()

	snowFile, err := os.Open(snowFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open snowFile %v", snowFilePath)
	}

	totalWords := snowFileLength / sbconst.SnowMerkleHashLen
	deckCount := getNumberOfDecks(totalWords)

	deckMap := make(map[int64]ReadAtCloser, deckCount)

	h := int64(1)
	for i := 0; i < deckCount; i++ {
		h = h * sbconst.DeckEntries
		letter := byte('a' + i)

		deckPath := path.Join(snowPath, snowBase+".deck."+string(letter))
		deckStat, err := os.Stat(deckPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to stat deckFile %v", deckPath)
		}

		expectedLength := snowFileLength / h
		if deckStat.Size() != expectedLength {
			return nil, fmt.Errorf("unexpected length on %v.deck.%v", snowBase, letter)
		}

		deckFile, err := os.Open(deckPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open deckFile %v", deckPath)
		}

		deckMap[h] = deckFile
	}

	return &snowMerkleProof{
		snowFile:       snowFile,
		snowFileLength: snowFileLength,
		deckMap:        deckMap,
		totalWords:     totalWords,
	}, nil
}
