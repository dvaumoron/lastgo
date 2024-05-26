package htmlquery

import (
	"bytes"
	_ "embed"
	"slices"
	"testing"
)

//go:embed testdata/dl.html
var dlData []byte

func TestExtractAttrs(t *testing.T) {
	t.Parallel()

	artifactoryReader := bytes.NewReader(dlData)

	extractor := NonEmptyFilter(BasicSelectionExtractor("href"))
	extracted, err := extractList(artifactoryReader, "a.downloadBox", extractor)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if !slices.Equal(extracted, []string{
		"https://go.dev/dl/go1.22.3.windows-amd64.msi", "https://go.dev/dl/go1.22.3.darwin-arm64.pkg", "https://go.dev/dl/go1.22.3.darwin-amd64.pkg",
		"https://go.dev/dl/go1.22.3.linux-amd64.tar.gz", "https://go.dev/dl/go1.22.3.src.tar.gz",
	}) {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestExtractTexts(t *testing.T) {
	t.Parallel()

	artifactoryReader := bytes.NewReader(dlData)

	extractor := NonEmptyFilter(BasicSelectionExtractor("#text"))
	extracted, err := extractList(artifactoryReader, "div.toggleVisible div.expanded span", extractor)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if !slices.Equal(extracted, []string{"go1.22.3", "go1.21.10"}) {
		t.Error("Unmatching results, get :", extracted)
	}
}
