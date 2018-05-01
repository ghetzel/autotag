package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/pathutil"
	"github.com/ghetzel/go-stockutil/stringutil"
)

var AutotagFile string = `autotag.list`

type Scanner struct {
	PatternFile string
	overrides   map[string]interface{}
}

func NewScanner() *Scanner {
	return &Scanner{
		overrides: make(map[string]interface{}),
	}
}

func (self *Scanner) Override(key string, value interface{}) {
	self.overrides[key] = value
}

func (self *Scanner) Scan(roots ...string) <-chan []*FileMatch {
	var count int
	var matchchan = make(chan []*FileMatch)

	go func() {
		defer close(matchchan)

		for _, relroot := range roots {
			if root, err := filepath.Abs(relroot); err == nil {
				var rules Rules
				var patternFile string

				if self.PatternFile == `` {
					patternFile = root
				} else {
					patternFile = self.PatternFile
				}

				if r, err := self.LoadRulesFromPath(patternFile); err == nil {
					rules = r
				} else {
					log.Error(err)
					return
				}

				matches := make([]*FileMatch, 0)
				var lastParent string

				if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
					count += 1

					if lastParent != `` && filepath.Dir(path) != lastParent && len(matches) > 0 {
						matchchan <- matches
						matches = nil
					}

					if err == nil && info.Mode().IsRegular() {
						if rule, m := rules.Match(path); m != nil {
							fm := &FileMatch{
								Path: path,
								Rule: rule,
								Tags: make(map[string]interface{}),
							}

							// extract values from named captures
							for c, v := range m.NamedCaptures() {
								fm.Tags[c] = stringutil.Autotype(v)
							}

							// apply overrides
							for o, k := range self.overrides {
								fm.Tags[o] = k
							}

							matches = append(matches, fm)
						} else {
							log.Debugf("SCAN: %v", path)
						}
					}

					lastParent = filepath.Dir(path)
					return nil
				}); err == nil {
					if len(matches) > 0 {
						matchchan <- matches
					}
				} else {
					log.Error(err)
					return
				}

			} else {
				log.Error(err)
				return
			}
		}

		log.Infof("Scan completed, read %d files", count)
		return
	}()

	return matchchan
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
