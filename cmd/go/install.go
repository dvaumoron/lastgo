package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dvaumoron/lastgo/pkg/datefile"
	"github.com/dvaumoron/lastgo/pkg/goversion"
	"github.com/dvaumoron/lastgo/pkg/htmlquery"
	"golang.org/x/net/html"
)

const archiveName = "archive"

type versionDesc struct {
	version        string
	downloadURL    string
	sha256checksum string
}

func getLastVersion(conf config) versionDesc {
	var builder strings.Builder
	builder.WriteByte('.')
	builder.WriteString(runtime.GOOS)
	builder.WriteByte('-')
	builder.WriteString(runtime.GOARCH)
	os_arch := builder.String()

	getHref := htmlquery.BasicSelectionExtractor("href")
	getInnerText := htmlquery.BasicSelectionExtractor("#text")

	found := false
	versions, err := htmlquery.Request(conf.downloadURL, "div.toggleVisible div.expanded tbody tr", func(s *goquery.Selection) (versionDesc, bool) {
		if found {
			return versionDesc{}, false
		}

		downloadURL := ""
		s.Find("a.download").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			downloadURL = getHref(s)
			return false
		})

		version := goversion.Find(downloadURL)
		if version == "" {
			return versionDesc{}, false
		}

		splitted := strings.Split(downloadURL, version)
		if len(splitted) < 2 || !strings.HasPrefix(splitted[1], os_arch) {
			return versionDesc{}, false
		}

		sha256checksum := ""
		s.Find("tt").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			sha256checksum = getInnerText(s)
			return false
		})

		found = true
		desc := versionDesc{
			version:        version,
			downloadURL:    downloadURL,
			sha256checksum: sha256checksum,
		}

		return desc, true
	})

	if err != nil || len(versions) == 0 {
		fmt.Println("Unable to retrieve last version :", err)
		os.Exit(1)
	}

	datefile.Write(conf.dateFilePath)

	return versions[0]
}

func install(installPath string, desc versionDesc) error {

	// TODO

	return os.MkdirAll(filepath.Join(installPath, desc.version), 755)
}

func notEqualInnerText(node *html.Node, value string) bool {
	return !strings.EqualFold(strings.TrimSpace(node.FirstChild.Data), value)
}
