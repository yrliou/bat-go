package cmd

import (
	"context"
	"time"

	"github.com/brave-intl/bat-go/payment"
	appctx "github.com/brave-intl/bat-go/utils/context"
	srv "github.com/brave-intl/bat-go/utils/service"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

func init() {
	jobWorkersCmd.AddCommand(voteDrainJobCmd)
}

var voteDrainJobCmd = &cobra.Command{
	Use:   "vote-drain-job",
	Short: "start up the vote drain job",
	Run:   runVoteDrainJob,
}

func runVoteDrainJob(cmd *cobra.Command, args []string) {

	// setup context and logger first
	_, logger := setupLogger(ctx)
	logger.Info().Msg("starting the vote drain job workers...")
	logger.Info().Msg(databaseURL)

	// need to setup service involved
	paymentPG, err := payment.NewPostgres(databaseURL, false, "payment_db")
	if err != nil {
		sentry.CaptureMessage(err.Error())
		sentry.Flush(time.Second * 2)
		logger.Panic().Err(err).Msg("Must be able to init postgres connection to start")
	}
	service, err := payment.InitService(ctx, paymentPG)
	if err != nil {
		sentry.CaptureMessage(err.Error())
		sentry.Flush(time.Second * 2)
		logger.Panic().Err(err).Msg("Payment service initialization failed")
	}

	// create a job to run
	var job = srv.Job{
		Func:    service.RunNextVoteDrainJob,
		Cadence: jobCadence,
		Workers: jobWorkers,
	}

	// setup context
	slogger := logger.With().
		Str("job_name", "vote-drain-worker").
		Int("workers", jobWorkers).
		Dur("cadence", jobCadence).
		Logger()

	slogger.Info().Msg("kicking off workers!")

	// run the job
	for i := 0; i < job.Workers; i++ {
		slogger = slogger.With().
			Int("worker_number", i).
			Logger()

		ctx = context.WithValue(ctx, appctx.LoggerCTXKey, slogger)
		slogger.Info().Msg("setup logger in context go run!")

		jobWorker(ctx, job.Func, job.Cadence)
	}
}
