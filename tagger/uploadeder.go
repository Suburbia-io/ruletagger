package tagger

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Uploader struct {
	count     int
	total     int
	buf       *bytes.Buffer
	w         *csv.Writer
	apiScheme string
	apiHost   string
	apiKey    string
	datasetID string
	tagAppID  string
}

func NewUploader(args MainArgs) *Uploader {
	buf := &bytes.Buffer{}

	u := &Uploader{
		buf:       buf,
		w:         csv.NewWriter(buf),
		apiScheme: args.APIScheme,
		apiHost:   args.APIHost,
		apiKey:    args.APIKey,
		datasetID: args.DatasetID,
		tagAppID:  args.TagAppID,
	}

	u.reset()
	return u
}

func (u *Uploader) reset() {
	u.count = 0
	u.buf.Reset()
	u.w = csv.NewWriter(u.buf)
	u.w.Write([]string{"fingerprint", "tag_type", "tag", "confidence"})
}

func (u *Uploader) Process(fp *Fingerprint) {
	fmtFloat := func(x float64) string {
		return strconv.FormatFloat(x, 'f', -1, 64)
	}

	u.count++
	u.total++
	u.w.Write([]string{fp.Fingerprint, "brand", fp.Brand, fmtFloat(fp.BrandConfidence)})
	u.w.Write([]string{fp.Fingerprint, "category", fp.Category, fmtFloat(fp.CategoryConfidence)})
	u.w.Write([]string{fp.Fingerprint, "unit_measure", fp.UnitMeasure, fmtFloat(fp.UnitMeasureConfidence)})

	if u.count >= BUF_SIZE {
		u.uploadData()
	}
}

func (u *Uploader) Done() {
	u.uploadData()
}

func (ul *Uploader) uploadData() {
	if ul.count <= 0 {
		return
	}

	ul.w.Flush()
	defer ul.reset()

	// Construct request.
	u := &url.URL{
		Scheme: ul.apiScheme,
		Host:   ul.apiHost,
		Path:   "/admin/api/fingerprintTagUpsertCSV",
	}

	// Query params.
	q := url.Values{
		"datasetID": []string{ul.datasetID},
		"tagAppID":  []string{ul.tagAppID},
	}
	u.RawQuery = q.Encode()

	log.Printf("Uploading: %8d", ul.count)
	log.Printf("Total:     %8d", ul.total)

	req, err := http.NewRequest(http.MethodPost, u.String(), ul.buf)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+ul.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to upload data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Upload faled: [%d] %s", resp.StatusCode, resp.Status)
	}

	callResp := struct {
		OK   bool            `json:"ok"`
		Data json.RawMessage `json:"data"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&callResp); err != nil {
		log.Fatalf("Failed to decode body: %v", err)
	}
	if !callResp.OK {
		log.Fatalf("Upload failed: %s", callResp.Data)
	}
}
