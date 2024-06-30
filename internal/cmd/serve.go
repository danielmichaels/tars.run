package cmd

import (
	"context"

	"github.com/danielmichaels/shortlink-go/internal/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	zlog "github.com/danielmichaels/shortlink-go/internal/logger"
	"github.com/danielmichaels/shortlink-go/internal/server"
	"github.com/spf13/cobra"
)

func ServeCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Args:  cobra.NoArgs,
		Short: "Start the web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.AppConfig()
			logger := zlog.NewLogger(cfg.Names.AppName, cfg)
			db, err := data.OpenDB(cfg)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to open database. exiting")
			}
			logger.Info().Msg("database connection established")
			app := &server.Application{
				Config: cfg,
				Logger: logger,
				Models: data.NewModels(db),
			}

			err = app.Serve(ctx)
			if err != nil {
				app.Logger.Error().Err(err).Msg("server failed to start")
			}
			return nil
		},
	}
	return cmd
}
