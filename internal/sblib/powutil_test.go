package sblib_test

import (
	"encoding/hex"
	"testing"

	"go.snowblossom/internal/sblib"
)

func TestNextSnowFieldIndex(t *testing.T) {
	context := make([]byte, 32)
	for i := range context {
		context[i] = byte(i)
	}

	v := sblib.GetNextSnowFieldIndex(context, 1024*1024*1024*32, sblib.NewMessageDigest())
	if v != 19323217346 {
		t.Error("unexpected GetNextSnowFieldIndex value", v)
	}
}

func TestGetNextContext(t *testing.T) {
	// from new Random(9182L)
	context, _ := hex.DecodeString("721ffffa23428ac81119f9ae9f4e3f2d79c4d696b78fd6a998041b7ec9cf4bd3")
	word, _ := hex.DecodeString("81a1de38936b38a9205816e1fefc2d1ebb446830eb765e8b32b7b92cc67a51f0")

	newContext := sblib.GetNextContext(context, word, sblib.NewMessageDigest())
	if hex.EncodeToString(newContext) != "db176f021b168412a05904e7e729a072943cc8a4f266a4579654b8aa0fecfb8b" {
		t.Error("unexpected GetNextContext value")
	}
}
