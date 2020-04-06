package tagger

import (
	"database/sql"
	"os"
)

type OutputSqlite struct {
	db   *sql.DB
	tx   *sql.Tx
	stmt *sql.Stmt

	count int
}

const outputSqliteMigration = `
CREATE TABLE fingerprints(
  fingerprint TEXT NOT NULL PRIMARY KEY,
  raw_text    TEXT NOT NULL,
  count       INT  NOT NULL,

  brand            TEXT NOT NULL,
  brand_confidence REAL NOT NULL,

  category            TEXT NOT NULL,
  category_confidence REAL NOT NULL,

  unit_measure            TEXT NOT NULL,
  unit_measure_confidence REAL NOT NULL
);`

const outputSqliteInsertStmt = `INSERT INTO fingerprints(
 fingerprint,
 raw_text,
 count,
 brand,
 brand_confidence,
 category,
 category_confidence,
 unit_measure,
 unit_measure_confidence
)VALUES(?,?,?,?,?,?,?,?,?)`

func NewOutputSqlite(path string) *OutputSqlite {
	os.RemoveAll(path)
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_sync=OFF")
	if err != nil {
		panic(err)
	}

	if _, err = db.Exec(outputSqliteMigration); err != nil {
		panic(err)
	}

	o := &OutputSqlite{
		db: db,
	}

	o.commit()
	return o
}

func (o *OutputSqlite) commit() {
	if o.tx != nil {
		if err := o.tx.Commit(); err != nil {
			panic(err)
		}
	}

	tx, err := o.db.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare(outputSqliteInsertStmt)
	if err != nil {
		panic(err)
	}

	o.count = 0
	o.tx = tx
	o.stmt = stmt
}

func (o *OutputSqlite) Process(fp *Fingerprint) {
	_, err := o.stmt.Exec(
		fp.Fingerprint,
		fp.RawText,
		fp.Count,
		fp.Brand,
		fp.BrandConfidence,
		fp.Category,
		fp.CategoryConfidence,
		fp.UnitMeasure,
		fp.UnitMeasureConfidence)
	if err != nil {
		panic(err)
	}

	o.count++
	if o.count > 2048 {
		o.commit()
	}
}

func (o *OutputSqlite) Done() {
	o.commit()

	if err := o.db.Close(); err != nil {
		panic(err)
	}
}
