package tagger

import (
	"log"
	"sort"
)

type UnknownTagRemover struct {
	tagDB   *TagDB
	removed map[string]map[string]int
}

func NewUnknownTagRemover(tagDBPath string) *UnknownTagRemover {
	return &UnknownTagRemover{
		tagDB:   NewTagDB(tagDBPath),
		removed: map[string]map[string]int{},
	}
}

func (t *UnknownTagRemover) Process(fp *Fingerprint) {
	// Brand
	if fp.Brand != "" && !t.tagDB.TagExists("brand", fp.Brand) {
		t.put("brand", fp.Brand, fp.Count)
		fp.Brand = ""
	}

	// Category
	if fp.Category != "" && !t.tagDB.TagExists("category", fp.Category) {
		t.put("category", fp.Category, fp.Count)
		fp.Category = ""
	}

	// Unit Measure
	if fp.UnitMeasure != "" && !t.tagDB.TagExists("unit_measure", fp.UnitMeasure) {
		t.put("unit_measure", fp.UnitMeasure, fp.Count)
		fp.UnitMeasure = ""
	}
}

func (t *UnknownTagRemover) Done() {
	log.Printf("***** Tags Removed *****")
	for tagType, removed := range t.removed {
		type tagCount struct {
			Tag   string
			Count int
		}

		var tagCounts []tagCount
		for t, c := range removed {
			tagCounts = append(tagCounts, tagCount{t, c})
		}

		sort.Slice(tagCounts, func(i, j int) bool {
			return tagCounts[i].Count > tagCounts[j].Count
		})

		for _, tc := range tagCounts {
			log.Printf("%s - %9d - %s", tagType, tc.Count, tc.Tag)
		}
	}
}

func (t *UnknownTagRemover) put(tagType, tag string, count int) {
	if _, ok := t.removed[tagType]; !ok {
		t.removed[tagType] = map[string]int{}
	}
	if _, ok := t.removed[tagType][tag]; !ok {
		t.removed[tagType][tag] = 0
	}
	t.removed[tagType][tag] += count
}
