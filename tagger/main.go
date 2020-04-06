package tagger

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const BUF_SIZE = 16384

type MainArgs struct {
	APIScheme      string // Either "http" or "https".
	APIHost        string // The Suburbia API host.
	APIKey         string // The Suburbia API key.
	DatasetID      string // The dataset ID for uploads.
	TagAppID       string // The tag app ID for uploads.
	InputDBPath    string // Empty to download from dashboard.
	OutputDBPath   string // Empty to ignore.
	Upload         bool   // If true, upload to the dashboard server.
	MinCount       int    // If > 0, ignore lines with a lower count.
	CPUProfilePath string // Where to save the CPU profile.
}

func Main(args MainArgs) {

	if args.InputDBPath == "" {
		args.InputDBPath = downloadDb(args.APIKey, args.APIHost, args.DatasetID)
	}

	if args.CPUProfilePath != "" {
		f, err := os.Create(args.CPUProfilePath)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// The order of the pipeline matters as some steps may depend on previous
	// values.
	procs := []Processor{
		NewExampleTagger(),

		// Add other taggers here.
		//NewOtherTagger(),
	}

	if args.OutputDBPath != "" {
		procs = append(procs, NewOutputSqlite(args.OutputDBPath))
	}

	procs = append(procs, NewUnknownTagRemover(args.InputDBPath))

	if args.Upload {
		procs = append(procs, NewUploader(args))
	}

	runPipeline(args.InputDBPath, args.MinCount, procs)
}

func runPipeline(dbPath string, minCount int, procs []Processor) {
	out := make(chan []*Fingerprint, 2)
	go readSqlite(dbPath, minCount, out)

	for _, proc := range procs {
		in := out
		out = make(chan []*Fingerprint, 2)
		go runProcessor(in, out, proc)
	}

	// Wait for final processor to finish.
	for _ = range out {
	}
}

func runProcessor(in, out chan []*Fingerprint, proc Processor) {
	for l := range in {
		for _, fp := range l {
			proc.Process(fp)
		}
		out <- l
	}
	proc.Done()
	close(out)
}

func readSqlite(path string, minCount int, out chan []*Fingerprint) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	query := `SELECT fingerprint,raw_text,count,brand FROM fingerprints`
	if minCount > 0 {
		query += fmt.Sprintf(` WHERE count > %d`, minCount)
	}
	query += ` ORDER BY count DESC`

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	l := make([]*Fingerprint, 0, BUF_SIZE)
	push := func() {
		out <- l
		l = make([]*Fingerprint, 0, BUF_SIZE)
	}
	i := 0

	for rows.Next() {
		i++
		if i%100000 == 0 {
			log.Printf("Read rows: %8d", i)
		}

		fp := new(Fingerprint)
		err := rows.Scan(
			&fp.Fingerprint,
			&fp.RawTextWithSupCat,
			&fp.Count,
			&fp.BrandCons)
		if err != nil {
			panic(err)
		}

		// Remove supplemental category if present.
		fp.RawText = fp.RawTextWithSupCat
		rootIdx := strings.Index(fp.RawText, " root - >")
		if rootIdx >= 0 {
			fp.RawText = fp.RawText[:rootIdx]
		}

		l = append(l, fp)
		if len(l) == BUF_SIZE {
			push()
		}
	}
	push()

	close(out)
}
