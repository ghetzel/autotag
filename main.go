package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/maputil"
)

func main() {
	app := cli.NewApp()
	app.Name = `autotag`
	app.Usage = `Automatically tag media files based on filename patterns.`
	app.Version = `0.0.2`
	app.EnableBashCompletion = false

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `info`,
			EnvVar: `LOGLEVEL`,
		},
		cli.BoolFlag{
			Name:  `dry-run, n`,
			Usage: `Don't modify files, just print what changes would be made.`,
		},
		cli.BoolFlag{
			Name:  `yes, y`,
			Usage: `Assume 'yes' to any prompts.`,
		},
		cli.BoolFlag{
			Name:  `all, A`,
			Usage: `Process all matches at once instead of grouping by parent directory.`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))

		log.Infof("Starting %s %s", c.App.Name, c.App.Version)
		return nil
	}

	app.Action = func(c *cli.Context) {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		var allMatches []*FileMatch
		var roots []string

		if c.NArg() > 0 {
			roots = []string(c.Args())
		} else {
			roots = []string{`.`}
		}

		for matches := range NewScanner().Scan(roots...) {
			if c.Bool(`all`) {
				allMatches = append(allMatches, matches...)
			} else {
				processMatches(c, tw, matches)
			}
		}

		if len(allMatches) > 0 {
			processMatches(c, tw, allMatches)
		}
	}

	app.Run(os.Args)
}

func processMatches(c *cli.Context, tw *tabwriter.Writer, matches []*FileMatch) {
	if len(matches) == 0 {
		return
	}

	if c.Bool(`all`) {
		log.Noticef("%d Matches", len(matches))
	} else {
		log.Noticef("Directory: %v", filepath.Dir(matches[0].Path))
	}

	fmt.Fprintln(tw, "DISC\tTRACK\tTITLE\tARTIST\tALBUM")

	for _, match := range matches {
		p := maputil.M(match.Tags)

		fmt.Fprintf(
			tw,
			"%v\t%v\t%v\t%v\t%v\n",
			p.Int(`disc`),
			p.Int(`track`),
			p.String(`title`),
			p.String(`artist`),
			p.String(`album`),
		)
	}

	tw.Flush()

	var proceed bool

	if c.Bool(`dry-run`) {
		log.Noticef("[dry-run] Would update %d files", len(matches))
		return
	} else if c.Bool(`yes`) {
		proceed = true
	} else {
		proceed = log.Confirm("\nProceed with applying changes? (y/n): ")
	}

	if proceed {
		for _, match := range matches {
			if err := match.Apply(); err == nil {
				log.Infof("%v: updated file", match.Path)
			} else {
				log.Warningf("%v: %v", match.Path, err)
			}
		}
	}
}
