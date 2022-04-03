package sblib_test

import (
	"encoding/hex"
	"fmt"
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

func TestLessThanTarget(t *testing.T) {
	tests := []struct {
		foundHash string
		target    string
		want      bool
	}{
		{
			"000001F5FA05E2E9599E976850469CB6A4084F9EE156F4D0894E4E281B00DE71",
			"0000040000000000000000000000000000000000000000000000000000000000",
			true,
		},
		{
			"000011F5FA05E2E9599E976850469CB6A4084F9EE156F4D0894E4E281B00DE71",
			"0000040000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"000003F5FA05E2E9599E976850469CB6A4084F9EE156F4D0894E4E281B00DE71",
			"0000040000000000000000000000000000000000000000000000000000000000",
			true,
		},
		{
			"0000040000000000000000000000000000000000000000000000000000000000",
			"0000040000000000000000000000000000000000000000000000000000000000",
			false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%v: %v vs %v is less: %v", i, tt.foundHash, tt.target, tt.want), func(t *testing.T) {
			foundHashBytes, _ := hex.DecodeString(tt.foundHash)
			targetBytes, _ := hex.DecodeString(tt.target)
			if got := sblib.LessThanTarget(foundHashBytes, targetBytes); got != tt.want {
				t.Errorf("LessThanTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
