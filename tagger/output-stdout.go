package tagger

import (
	"fmt"
	"strings"
)

type OutputStdout struct{}

func NewOutputStdout() *OutputStdout {
	return &OutputStdout{}
}

func (o *OutputStdout) Process(fp *Fingerprint) {
	fmt.Println(strings.Join([]string{
		fp.Fingerprint,
		o.fmtStr(fp.Brand),
		o.fmtStr(fp.UnitMeasure),
		fp.RawText},
		"\t"))
}

func (*OutputStdout) Done() {}

func (*OutputStdout) fmtStr(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
