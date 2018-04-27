package main

import (
	"regexp"

	"github.com/ghetzel/go-stockutil/rxutil"
)

type Rules []*Rule

func (self Rules) Match(filepath string) *rxutil.MatchResult {
	for _, rule := range self {
		if match := rule.Match(filepath); match != nil {
			return match
		}
	}

	return nil
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
