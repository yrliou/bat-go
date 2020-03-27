package cmd

import (
	"context"
	"net/http"
	"time"

	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave-intl/bat-go/payment"
	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/brave-intl/bat-go/utils/handlers"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/hlog"
	"github.com/spf13/cobra"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "serve a the payments microservice",
	Run:   runServePayments,
}

func runServePayments(cmd *cobra.Command, args []string) {
	// setup context and logger first
	serverCtx, logger := setupLogger(ctx)
	slogger := logger.With().
		Str("service", "payments").
		Str("addr", addr).
		Logger()
	slogger.Info().Msg("starting the payments service...")

	// our custom logger on ctx
	ctx = context.WithValue(ctx, appctx.LoggerCTXKey, slogger)

	// need to setup service involved
	paymentPG, err := payment.NewPostgres(databaseURL, false, "payment_db")
	if err != nil {
		sentry.CaptureMessage(err.Error())
		sentry.Flush(time.Second * 2)
		logger.Panic().Err(err).Msg("Must be able to init postgres connection to start")
	}

	// initialize our service
	service, err := payment.InitService(ctx, paymentPG)
	if err != nil {
		sentry.CaptureMessage(err.Error())
		sentry.Flush(time.Second * 2)
		logger.Panic().Err(err).Msg("Payment service initialization failed")
	}

	// get router from context.
	r, err := appctx.RouterFromContext(ctx)

	// logging based middlewares for service
	r.Use(
		hlog.NewHandler(slogger),
		hlog.UserAgentHandler("user_agent"),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		middleware.RequestLogger(&slogger))

	// the routes for the payment service
	r.Mount("/v1/orders", payment.Router(service))
	r.Mount("/v1/votes", payment.VoteRouter(service))
	// routes common to all services
	r.Get("/metrics", middleware.Metrics())
	r.Get("/health-check", handlers.HealthCheckHandler(version, buildTime, commit))

	logger.Info().
		Str("version", version).
		Str("commit", commit).
		Str("build_time", buildTime).
		Msg("server starting...")

	// start server
	if err := http.ListenAndServe(
		addr,
		chi.ServerBaseContext(serverCtx, r)); err != nil {
		sentry.CaptureMessage(err.Error())
		sentry.Flush(time.Second * 2)
		logger.Panic().Err(err).Msg("HTTP server start failed!")
	}

	logger.Info().
		Msg("shutting down.")
}
