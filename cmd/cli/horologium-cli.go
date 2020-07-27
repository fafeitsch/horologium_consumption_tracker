package main

import (
	"github.com/fafeitsch/Horologium/horologium"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.App{
		Name:                 "Horologium",
		Description:          "Horologium reads consumption files and reports the consumption as well as the generated costs on a monthly basis.",
		Authors:              []*cli.Author{{Name: "Fabian Feitsch", Email: "info@fafeitsch.de"}},
		Copyright:            "MIT License",
		Usage:                "horologium-cli [OPTIONS] DATA_FILE",
		Version:              "1.0.0",
		Commands:             nil,
		EnableBashCompletion: true,
		Action: func(context *cli.Context) error {
			filename := context.Args().Get(0)
			reader, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer func() {
				_ = reader.Close()
			}()
			series, err := horologium.LoadFromReader(reader)
			if err != nil {
				return err
			}
			stats := series.MonthlyStatistics(time.Now().AddDate(0, -6, 0), time.Now())
			horologium.MonthlyStatistics(os.Stdout, stats)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
