package main

import (
	"regexp"

	"github.com/ghetzel/go-stockutil/rxutil"
)

type Rules []*Rule

func (self Rules) Match(filepath string) (*Rule, *rxutil.MatchResult) {
	for _, rule := range self {
		if match := rule.Match(filepath); match != nil {
			return rule, match
		}
	}

	return nil, nil
}

type Rule struct {
	rx *regexp.Regexp
}

func NewRule(rx *regexp.Regexp) *Rule {
	return &Rule{
		rx: rx,
	}
}

func (self *Rule) Match(filepath string) *rxutil.MatchResult {
	return rxutil.Match(self.rx, filepath)
}
