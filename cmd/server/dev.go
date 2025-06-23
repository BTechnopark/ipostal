package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BTechnopark/ipostal/config"
	"github.com/urfave/cli/v2"
)

func DevServer() *cli.Command {
	return &cli.Command{
		Name:        "Indonesian Postal Code Development Server",
		Description: "Indonesian Postal Code Development Server",
		Aliases:     []string{"dev"},
		Before: func(ctx *cli.Context) error {
			envs := map[string]string{
				"DEV_MODE":       "1",
				"CACHE_DURATION": "1m",
			}
			for key, value := range envs {
				os.Setenv(key, value)
				slog.Info("Set Environment", slog.String(key, value))
			}

			return nil
		},
		Action: func(ctx *cli.Context) error {
			var err error

			sdk := SetUpSdk()
			err = CreateApi(sdk)
			if err != nil {
				return err
			}

			port := config.GetEnv("PORT", "8000")
			sdk.GetGinEngine().Run(fmt.Sprintf("localhost:%s", port))

			return err
		},
	}
}
