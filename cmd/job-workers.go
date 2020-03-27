package cmd

import (
	"context"
	"time"

	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	jobCadence time.Duration
	jobWorkers int
)

func init() {
	// add persistent flags for all job-workers
	// --job-cadence 5s # this indicates 5 seconds between job runs
	jobWorkersCmd.PersistentFlags().DurationVar(
		&jobCadence, "job-cadence", 5*time.Second, "This job should run within this frequency")
	// bind viper to our persistent flag
	viper.BindPFlag("job-cadence", jobWorkersCmd.Flags().Lookup("job-cadence"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("job-cadence", "JOB_CADENCE")
	// --job-workers 5 # this indicates the number of workers to run
	jobWorkersCmd.PersistentFlags().IntVar(
		&jobWorkers, "job-workers", 1, "This is the number of workers to run")
	// bind viper to our persistent flag
	viper.BindPFlag("job-workers", jobWorkersCmd.Flags().Lookup("job-workers"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("job-workers", "JOB_WORKERS")
}

var jobWorkersCmd = &cobra.Command{
	Use:   "job-workers",
	Short: "start up job workers",
	Run:   runJobWorkers,
}

func runJobWorkers(cmd *cobra.Command, args []string) {
	// setup context and logger first
	_, logger := setupLogger(ctx)
	logger.Info().Msg("starting the job workers...")
}

// perform the job runs
func jobWorker(ctx context.Context, job func(context.Context) (bool, error), duration time.Duration) {
	for {
		logger, err := appctx.LoggerFromContext(ctx)
		logger.Info().Msg("about to run job!")

		_, err = job(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Job failed!")

			sentry.CaptureMessage(err.Error())
			sentry.Flush(time.Second * 2)
		}
		// regardless if attempted or not, wait for the duration until retrying
		logger.Info().Msg("waiting for next job run!")
		<-time.After(duration)
	}
}
