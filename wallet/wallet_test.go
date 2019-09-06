package wallet

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletIsNew(t *testing.T) {
	assert := assert.New(t)

	walletID := uuid.NewV4()
	providerID := uuid.NewV4()

	wallet := Info{
		ID:          walletID.String(),
		Provider:    "uphold",
		ProviderID:  providerID.String(),
		AltCurrency: nil,
		PublicKey:   "-",
		LastBalance: nil,
		CreatedAt:   time.Now().AddDate(0, 0, -1),
	}

	assert.Equal(false, wallet.IsNew())
	wallet.CreatedAt = time.Now()
	assert.Equal(true, wallet.IsNew())
}
