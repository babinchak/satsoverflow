package question

import (
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcutil/bech32"
)

func Question_bech32(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	hash = hash[:8]
	conv, err := bech32.ConvertBits(hash, 8, 5, true)
	if err != nil {
		log.Fatalf("Error converting bits: %v\n", err)
	}
	encoded, err := bech32.Encode("question", conv)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	return encoded
}
