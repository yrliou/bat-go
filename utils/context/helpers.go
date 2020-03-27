package context

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

// LoggerFromContext - get the logger from the context
func LoggerFromContext(ctx context.Context) (zerolog.Logger, error) {
	if logger, ok := ctx.Value(LoggerCTXKey).(zerolog.Logger); ok {
		return logger, nil
	}
	return zerolog.Logger{}, ErrLoggerNotInContext
}

// RouterFromContext - get the logger from the context
func RouterFromContext(ctx context.Context) (chi.Router, error) {
	if router, ok := ctx.Value(ServiceRouterCTXKey).(chi.Router); ok {
		return router, nil
	}
	return nil, ErrServiceRouterNotInContext
}

// stringFromContext - generic string from context with custom error
func stringFromContext(ctx context.Context, k ctxKey, failureErr error) (string, error) {
	if v, ok := ctx.Value(k).(string); ok {
		return v, nil
	}
	return "", failureErr
}

// BatSettlementAddressFromContext - helper to get the uphold http proxy from the context
func BatSettlementAddressFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, BatSettlementAddressCTXKey, ErrBatSettlementAddressNotInContext)
}

// ChallengeBypassServerFromContext - helper to get the uphold http proxy from the context
func ChallengeBypassServerFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ChallengeBypassServerCTXKey, ErrChallengeBypassServerNotInContext)
}

// DatabaseMigrationsURLFromContext - helper to get the uphold http proxy from the context
func DatabaseMigrationsURLFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, DatabaseMigrationsURLCTXKey, ErrDatabaseMigrationsURLNotInContext)
}

// DatabaseURLFromContext - helper to get the uphold http proxy from the context
func DatabaseURLFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, DatabaseURLCTXKey, ErrDatabaseURLNotInContext)
}

// DebugFromContext - helper to get the uphold http proxy from the context
func DebugFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, DebugCTXKey, ErrDebugNotInContext)
}

// ED25519PrivateKeyFromContext - helper to get the uphold http proxy from the context
func ED25519PrivateKeyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ED25519PrivateKeyCTXKey, ErrED25519PrivateKeyNotInContext)
}

// ED25519PublicKeyFromContext - helper to get the uphold http proxy from the context
func ED25519PublicKeyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ED25519PublicKeyCTXKey, ErrED25519PublicKeyNotInContext)
}

// EnvironmentFromContext - helper to get the uphold http proxy from the context
func EnvironmentFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, EnvironmentCTXKey, ErrEnvironmentNotInContext)
}

// FeatureOrdersFromContext - helper to get the uphold http proxy from the context
func FeatureOrdersFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, FeatureOrdersCTXKey, ErrFeatureOrdersNotInContext)
}

// GrantDBInstanceClassFromContext - helper to get the uphold http proxy from the context
func GrantDBInstanceClassFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, GrantDBInstanceClassCTXKey, ErrGrantDBInstanceClassNotInContext)
}

// GrantSignatorPublicKeyFromContext - helper to get the uphold http proxy from the context
func GrantSignatorPublicKeyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, GrantSignatorPublicKeyCTXKey, ErrGrantSignatorPublicKeyNotInContext)
}

// GrantWalletCardIDFromContext - helper to get the uphold http proxy from the context
func GrantWalletCardIDFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, GrantWalletCardIDCTXKey, ErrGrantWalletCardIDNotInContext)
}

// GrantWalletPrivateKeyFromContext - helper to get the uphold http proxy from the context
func GrantWalletPrivateKeyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, GrantWalletPrivateKeyCTXKey, ErrGrantWalletPrivateKeyNotInContext)
}

// GrantWalletPublicKeyFromContext - helper to get the uphold http proxy from the context
func GrantWalletPublicKeyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, GrantWalletPublicKeyCTXKey, ErrGrantWalletPublicKeyNotInContext)
}

// KafkaBrokersFromContext - helper to get the uphold http proxy from the context
func KafkaBrokersFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, KafkaBrokersCTXKey, ErrKafkaBrokersNotInContext)
}

// LedgerServerFromContext - helper to get the uphold http proxy from the context
func LedgerServerFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, LedgerServerCTXKey, ErrLedgerServerNotInContext)
}

// RateAuthFromContext - helper to get the uphold http proxy from the context
func RateAuthFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, RateAuthCTXKey, ErrRateAuthNotInContext)
}

// RatiosServerFromContext - helper to get the uphold http proxy from the context
func RatiosServerFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, RatiosServerCTXKey, ErrRatiosServerNotInContext)
}

// ReputationServerFromContext - helper to get the uphold http proxy from the context
func ReputationServerFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ReputationServerCTXKey, ErrReputationServerNotInContext)
}

// ReputationTokenFromContext - helper to get the uphold http proxy from the context
func ReputationTokenFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ReputationTokenCTXKey, ErrReputationTokenNotInContext)
}

// ReadOnlyDatabaseURLFromContext - helper to get the uphold http proxy from the context
func ReadOnlyDatabaseURLFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ReadOnlyDatabaseURLCTXKey, ErrReadOnlyDatabaseURLNotInContext)
}

// UpholdAccessTokenFromContext - helper to get the uphold http proxy from the context
func UpholdAccessTokenFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, UpholdAccessTokenCTXKey, ErrUpholdAccessTokenNotInContext)
}

// UpholdEnvironmentFromContext - helper to get the uphold http proxy from the context
func UpholdEnvironmentFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, UpholdEnvironmentCTXKey, ErrUpholdEnvironmentNotInContext)
}

// UpholdHTTPProxyFromContext - helper to get the uphold http proxy from the context
func UpholdHTTPProxyFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, UpholdHTTPProxyCTXKey, ErrUpholdHTTPProxyNotInContext)
}

// KafkaSSLCertificateLocationFromContext - helper to get the uphold http proxy from the context
func KafkaSSLCertificateLocationFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, KafkaSSLCertificateLocationCTXKey, ErrKafkaSSLCertificateLocationNotInContext)
}

// KafkaSSLKeyLocationFromContext - helper to get the uphold http proxy from the context
func KafkaSSLKeyLocationFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, KafkaSSLKeyLocationCTXKey, ErrKafkaSSLKeyLocationNotInContext)
}

// KafkaSSLKeyPasswordFromContext - helper to get the uphold http proxy from the context
func KafkaSSLKeyPasswordFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, KafkaSSLKeyPasswordCTXKey, ErrKafkaSSLKeyPasswordNotInContext)
}

// KafkaSSLCALocationFromContext - helper to get the uphold http proxy from the context
func KafkaSSLCALocationFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, KafkaSSLCALocationCTXKey, ErrKafkaSSLCALocationNotInContext)
}

// ServiceAddrFromContext - helper to get the service addr from the context
func ServiceAddrFromContext(ctx context.Context) (string, error) {
	return stringFromContext(ctx, ServiceAddrCTXKey, ErrServiceAddrNotInContext)
}

// GetConfValue - get a configuration value...
// To bridge the gap between using os.Getenv directly, and now using the context values, this is the helper wrapper
// from which everything can call.  It will look in context first, and if not found fall back on the os.Getenv
func GetConfValue(ctx context.Context, k string) (string, error) {
	if key, ok := stringToCTXKey[k]; ok {
		// attempt to pull value from context
		if fn, ok := confToGetFn[key]; ok {
			v, err := fn(ctx)
			if err == nil {
				// we have something from the context
				return v, nil
			}
			if err != nil && !errors.Is(err, ErrNotInContext) {
				// real error, bail?
				return "", fmt.Errorf("failed to get conf value: %w", err)
			}
		}
	}
	// default to get environment value
	return os.Getenv(k), nil
}
