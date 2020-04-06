package tagger

import (
	"github.com/Suburbia-io/tagengine"
)

type ruleSet struct {
	mixedTag      string
	tagItemCounts map[string]int
	*tagengine.RuleSet
}

func newRuleSet(mixedTag string) *ruleSet {
	return &ruleSet{
		mixedTag:      mixedTag,
		tagItemCounts: map[string]int{},
		RuleSet:       tagengine.NewRuleSet(),
	}
}

func (rs *ruleSet) Match(
	fp *Fingerprint,
	tagDest *string,
	confDest *float64,
) []tagengine.Match {
	matches := rs.RuleSet.Match(fp.RawText)

	if len(matches) == 0 {
		return nil
	}

	if len(matches) > 1 {
		*tagDest = rs.mixedTag
		*confDest = 1
		return matches
	}

	tag := matches[0].Tag
	if _, ok := rs.tagItemCounts[tag]; !ok {
		rs.tagItemCounts[tag] = 0
	}
	rs.tagItemCounts[tag] += fp.Count

	*tagDest = tag
	*confDest = matches[0].Confidence
	return matches
}
