package clustername

import (
	"fmt"

	"github.com/kris-nova/kubicorn/cutil/rand"
)

// GetClusterName returns a cluster name based on a provided
// provider shorthand.
func GetClusterName(providerShorthand string) string {
	return fmt.Sprintf("e2e-%s-%s", providerShorthand, randStringRunes(6))
}

// randStringRunes returns random alphanumeric string with the given length.
func randStringRunes(length int) string {
	var letterRunes = []rune("0123456789abcdef")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.GenerateRandomInt(0, len(letterRunes))]
	}
	return string(b)
}
