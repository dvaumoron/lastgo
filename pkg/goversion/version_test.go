package goversion_test

import (
	"testing"

	"github.com/dvaumoron/lastgo/pkg/goversion"
)

func TestFindEmpty(t *testing.T) {
	t.Parallel()

	if version := goversion.Find("index.json"); version != "" {
		t.Error("Should not find a version, get :", version)
	}
}

func TestFindURL(t *testing.T) {
	t.Parallel()

	if version := goversion.Find("https://go.dev/dl/go1.22.3.darwin-arm64.tar.gz"); version != "go1.22.3" {
		t.Error("Unexpected result, get :", version)
	}
}
