package tagger

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func writeRulesetStats(prefix string, rs *ruleSet) {
	tagCounts := writeRuleStats(prefix, rs)
	writeTagLineStats(prefix, rs, tagCounts)
}

func writeTagLineStats(prefix string, rs *ruleSet, tagCounts map[string]int) {
	tags := []string{}
	for s := range tagCounts {
		tags = append(tags, s)
	}

	sort.Slice(tags, func(i, j int) bool {
		iTag := tags[i]
		jTag := tags[j]
		iCount := rs.tagItemCounts[iTag]
		jCount := rs.tagItemCounts[jTag]
		if iCount != jCount {
			return jCount < iCount
		}
		return iTag < jTag
	})

	f, err := os.Create(prefix + "-tag-stats.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	write := func(row []string) {
		if err := w.Write(row); err != nil {
			panic(err)
		}
	}

	write([]string{"tag", "lines", "count"})

	for _, tag := range tags {
		write([]string{
			tag,
			strconv.FormatInt(int64(tagCounts[tag]), 10),
			strconv.FormatInt(int64(rs.tagItemCounts[tag]), 10),
		})
	}
}

func writeRuleStats(prefix string, rs *ruleSet) map[string]int {
	tagCounts := map[string]int{}

	f, err := os.Create(prefix + "-rule-stats.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	write := func(row []string) {
		if err := w.Write(row); err != nil {
			panic(err)
		}
	}

	write([]string{"tag", "first_count", "match_count", "includes", "excludes"})

	for _, r := range rs.ListRules() {
		write([]string{
			r.Tag,
			strconv.FormatInt(int64(r.FirstCount), 10),
			strconv.FormatInt(int64(r.MatchCount), 10),
			fmt.Sprintf("%#v", r.Includes),
			fmt.Sprintf("%#v", r.Excludes),
		})
		if _, ok := tagCounts[r.Tag]; !ok {
			tagCounts[r.Tag] = 0
		}
		tagCounts[r.Tag] += r.FirstCount
	}

	return tagCounts
}
