package promotion

import (
	"context"
	"time"

	"github.com/brave-intl/bat-go/utils/ads"
	"github.com/brave-intl/bat-go/wallet"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Promotion includes information about a particular promotion
type Promotion struct {
	ID                  uuid.UUID       `json:"id" db:"id"`
	CreatedAt           time.Time       `json:"createdAt" db:"created_at"`
	ExpiresAt           time.Time       `json:"expiresAt" db:"expires_at"`
	Version             int             `json:"version" db:"version"`
	SuggestionsPerGrant int             `json:"suggestionsPerGrant" db:"suggestions_per_grant"`
	ApproximateValue    decimal.Decimal `json:"approximateValue" db:"approximate_value"`
	Type                string          `json:"type" db:"promotion_type"`
	RemainingGrants     int             `json:"-" db:"remaining_grants"`
	Active              bool            `json:"-" db:"active"`
	Available           bool            `json:"available" db:"available"`
	//ClaimableUntil      time.Time
	//PublicKeys          []string
}

// GetOrCreateWallet attempts to retrieve wallet info from the local datastore, falling back to the ledger
func (service *Service) GetOrCreateWallet(ctx context.Context, walletID uuid.UUID) (*wallet.Info, error) {
	wallet, err := service.datastore.GetWallet(walletID)
	if err != nil {
		return nil, errors.Wrap(err, "Error looking up wallet")
	}

	if wallet == nil {
		wallet, err = service.ledgerClient.GetWallet(ctx, walletID)
		if err != nil {
			return nil, errors.Wrap(err, "Error looking up wallet")
		}
		if wallet != nil {
			err = service.datastore.InsertWallet(wallet)
			if err != nil {
				return nil, errors.Wrap(err, "Error saving wallet")
			}
		}
	}
	return wallet, nil
}

// GetAvailablePromotions first looks up the wallet and then retrieves available promotions
func (service *Service) GetAvailablePromotions(
	ctx context.Context,
	walletID uuid.UUID,
	availableTypes map[string]bool,
) ([]Promotion, error) {
	wallet, err := service.GetOrCreateWallet(ctx, walletID)
	if err != nil {
		return []Promotion{}, err
	}
	if wallet.IsNew() {
		return []Promotion{}, errors.New("wallet is too new to receive grants")
	}
	promotions, err := service.datastore.GetAvailablePromotionsForWallet(wallet)
	if err != nil {
		return []Promotion{}, err
	}
	filtered := []Promotion{}
	for _, promotion := range promotions {
		if availableTypes[promotion.Type] {
			filtered = append(filtered, promotion)
		}
	}
	return filtered, nil
}

func (service *Service) PromotionTypesAvailable(countryCode string) (map[string]bool, error) {
	cohorts := map[string]bool{
		"control": true,
	}
	countriesAdsAvailable, err := ads.AvailableIn()
	if err != nil {
		return cohorts, err
	}
	for _, countryAdsAvailable := range countriesAdsAvailable {
		if countryCode == countryAdsAvailable {
			cohorts["ads"] = true
		}
	}
	if len(cohorts) == 0 {
		cohorts["ugp"] = true
	}
	return cohorts, nil
}
