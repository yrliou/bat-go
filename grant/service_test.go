package grant

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSettlementBalance(t *testing.T) {
	if os.Getenv("UPHOLD_ACCESS_TOKEN") == "" {
		t.Skip("skipping test; UPHOLD_ACCESS_TOKEN not set")
	}
	assert := assert.New(t)

	current, err := GetSettlementBalance()
	assert.NoError(err)
	assert.Equal(true, current.IsPositive(), "value should be positive")
}
