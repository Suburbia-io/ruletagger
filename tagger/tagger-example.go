package tagger

type ExampleTagger struct {
	rules *ruleSet
}

func NewExampleTagger() *ExampleTagger {
	t := &ExampleTagger{
		// Can also be an empty string.
		rules: newRuleSet("multiple-matches"),
	}

	t.initRules()

	return t
}

func (t *ExampleTagger) Process(fp *Fingerprint) {
	if fp.Brand != "" {
		return
	}

	fp.BrandMatches = t.rules.Match(fp, &fp.Brand, &fp.BrandConfidence)
}

func (t *ExampleTagger) Done() {
	writeRulesetStats("Example", t.rules)
}

func (t *ExampleTagger) initRules() {

	t.rules.Add(
		RuleGroup("example").
			Inc("example").
			Inc("ex."),
	)
}
