// Package main contains the CLI program for Horologium.
package main

import (
	"github.com/fafeitsch/Horologium/horologium"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	var months int
	monthsFlag := cli.IntFlag{Name: "lastMonths", Value: 6, Usage: "The number of last full months to show in the statistics (exlcuding the current month).", Destination: &months}
	app := cli.App{
		Name:                 "Horologium",
		Description:          "Horologium reads consumption files and reports the consumption as well as the generated costs on a monthly basis.",
		Authors:              []*cli.Author{{Name: "Fabian Feitsch", Email: "info@fafeitsch.de"}},
		Copyright:            "MIT License",
		Usage:                "horologium [OPTIONS] DATA_FILE",
		Version:              "1.0.0",
		Commands:             []*cli.Command{},
		EnableBashCompletion: true,
		Flags:                []cli.Flag{&monthsFlag},
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
			beforeMonths := time.Now().AddDate(0, int(-math.Abs(float64(months))), 0)
			start := horologium.CreateDate(beforeMonths.Year(), int(beforeMonths.Month()), 1)
			stats := series.MonthlyStatistics(start, time.Now())
			stats.Render(os.Stdout)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
