package datasources

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/charmbracelet/log"
)

func ComputeHash(doc string) Hash {
	// Dumb af, but it's a cheap way to specific the most generic thing
	// you can :-/
	var v interface{}
	err := json.Unmarshal([]byte(doc), &v) // NB: You should handle errors :-/

	if err != nil {
		log.Error("Failed to marshal sorted body", "error", err)
		return Hash("")
	}

	cdoc, _ := json.Marshal(v)
	sum := sha256.Sum256(cdoc)
	return Hash(hex.EncodeToString(sum[0:]))
}
