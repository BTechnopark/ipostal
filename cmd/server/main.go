package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BTechnopark/ipostal/config"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:        "Indonesial Postal Code",
		Description: "Indonesial Postal Code",
		Commands: []*cli.Command{
			DevServer(),
		},
		Action: func(ctx *cli.Context) error {
			var err error

			sdk := SetUpSdk()
			CreateApi(sdk)

			port := config.GetEnv("PORT", "3000")
			sdk.GetGinEngine().Run(fmt.Sprintf(":%s", port))

			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panicln(err)
	}
}
