package cmd

import (
	"context"

	"github.com/danielmichaels/shortlink-go/internal/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	zlog "github.com/danielmichaels/shortlink-go/internal/logger"
	"github.com/danielmichaels/shortlink-go/internal/server"
	"github.com/danielmichaels/shortlink-go/internal/store"
	"github.com/danielmichaels/shortlink-go/internal/templates"
	"github.com/spf13/cobra"
)

func ServeCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Args:  cobra.NoArgs,
		Short: "Start the web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.AppConfig()
			logger := zlog.NewLogger("tars", cfg)
			db, err := store.OpenDB(cfg)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to open database. exiting")
			}
			logger.Info().Msg("database connection established")
			templateCache, err := templates.NewTemplateCache()
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to create a template cache")
			}
			app := &server.Application{
				Config:   cfg,
				Logger:   logger,
				Models:   data.NewModels(db),
				Template: templateCache,
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
