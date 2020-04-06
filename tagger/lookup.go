package tagger

import "database/sql"

type TagDB struct {
	// Map from tag_type -> map from tag -> empty struct.
	m map[string]map[string]struct{}
}

func NewTagDB(dbPath string) *TagDB {
	m := map[string]map[string]struct{}{}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	tagTypeRows, err := db.Query(`SELECT tag_type_id,tag_type FROM tag_types`)
	if err != nil {
		panic(err)
	}

	for tagTypeRows.Next() {
		var tagTypeID, tagType string
		if err := tagTypeRows.Scan(&tagTypeID, &tagType); err != nil {
			panic(err)
		}

		m[tagType] = map[string]struct{}{}

		tagRows, err := db.Query(
			`SELECT tag FROM tags WHERE tag_type_id=?`, tagTypeID)
		if err != nil {
			panic(err)
		}
		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				panic(err)
			}
			m[tagType][tag] = struct{}{}
		}
	}

	return &TagDB{m: m}
}

func (tdb *TagDB) TagExists(tagType, tag string) bool {
	m, ok := tdb.m[tagType]
	if !ok {
		return false
	}
	_, ok = m[tag]
	return ok
}
