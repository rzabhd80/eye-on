package main

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	app := cli.App{
		Name:  "Eye-on Challenge",
		Usage: "exchange and vendor agnostic transaction manager",
		Commands: []*cli.Command{
			{Name: "api", Usage: "setup api server", Action: func(ctx *cli.Context) error {
				err := apiService(ctx, logger)
				if err != nil {
					logger.Fatal("failed to setup api", zap.Error(err))
					return err
				}
				return nil
			}},
		},
	}
	er := app.Run(os.Args)
	if er != nil {
		logger.Error(er.Error())
	}
}
