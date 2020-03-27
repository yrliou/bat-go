package context

import (
	"errors"
	"fmt"
)

var (
	// ErrNotInContext - base "not in context" error message from which others are derived
	ErrNotInContext = errors.New("failed to get value from context")

	// ErrLoggerNotInContext - failed to get logger from context error
	ErrLoggerNotInContext = fmt.Errorf("failed to get logger: %w", ErrNotInContext)
	// ErrDatastoreNotInContext - failed to get datastore from context error
	ErrDatastoreNotInContext = fmt.Errorf("failed to get datastore: %w", ErrNotInContext)
	// ErrBatSettlementAddressNotInContext - failed to get bat settlement address from context error
	ErrBatSettlementAddressNotInContext = fmt.Errorf("failed to get bat settlement address: %w", ErrNotInContext)
	// ErrChallengeBypassServerNotInContext - failed to get challenge bypass server from context error
	ErrChallengeBypassServerNotInContext = fmt.Errorf("failed to get challenge bypass server: %w", ErrNotInContext)
	// ErrDatabaseMigrationsURLNotInContext - failed to get database migrations url from context error
	ErrDatabaseMigrationsURLNotInContext = fmt.Errorf("failed to get database migrations url: %w", ErrNotInContext)
	// ErrDatabaseURLNotInContext - failed to get from context
	ErrDatabaseURLNotInContext = fmt.Errorf("failed to get database url: %w", ErrNotInContext)
	// ErrDebugNotInContext - failed to get from context
	ErrDebugNotInContext = fmt.Errorf("failed to get debug: %w", ErrNotInContext)
	// ErrED25519PrivateKeyNotInContext - failed to get from context
	ErrED25519PrivateKeyNotInContext = fmt.Errorf("failed to get ed25519 private key: %w", ErrNotInContext)
	// ErrED25519PublicKeyNotInContext - failed to get from context
	ErrED25519PublicKeyNotInContext = fmt.Errorf("failed to get ed25519 public key: %w", ErrNotInContext)
	// ErrEnvironmentNotInContext - failed to get from context
	ErrEnvironmentNotInContext = fmt.Errorf("failed to get environment: %w", ErrNotInContext)
	// ErrFeatureOrdersNotInContext - failed to get from context
	ErrFeatureOrdersNotInContext = fmt.Errorf("failed to get feature orders: %w", ErrNotInContext)
	// ErrGrantDBInstanceClassNotInContext - failed to get from context
	ErrGrantDBInstanceClassNotInContext = fmt.Errorf("failed to get grant db instance class: %w", ErrNotInContext)
	// ErrGrantSignatorPublicKeyNotInContext - failed to get from context
	ErrGrantSignatorPublicKeyNotInContext = fmt.Errorf("failed to get grant signator public key: %w", ErrNotInContext)
	// ErrGrantWalletCardIDNotInContext - failed to get from context
	ErrGrantWalletCardIDNotInContext = fmt.Errorf("failed to get grant wallet card id: %w", ErrNotInContext)
	// ErrGrantWalletPrivateKeyNotInContext - failed to get from context
	ErrGrantWalletPrivateKeyNotInContext = fmt.Errorf("failed to get grant wallet private key: %w", ErrNotInContext)
	// ErrGrantWalletPublicKeyNotInContext - failed to get from context
	ErrGrantWalletPublicKeyNotInContext = fmt.Errorf("failed to get grant wallet public key: %w", ErrNotInContext)
	// ErrKafkaBrokersNotInContext - failed to get from context
	ErrKafkaBrokersNotInContext = fmt.Errorf("failed to get kafka brokers: %w", ErrNotInContext)
	// ErrLedgerServerNotInContext - failed to get from context
	ErrLedgerServerNotInContext = fmt.Errorf("failed to get ledger server: %w", ErrNotInContext)
	// ErrRateAuthNotInContext - failed to get from context
	ErrRateAuthNotInContext = fmt.Errorf("failed to get rate auth: %w", ErrNotInContext)
	// ErrRatiosServerNotInContext - failed to get from context
	ErrRatiosServerNotInContext = fmt.Errorf("failed to get ratios server: %w", ErrNotInContext)
	// ErrReputationServerNotInContext - failed to get from context
	ErrReputationServerNotInContext = fmt.Errorf("failed to get reputation server: %w", ErrNotInContext)
	// ErrReputationTokenNotInContext - failed to get from context
	ErrReputationTokenNotInContext = fmt.Errorf("failed to get reputation token: %w", ErrNotInContext)
	// ErrReadOnlyDatabaseURLNotInContext - failed to get from context
	ErrReadOnlyDatabaseURLNotInContext = fmt.Errorf("failed to get read only database url: %w", ErrNotInContext)
	// ErrUpholdAccessTokenNotInContext - failed to get from context
	ErrUpholdAccessTokenNotInContext = fmt.Errorf("failed to get uphold access token: %w", ErrNotInContext)
	// ErrUpholdEnvironmentNotInContext - failed to get from context
	ErrUpholdEnvironmentNotInContext = fmt.Errorf("failed to get uphold environment: %w", ErrNotInContext)
	// ErrUpholdHTTPProxyNotInContext - failed to get from context
	ErrUpholdHTTPProxyNotInContext = fmt.Errorf("failed to get uphold http proxy: %w", ErrNotInContext)
	// ErrKafkaSSLCertificateLocationNotInContext - failed to get from context
	ErrKafkaSSLCertificateLocationNotInContext = fmt.Errorf("failed to get kafka ssl cert location: %w", ErrNotInContext)
	// ErrKafkaSSLKeyLocationNotInContext - failed to get from context
	ErrKafkaSSLKeyLocationNotInContext = fmt.Errorf("failed to get kafka ssl key location: %w", ErrNotInContext)
	// ErrKafkaSSLKeyPasswordNotInContext - failed to get from context
	ErrKafkaSSLKeyPasswordNotInContext = fmt.Errorf("failed to get kafka ssl key password: %w", ErrNotInContext)
	// ErrKafkaSSLCALocationNotInContext - failed to get from context
	ErrKafkaSSLCALocationNotInContext = fmt.Errorf("failed to get kafka ssl ca location: %w", ErrNotInContext)
	// ErrServiceAddrNotInContext - failed to get from context
	ErrServiceAddrNotInContext = fmt.Errorf("failed to get service address: %w", ErrNotInContext)
	// ErrServiceRouterNotInContext - failed to get from context
	ErrServiceRouterNotInContext = fmt.Errorf("failed to get service router: %w", ErrNotInContext)
)
