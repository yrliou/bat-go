// package cmd - command line interaction for bat-go services and processes
package cmd

import (
	"context"

	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	databaseURL                 string
	challengeBypassServer       string
	ledgerServer                string
	kafkaSSLCertificateLocation string
	kafkaSSLKeyLocation         string
	kafkaSSLKeyPassword         string
	kafkaSSLCALocation          string
	ctx                         = context.Background()
)

func setupRootCmd(cmd *cobra.Command, args []string) {
	// setup our base context for our application with these values
	ctx = context.WithValue(ctx, appctx.DatabaseURLCTXKey, databaseURL)
	ctx = context.WithValue(ctx, appctx.ChallengeBypassServerCTXKey, challengeBypassServer)
	ctx = context.WithValue(ctx, appctx.LedgerServerCTXKey, ledgerServer)
	ctx = context.WithValue(ctx, appctx.KafkaSSLCertificateLocationCTXKey, kafkaSSLCertificateLocation)
	ctx = context.WithValue(ctx, appctx.KafkaSSLKeyLocationCTXKey, kafkaSSLKeyLocation)
	ctx = context.WithValue(ctx, appctx.KafkaSSLKeyPasswordCTXKey, kafkaSSLKeyPassword)
	ctx = context.WithValue(ctx, appctx.KafkaSSLCALocationCTXKey, kafkaSSLCALocation)
}

func init() {
	// env: DATABASE_URL
	// --database-url postgres://.... # this indicates the database to connect to
	rootCmd.PersistentFlags().StringVar(
		&databaseURL, "database-url", "", "This is the main database url")
	// bind viper to our persistent flag
	viper.BindPFlag("database-url", rootCmd.Flags().Lookup("database-url"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("database-url", "DATABASE_URL")

	// env: CHALLENGE_BYPASS_SERVER
	// --challenge-bypass-url http://.... # this indicates the challenge bypass server to connect to
	rootCmd.PersistentFlags().StringVar(
		&challengeBypassServer, "challenge-bypass-server", "", "This is the challenge bypass server url")
	// bind viper to our persistent flag
	viper.BindPFlag("challenge-bypass-server", rootCmd.Flags().Lookup("challenge-bypass-server"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("challenge-bypass-server", "CHALLENGE_BYPASS_SERVER")

	// env: LEDGER_SERVER
	// --ledger-url http://.... # this indicates the ledger server to connect to
	rootCmd.PersistentFlags().StringVar(
		&ledgerServer, "ledger-server", "", "This is the ledger server url")
	// bind viper to our persistent flag
	viper.BindPFlag("ledger-server", rootCmd.Flags().Lookup("ledger-server"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("ledger-server", "LEDGER_SERVER")

	// env: KAFKA_SSL_CERTIFICATE_LOCATION
	// --kafka-ssl-certificate-location file://.... # this indicates the file location of the kafka ssl cert
	rootCmd.PersistentFlags().StringVar(
		&kafkaSSLCertificateLocation, "kafka-ssl-certificate-location", "", "location of the kafka ssl cert")
	// bind viper to our persistent flag
	viper.BindPFlag("kafka-ssl-certificate-location", rootCmd.Flags().Lookup("kafka-ssl-certificate-location"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("kafka-ssl-certificate-location", "KAFKA_SSL_CERTIFICATE_LOCATION")

	// env: KAFKA_SSL_KEY_LOCATION
	// --kafka-ssl-key-location file://.... # this indicates the file location of the kafka ssl client key
	rootCmd.PersistentFlags().StringVar(
		&kafkaSSLKeyLocation, "kafka-ssl-key-location", "", "location of the kafka ssl client cert")
	// bind viper to our persistent flag
	viper.BindPFlag("kafka-ssl-key-location", rootCmd.Flags().Lookup("kafka-ssl-key-location"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("kafka-ssl-key-location", "KAFKA_SSL_KEY_LOCATION")

	// env: KAFKA_SSL_CA_LOCATION
	// --kafka-ssl-key-location file://.... # this indicates the file location of the kafka ssl client key
	rootCmd.PersistentFlags().StringVar(
		&kafkaSSLCALocation, "kafka-ssl-ca-location", "", "location of the kafka ssl signing ca cert")
	// bind viper to our persistent flag
	viper.BindPFlag("kafka-ssl-ca-location", rootCmd.Flags().Lookup("kafka-ssl-ca-location"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("kafka-ssl-ca-location", "KAFKA_SSL_CA_LOCATION")

	// env: KAFKA_SSL_KEY_PASSWORD
	// --kafka-ssl-key-password # this indicates the kafka ssl client key password
	rootCmd.PersistentFlags().StringVar(
		&kafkaSSLKeyPassword, "kafka-ssl-key-password", "", "the kafka ssl client key password")
	// bind viper to our persistent flag
	viper.BindPFlag("kafka-ssl-key-password", rootCmd.Flags().Lookup("kafka-ssl-key-password"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("kafka-ssl-key-password", "KAFKA_SSL_KEY_PASSWORD")
}

var (
	rootCmd = &cobra.Command{
		Use:              "bat-go",
		Short:            "bat-go command entrypoint",
		Long:             "Command line tool for bat-go services and processes",
		PersistentPreRun: setupRootCmd,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
