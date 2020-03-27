package context

import "context"

// ctxKey - a type for context keys
type ctxKey int

// String - return the string associated with this context key
func (k ctxKey) String() string {
	return ctxKeyToString[k]
}

// ContextKeyFromString - get the context key from a string value
func ContextKeyFromString(k string) (ctxKey, bool) {
	v, ok := stringToCTXKey[k]
	return v, ok
}

const (
	// DatastoreCTXKey - the context key for getting the datastore
	DatastoreCTXKey ctxKey = iota + 1
	// LoggerCTXKey - the context key for getting the logger
	LoggerCTXKey
	// BatSettlementAddressCTXKey - the context key for getting the settlement address
	BatSettlementAddressCTXKey
	// ChallengeBypassServerCTXKey - the context key for getting the challenge bypass server address
	ChallengeBypassServerCTXKey
	// DatabaseMigrationsURLCTXKey - the context key for getting the database migrations url
	DatabaseMigrationsURLCTXKey
	// DatabaseURLCTXKey - the context key for getting the database connection string
	DatabaseURLCTXKey
	// DebugCTXKey - the context key for getting telling if we are in debug mode
	DebugCTXKey
	// ED25519PrivateKeyCTXKey - the context key for getting the ed25519 private key
	ED25519PrivateKeyCTXKey
	// ED25519PublicKeyCTXKey - the context key for getting the ed25519 public key
	ED25519PublicKeyCTXKey
	// EnvironmentCTXKey - the context key for getting the environment
	EnvironmentCTXKey
	// FeatureOrdersCTXKey - the context key for getting knowing if the orders feature is active
	FeatureOrdersCTXKey
	// GrantDBInstanceClassCTXKey - the context key for getting the instance class of the grant db
	GrantDBInstanceClassCTXKey
	// GrantSignatorPublicKeyCTXKey - the context key for getting the grant signator public key
	GrantSignatorPublicKeyCTXKey
	// GrantWalletCardIDCTXKey - the context key for getting the grant wallet card id
	GrantWalletCardIDCTXKey
	// GrantWalletPrivateKeyCTXKey - the context key for getting the grant wallet private key
	GrantWalletPrivateKeyCTXKey
	// GrantWalletPublicKeyCTXKey - the context key for getting the grant wallet public key
	GrantWalletPublicKeyCTXKey
	// KafkaBrokersCTXKey - the context key for getting the kafka brokers
	KafkaBrokersCTXKey
	// LedgerServerCTXKey - the context key for getting the ledger server
	LedgerServerCTXKey
	// RateAuthCTXKey - the context key for getting the auth rate
	RateAuthCTXKey
	// RatiosServerCTXKey - the context key for getting the ratios server
	RatiosServerCTXKey
	// ReputationServerCTXKey - the context key for getting the reputation server
	ReputationServerCTXKey
	// ReputationTokenCTXKey - the context key for getting the reputation token
	ReputationTokenCTXKey
	// ReadOnlyDatabaseURLCTXKey - the context key for getting the read only database url
	ReadOnlyDatabaseURLCTXKey
	// UpholdAccessTokenCTXKey - the context key for getting the uphold access token
	UpholdAccessTokenCTXKey
	// UpholdEnvironmentCTXKey - the context key for getting the uphold environment
	UpholdEnvironmentCTXKey
	// UpholdHTTPProxyCTXKey - the context key for getting the uphold http proxy
	UpholdHTTPProxyCTXKey
	// KafkaSSLCertificateLocationCTXKey - location of the kafka ssl cert
	KafkaSSLCertificateLocationCTXKey
	// KafkaSSLKeyLocationCTXKey - location of the kafka ssl key
	KafkaSSLKeyLocationCTXKey
	// KafkaSSLKeyPasswordCTXKey - value of the kafka ssl key password
	KafkaSSLKeyPasswordCTXKey
	// KafkaSSLCALocationCTXKey - location of the kafka ssl CA
	KafkaSSLCALocationCTXKey
)

var (
	stringToCTXKey = map[string]ctxKey{
		"datastore":                      DatastoreCTXKey,
		"BAT_SETTLEMENT_ADDRESS":         BatSettlementAddressCTXKey,
		"CHALLENGE_BYPASS_SERVER":        ChallengeBypassServerCTXKey,
		"DATABASE_MIGRATIONS_URL":        DatabaseMigrationsURLCTXKey,
		"DATABASE_URL":                   DatabaseURLCTXKey,
		"DEBUG":                          DebugCTXKey,
		"ED25519_PRIVATE_KEY":            ED25519PrivateKeyCTXKey,
		"ED25519_PUBLIC_KEY":             ED25519PublicKeyCTXKey,
		"ENV":                            EnvironmentCTXKey,
		"FEATURE_ORDERS":                 FeatureOrdersCTXKey,
		"GRANT_DB_INSTANCE_CLASS":        GrantDBInstanceClassCTXKey,
		"GRANT_SIGNATOR_PUBLIC_KEY":      GrantSignatorPublicKeyCTXKey,
		"GRANT_WALLET_CARD_ID":           GrantWalletCardIDCTXKey,
		"GRANT_WALLET_PRIVATE_KEY":       GrantWalletPrivateKeyCTXKey,
		"GRANT_WALLET_PUBLIC_KEY":        GrantWalletPublicKeyCTXKey,
		"KAFKA_BROKERS":                  KafkaBrokersCTXKey,
		"LEDGER_SERVER":                  LedgerServerCTXKey,
		"RATE_AUTH":                      RateAuthCTXKey,
		"RATIOS_SERVER":                  RatiosServerCTXKey,
		"REPUTATION_SERVER":              ReputationServerCTXKey,
		"REPUTATION_TOKEN":               ReputationTokenCTXKey,
		"RO_DATABASE_URL":                ReadOnlyDatabaseURLCTXKey,
		"UPHOLD_ACCESS_TOKEN":            UpholdAccessTokenCTXKey,
		"UPHOLD_ENVIRONMENT":             UpholdEnvironmentCTXKey,
		"UPHOLD_HTTP_PROXY":              UpholdHTTPProxyCTXKey,
		"KAFKA_SSL_CERTIFICATE_LOCATION": KafkaSSLCertificateLocationCTXKey,
		"KAFKA_SSL_KEY_LOCATION":         KafkaSSLKeyLocationCTXKey,
		"KAFKA_SSL_KEY_PASSWORD":         KafkaSSLKeyPasswordCTXKey,
		"KAFKA_SSL_CA_LOCATION":          KafkaSSLCALocationCTXKey,
	}

	ctxKeyToString = map[ctxKey]string{
		DatastoreCTXKey:                   "datastore",
		BatSettlementAddressCTXKey:        "BAT_SETTLEMENT_ADDRESS",
		ChallengeBypassServerCTXKey:       "CHALLENGE_BYPASS_SERVER",
		DatabaseMigrationsURLCTXKey:       "DATABASE_MIGRATIONS_URL",
		DatabaseURLCTXKey:                 "DATABASE_URL",
		DebugCTXKey:                       "DEBUG",
		ED25519PrivateKeyCTXKey:           "ED25519_PRIVATE_KEY",
		ED25519PublicKeyCTXKey:            "ED25519_PUBLIC_KEY",
		EnvironmentCTXKey:                 "ENV",
		FeatureOrdersCTXKey:               "FEATURE_ORDERS",
		GrantDBInstanceClassCTXKey:        "GRANT_DB_INSTANCE_CLASS",
		GrantSignatorPublicKeyCTXKey:      "GRANT_SIGNATOR_PUBLIC_KEY",
		GrantWalletCardIDCTXKey:           "GRANT_WALLET_CARD_ID",
		GrantWalletPrivateKeyCTXKey:       "GRANT_WALLET_PRIVATE_KEY",
		GrantWalletPublicKeyCTXKey:        "GRANT_WALLET_PUBLIC_KEY",
		KafkaBrokersCTXKey:                "KAFKA_BROKERS",
		LedgerServerCTXKey:                "LEDGER_SERVER",
		RateAuthCTXKey:                    "RATE_AUTH",
		RatiosServerCTXKey:                "RATIOS_SERVER",
		ReputationServerCTXKey:            "REPUTATION_SERVER",
		ReputationTokenCTXKey:             "REPUTATION_TOKEN",
		ReadOnlyDatabaseURLCTXKey:         "RO_DATABASE_URL",
		UpholdAccessTokenCTXKey:           "UPHOLD_ACCESS_TOKEN",
		UpholdEnvironmentCTXKey:           "UPHOLD_ENVIRONMENT",
		UpholdHTTPProxyCTXKey:             "UPHOLD_HTTP_PROXY",
		KafkaSSLCertificateLocationCTXKey: "KAFKA_SSL_CERTIFICATE_LOCATION",
		KafkaSSLKeyLocationCTXKey:         "KAFKA_SSL_KEY_LOCATION",
		KafkaSSLKeyPasswordCTXKey:         "KAFKA_SSL_KEY_PASSWORD",
		KafkaSSLCALocationCTXKey:          "KAFKA_SSL_CA_LOCATION",
	}

	confToGetFn = map[ctxKey]func(ctx context.Context) (string, error){
		BatSettlementAddressCTXKey:        BatSettlementAddressFromContext,
		ChallengeBypassServerCTXKey:       ChallengeBypassServerFromContext,
		DatabaseMigrationsURLCTXKey:       DatabaseMigrationsURLFromContext,
		DatabaseURLCTXKey:                 DatabaseURLFromContext,
		DebugCTXKey:                       DebugFromContext,
		ED25519PrivateKeyCTXKey:           ED25519PrivateKeyFromContext,
		ED25519PublicKeyCTXKey:            ED25519PublicKeyFromContext,
		EnvironmentCTXKey:                 EnvironmentFromContext,
		FeatureOrdersCTXKey:               FeatureOrdersFromContext,
		GrantDBInstanceClassCTXKey:        GrantDBInstanceClassFromContext,
		GrantSignatorPublicKeyCTXKey:      GrantSignatorPublicKeyFromContext,
		GrantWalletCardIDCTXKey:           GrantWalletCardIDFromContext,
		GrantWalletPrivateKeyCTXKey:       GrantWalletPrivateKeyFromContext,
		GrantWalletPublicKeyCTXKey:        GrantWalletPublicKeyFromContext,
		KafkaBrokersCTXKey:                KafkaBrokersFromContext,
		LedgerServerCTXKey:                LedgerServerFromContext,
		RateAuthCTXKey:                    RateAuthFromContext,
		RatiosServerCTXKey:                RatiosServerFromContext,
		ReputationServerCTXKey:            ReputationServerFromContext,
		ReputationTokenCTXKey:             ReputationTokenFromContext,
		ReadOnlyDatabaseURLCTXKey:         ReadOnlyDatabaseURLFromContext,
		UpholdAccessTokenCTXKey:           UpholdAccessTokenFromContext,
		UpholdEnvironmentCTXKey:           UpholdEnvironmentFromContext,
		UpholdHTTPProxyCTXKey:             UpholdHTTPProxyFromContext,
		KafkaSSLCertificateLocationCTXKey: KafkaSSLCertificateLocationFromContext,
		KafkaSSLKeyLocationCTXKey:         KafkaSSLKeyLocationFromContext,
		KafkaSSLKeyPasswordCTXKey:         KafkaSSLKeyPasswordFromContext,
		KafkaSSLCALocationCTXKey:          KafkaSSLCALocationFromContext,
	}
)
