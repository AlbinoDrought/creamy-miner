package extra

import (
	"sync/atomic"

	"go.snowblossom/internal/sblib"
	"go.snowblossom/internal/sbproto/snowblossom"
)

func MergeFields(allFields []map[int]sblib.SnowMerkleProof) map[int]sblib.SnowMerkleProof {
	byFieldID := map[int][]sblib.SnowMerkleProof{}

	for _, fields := range allFields {
		for fieldID, field := range fields {
			if byFieldID[fieldID] == nil {
				byFieldID[fieldID] = []sblib.SnowMerkleProof{}
			}

			byFieldID[fieldID] = append(byFieldID[fieldID], field)
		}
	}

	mergedFields := map[int]sblib.SnowMerkleProof{}

	for fieldID, fields := range byFieldID {
		if len(fields) == 1 {
			// only one field, don't wrap
			mergedFields[fieldID] = fields[0]
			continue
		}

		mergedFields[fieldID] = &mergedField{
			fields:      fields,
			fieldLength: uint32(len(fields)),
		}
	}

	return mergedFields
}

type mergedField struct {
	fields            []sblib.SnowMerkleProof
	fieldLength       uint32
	roundRobinCounter uint32
}

func (mf *mergedField) Close() error {
	var lastErr error
	for _, field := range mf.fields {
		if err := field.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (mf *mergedField) TotalWords() int64 {
	// should be the same between all fields, or this is dumb
	return mf.fields[0].TotalWords()
}

func (mf *mergedField) ReadWord(wordIndex int64, buffer []byte, currentDepth int) error {
	ctr := atomic.AddUint32(&mf.roundRobinCounter, 1) % mf.fieldLength
	return mf.fields[ctr].ReadWord(wordIndex, buffer, currentDepth)
}

func (mf *mergedField) GetProof(wordIndex int64) (*snowblossom.SnowPowProof, error) {
	ctr := atomic.AddUint32(&mf.roundRobinCounter, 1) % mf.fieldLength
	return mf.fields[ctr].GetProof(wordIndex)
}
