package main

import (
	"flag"
	"os"

	"github.com/Suburbia-io/ruletagger/tagger"
)

// Expects env variable DATASET_DASHBOARD_API_KEY to be set.
func main() {
	args := tagger.MainArgs{}
	args.APIKey = os.Getenv("DATASET_DASHBOARD_API_KEY")

	flag.StringVar(&args.APIScheme,
		"scheme", "https", "Either http or https.")
	flag.StringVar(&args.APIHost,
		"host", "", "API host name.")
	flag.StringVar(&args.DatasetID,
		"dataset-id", "", "Dataset ID.")
	flag.StringVar(&args.TagAppID,
		"app-id", "", "Tag application ID.")
	flag.StringVar(&args.InputDBPath,
		"input", "", "Sqlite database from dataset. Downloaded if not provided.")
	flag.StringVar(&args.OutputDBPath,
		"output", "", "Path to sqlite database for output.")
	flag.BoolVar(&args.Upload,
		"upload", false, "Set to upload data to dashboard.")
	flag.IntVar(&args.MinCount,
		"min-count", 0, "If > 0, ignore fingerprints with lower count.")
	flag.StringVar(&args.CPUProfilePath,
		"cpu-profile", "", "Path to CPU profile output.")

	flag.Parse()

	if args.APIKey == "" && args.InputDBPath == "" {
		panic("Set DATASET_DASHBOARD_API_KEY to download the sqlite db or use -input to use a local copy.")
	}
	if args.APIKey == "" && args.Upload {
		panic("Set DATASET_DASHBOARD_API_KEY to upload data.")
	}

	if args.APIHost == "" && args.InputDBPath == "" {
		panic("Add -host <hostname> to download the sqlite db or use -input to use a local copy.")
	}
	if args.APIHost == "" && args.Upload {
		panic("Add -host <hostname> to upload data.")
	}

	// TODO: Validate more arguments.

	tagger.Main(args)
}
