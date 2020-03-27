package cmd

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/brave-intl/bat-go/middleware"
	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	addr string
)

func init() {
	// add persistent flags for all serve subcommands
	// --addr=":8080" # the address to serve
	serveCmd.PersistentFlags().StringVar(
		&addr, "addr", ":3333", "The address to listen/serve on")
	// bind viper to our persistent flag
	viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
	// bind the environment variable incase it comes in that way
	viper.BindEnv("addr", "ADDR")
}

func setupServeCmd(cmd *cobra.Command, args []string) {
	// make sure to setup the root
	setupRootCmd(cmd, args)
	// setup our base context for our application with these values
	ctx = context.WithValue(ctx, appctx.ServiceAddrCTXKey, addr)
	// setup router with common middlewares for all microservices
	r := chi.NewRouter().With()
	r.Use(
		chiware.RequestID,
		chiware.RealIP,
		chiware.Heartbeat("/"),
		chiware.Timeout(60*time.Second),
		middleware.BearerToken,
		middleware.RateLimiter,
		middleware.RequestIDTransfer)

	ctx = context.WithValue(ctx, appctx.ServiceRouterCTXKey, r)
	// require govalidator
	govalidator.SetFieldsRequiredByDefault(true)
}

var serveCmd = &cobra.Command{
	Use:              "serve",
	Short:            "serve a particular microservice",
	Run:              runServe,
	PersistentPreRun: setupServeCmd,
}

func runServe(cmd *cobra.Command, args []string) {
	// setup context and logger first
	_, logger := setupLogger(ctx)
	logger.Info().Msg("starting the service...")
}
