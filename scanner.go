package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/pathutil"
)

var AutotagFile string = `autotag.list`

type Scanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (self *Scanner) Scan(root string) error {
	var count int

	if rules, err := self.LoadRulesFromPath(root); err == nil {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			count += 1

			if err == nil {
				if m := rules.Match(path); m != nil {
					out := ``

					for c, v := range m.NamedCaptures() {
						out += fmt.Sprintf("\n  % 12s: %v", c, v)
					}

					log.Infof("MATCH: %v%v", path, out)
				} else {
					log.Debugf("SCAN: %v", path)
				}
			}

			return nil
		}); err == nil {
			log.Infof("Scan completed, read %d files", count)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *Scanner) LoadRulesFromPath(root string) (Rules, error) {
	rules := make(Rules, 0)

	filename := filepath.Join(root, AutotagFile)

	// if the path exists
	if pathutil.FileExists(filename) {
		if file, err := os.Open(filename); err == nil {
			log.Debugf("Read pattern file: %v", filename)
			defer file.Close()

			lines := bufio.NewScanner(file)

			// parse the tagfile
			for lines.Scan() {
				line := lines.Text()
				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, `#`) || line == `` {
					continue
				}

				if rx, err := regexp.Compile(line); err == nil {
					rules = append(rules, NewRule(rx))
				} else {
					log.Warningf("%v: %v", filename, err)
				}
			}

			if err := lines.Err(); err != nil {
				log.Warningf("%v: %v", filename, err)
			}
		} else {
			return rules, err
		}
	}

	if dir := filepath.Dir(root); dir != root {
		if parentRules, err := self.LoadRulesFromPath(dir); err == nil {
			rules = append(rules, parentRules...)
		} else {
			return rules, err
		}
	}

	return rules, nil
}
