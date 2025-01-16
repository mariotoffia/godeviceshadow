package persistutils

import (
	"crypto/rand"
	"math/big"

	"fmt"
)

// Id generates a random 6-digit ID (e.g., 123456). If a prefix is provided,
// it is added to the start of the ID (e.g., "prefix-123456").
func Id(prefix ...string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))

	if err != nil {
		panic(fmt.Sprintf("failed to read crypto random: %v", err))
	}

	idStr := fmt.Sprintf("%06d", n.Int64())

	if len(prefix) > 0 {
		return prefix[0] + "-" + idStr
	}

	return idStr
}
