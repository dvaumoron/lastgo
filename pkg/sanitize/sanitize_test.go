package sanitize_test

import (
	"testing"

	"github.com/dvaumoron/lastgo/pkg/sanitize"
)

func TestArchivePathClean(t *testing.T) {
	t.Parallel()

	path, err := sanitize.ArchivePath("/home/test", "index.json")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if path != "/home/test/index.json" {
		t.Error("Unexpected result, get :", path)
	}
}

func TestArchivePathTainted(t *testing.T) {
	t.Parallel()

	if _, err := sanitize.ArchivePath("/home/test", "../index.json"); err == nil {
		t.Error("Should fail on tainted path")
	}
}
