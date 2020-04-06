package tagger

import "github.com/Suburbia-io/tagengine"

var Rule = tagengine.NewRule
var RuleGroup = tagengine.NewRuleGroup

type Fingerprint struct {
	Fingerprint string
	// RawTextWithSupCat can be removed when the supplemental category has been
	// removed from the fingerprints.
	RawTextWithSupCat     string
	RawText               string
	Count                 int
	BrandMatches          []tagengine.Match
	Brand                 string
	BrandConfidence       float64
	Category              string
	CategoryConfidence    float64
	UnitMeasure           string
	UnitMeasureConfidence float64
	BrandCons             string
}

type Processor interface {
	Process(*Fingerprint)
	Done()
}
