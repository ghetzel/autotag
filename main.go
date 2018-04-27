package main

import (
	"os"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/sliceutil"
)

func main() {
	app := cli.NewApp()
	app.Name = `autotag`
	app.Usage = `Automatically tag media files based on filename patterns.`
	app.Version = `0.0.1`
	app.EnableBashCompletion = false

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `debug`,
			EnvVar: `LOGLEVEL`,
		},
		cli.BoolFlag{
			Name:  `dry-run, n`,
			Usage: `Don't modify files, just print what changes would be made.`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))

		log.Infof("Starting %s %s", c.App.Name, c.App.Version)
		return nil
	}

	app.Action = func(c *cli.Context) {
		if err := NewScanner().Scan(sliceutil.OrString(c.Args().First(), `.`)); err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
