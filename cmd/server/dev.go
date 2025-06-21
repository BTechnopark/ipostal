package main

import (
	"fmt"
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
				"DEV_MODE": "1",
			}
			for key, value := range envs {
				os.Setenv(key, value)
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			var err error

			sdk := SetUpSdk()
			CreateApi(sdk)

			port := config.GetEnv("PORT", "3000")
			sdk.GetGinEngine().Run(fmt.Sprintf("localhost:%s", port))

			return err
		},
	}
}
