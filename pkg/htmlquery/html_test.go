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
