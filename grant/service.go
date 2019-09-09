package grant

import (
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/brave-intl/bat-go/utils/altcurrency"
	"github.com/brave-intl/bat-go/utils/httpsignature"
	"github.com/brave-intl/bat-go/wallet"
	"github.com/brave-intl/bat-go/wallet/provider"
	"github.com/brave-intl/bat-go/wallet/provider/uphold"
	"github.com/garyburd/redigo/redis"
	raven "github.com/getsentry/raven-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/ed25519"
)

const (
	lowerTxLimit        = 1
	upperTxLimit        = 120
	ninetyDaysInSeconds = 60 * 60 * 24 * 90
	productionEnv       = "production"
)

var (
	// SettlementDestination is the address of the settlement wallet
	SettlementDestination = os.Getenv("BAT_SETTLEMENT_ADDRESS")
	// GrantSignatorPublicKeyHex is the hex encoded public key of the keypair used to sign grants
	GrantSignatorPublicKeyHex    = os.Getenv("GRANT_SIGNATOR_PUBLIC_KEY")
	grantWalletPublicKeyHex      = os.Getenv("GRANT_WALLET_PUBLIC_KEY")
	grantWalletPrivateKeyHex     = os.Getenv("GRANT_WALLET_PRIVATE_KEY")
	grantWalletCardID            = os.Getenv("GRANT_WALLET_CARD_ID")
	grantPublicKey               ed25519.PublicKey
	grantWallet                  *uphold.Wallet
	refreshBalance               = true  // for testing we can disable balance refresh
	testSubmit                   = true  // for testing we can disable testing tx submit
	registerGrantInstrumentation = true  // for testing we can disable grant claim / redeem instrumentation registration
	safeMode                     = false // if set true disables grant redemption
	claimedGrantsCounter         = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "claimed_grants_total",
			Help: "Number of grants claimed since start.",
		},
		[]string{},
	)
	redeemedGrantsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redeemed_grants_total",
			Help: "Number of grants redeemed since start.",
		},
		[]string{"promotionId"},
	)
	fundsReceivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "funds_received_count",
			Help: "a count of the number of bat added to the settlement wallet",
		},
		[]string{},
	)
)

// Service contains datastore and redis connections as well as prometheus metrics
type Service struct {
	datastore                 Datastore
	redisPool                 *redis.Pool
	outstandingGrantCountDesc *prometheus.Desc
	completedGrantCountDesc   *prometheus.Desc
	grantWalletBalanceDesc    *prometheus.Desc
}

// InitService initializes the grant service
func InitService(datastore Datastore, redisPool *redis.Pool) (*Service, error) {
	gs := &Service{
		datastore: datastore,
		redisPool: redisPool,
		outstandingGrantCountDesc: prometheus.NewDesc(
			"outstanding_grants_total",
			"Outstanding grants that have been claimed and have not expired.",
			[]string{},
			prometheus.Labels{},
		),
		completedGrantCountDesc: prometheus.NewDesc(
			"completed_grants_total",
			"Completed grants that have been redeemed.",
			[]string{"promotionId"},
			prometheus.Labels{},
		),
		grantWalletBalanceDesc: prometheus.NewDesc(
			"grant_wallet_balance",
			"A gauge of the grant wallet remaining balance.",
			[]string{},
			prometheus.Labels{},
		),
	}

	var err error
	grantPublicKey, err = hex.DecodeString(GrantSignatorPublicKeyHex)
	if err != nil {
		return nil, err
	}

	if os.Getenv("ENV") == productionEnv && !refreshBalance {
		return nil, errors.New("refreshBalance must be true in production")
	}
	if os.Getenv("ENV") == productionEnv && !testSubmit {
		return nil, errors.New("testSubmit must be true in production")
	}

	if len(grantWalletCardID) > 0 {
		var info wallet.Info
		info.Provider = "uphold"
		info.ProviderID = grantWalletCardID
		{
			tmp := altcurrency.BAT
			info.AltCurrency = &tmp
		}

		var pubKey httpsignature.Ed25519PubKey
		var privKey ed25519.PrivateKey
		var err error

		pubKey, err = hex.DecodeString(grantWalletPublicKeyHex)
		if err != nil {
			return nil, err
		}
		privKey, err = hex.DecodeString(grantWalletPrivateKeyHex)
		if err != nil {
			return nil, err
		}

		grantWallet, err = uphold.New(info, privKey, pubKey)
		if err != nil {
			return nil, err
		}
	} else if os.Getenv("ENV") == productionEnv {
		return nil, errors.New("GRANT_WALLET_CARD_ID must be set in production")
	}

	if registerGrantInstrumentation {
		if datastore != nil {
			prometheus.MustRegister(gs)
		}

		prometheus.MustRegister(claimedGrantsCounter)
		prometheus.MustRegister(redeemedGrantsCounter)
		prometheus.MustRegister(fundsReceivedCounter)
	}

	return gs, nil
}

// Describe returns all descriptions of the collector.
// We implement this and the Collect function to fulfill the prometheus.Collector interface
func (gs *Service) Describe(ch chan<- *prometheus.Desc) {
	ch <- gs.outstandingGrantCountDesc
	ch <- gs.completedGrantCountDesc
	ch <- gs.grantWalletBalanceDesc
}

// Collect returns the current state of all metrics of the collector.
// We implement this and the Describe function to fulfill the prometheus.Collector interface
func (gs *Service) Collect(ch chan<- prometheus.Metric) {
	ogCount, err := gs.datastore.GetOutstandingGrantCount()
	if err != nil {
		raven.CaptureError(err, map[string]string{})
		return
	}
	ch <- prometheus.MustNewConstMetric(
		gs.outstandingGrantCountDesc,
		prometheus.GaugeValue,
		float64(ogCount),
	)
	redeemedCounts, err := gs.datastore.GetRedeemedCountByPromotion()
	if err != nil {
		raven.CaptureError(err, map[string]string{})
		return
	}
	for promotionID, completedCount := range redeemedCounts {
		ch <- prometheus.MustNewConstMetric(
			gs.completedGrantCountDesc,
			prometheus.GaugeValue,
			float64(completedCount),
			promotionID,
		)
	}

	balance, err := grantWallet.GetBalance(true)
	if err != nil {
		raven.CaptureError(err, map[string]string{})
		return
	}

	spendable, _ := grantWallet.GetWalletInfo().AltCurrency.FromProbi(balance.SpendableProbi).Float64()

	ch <- prometheus.MustNewConstMetric(
		gs.grantWalletBalanceDesc,
		prometheus.GaugeValue,
		spendable,
	)

	if isMasterMachine() {
		gs.pullWalletMetrics()
	}
}

func isMasterMachine() bool {
	return true
}

func (gs *Service) pullWalletMetrics() error {
	var lastBalance decimal.Decimal
	walletMetricKey := "settlement-wallet:balance"
	conn := gs.redisPool.Get()

	err := conn.Send("GET", walletMetricKey)
	if err != nil {
		return err
	}

	lastBalanceCachedInterface, err := conn.Receive()
	if err != nil {
		return err
	}

	lastBalanceCached := lastBalanceCachedInterface.(string)
	currentBalance, err := GetSettlementBalance()
	if err != nil {
		return err
	}

	if lastBalanceCached == "" {
		lastBalance = currentBalance
	} else {
		lastBalance, err = decimal.NewFromString(lastBalanceCached)
		if err != nil {
			return err
		}
	}

	delta, _ := currentBalance.Sub(lastBalance).Float64()
	fundsReceivedCounter.With(prometheus.Labels{}).Add(delta)

	conn.Send("SET", walletMetricKey, currentBalance.String())
	conn.Do("EXEC")
	return nil
}

func GetSettlementBalance() (decimal.Decimal, error) {
	zero, err := decimal.NewFromString("0")
	if err != nil {
		return zero, err
	}
	alt, err := altcurrency.FromString("BAT")
	if err != nil {
		return zero, err
	}
	settlementWallet, err := provider.GetWallet(wallet.Info{
		ID:          uuid.NewV4().String(),
		Provider:    "uphold",
		ProviderID:  SettlementDestination,
		PublicKey:   GrantSignatorPublicKeyHex,
		AltCurrency: &alt,
		LastBalance: nil,
		CreatedAt:   time.Unix(0, 0),
	})
	if err != nil {
		return zero, err
	}
	balance, err := settlementWallet.GetBalance(true)
	if err != nil {
		return zero, err
	}
	return balance.TotalProbi, nil
}
