/*
 *
 * Copyright 2024 lastgo authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dvaumoron/lastgo/pkg/datefile"
	"github.com/dvaumoron/lastgo/pkg/goversion"
	"github.com/dvaumoron/lastgo/pkg/htmlquery"
	sha256 "github.com/dvaumoron/lastgo/pkg/sha256"
	"github.com/dvaumoron/lastgo/pkg/uncompress"
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

		href := ""
		s.Find("a.download").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			href = getHref(s)
			return false
		})

		version := goversion.Find(href)
		if version == "" {
			return versionDesc{}, false
		}

		splitted := strings.Split(href, version)
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
			downloadURL:    absoluteURL(conf.downloadURL, href),
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

func absoluteURL(downloadURL string, href string) string {
	if strings.Contains(href, "://") {
		return href
	}

	baseURL, _ := url.Parse(downloadURL)
	baseURL.Path = ""
	return baseURL.JoinPath(href).String()
}

func install(rootPath string, desc versionDesc) error {
	response, err := http.Get(desc.downloadURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = sha256.Check(data, desc.sha256checksum); err != nil {
		return err
	}

	return uncompress.ToDir(data, desc.downloadURL, filepath.Join(rootPath, desc.version))
}
