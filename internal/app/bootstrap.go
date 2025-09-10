package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/config"
	"api-gateway/internal/observability"
	"api-gateway/internal/server"
)

func Run() {
	observability.InitLogging()
	cfg := config.Load()
	srv := server.New(cfg)

	go func() {
		observability.Logger().Info().Str("addr", srv.Addr).Msg("http server started")
		if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			observability.Logger().Fatal().Err(err).Msg("server listen failed")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		observability.Logger().Error().Err(err).Msg("server shutdown error")
	}
	observability.Logger().Info().Msg("server stopped")
}
