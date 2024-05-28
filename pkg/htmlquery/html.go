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
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func BasicSelectionExtractor(part string) func(*goquery.Selection) string {
	if part == "#text" {
		return selectionTextExtractor
	}

	return func(s *goquery.Selection) string {
		attr, _ := s.Attr(part)

		return strings.TrimSpace(attr)
	}
}

func NonEmptyFilter(basicExtractor func(*goquery.Selection) string) func(*goquery.Selection) (string, bool, bool) {
	return func(s *goquery.Selection) (string, bool, bool) {
		str := basicExtractor(s)
		return str, str != "", true
	}
}

func Request[T any](callURL string, selector string, extractor func(*goquery.Selection) (T, bool, bool)) ([]T, error) {
	response, err := http.Get(callURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return extractList(response.Body, selector, extractor)
}

func extractList[T any](reader io.Reader, selector string, extractor func(*goquery.Selection) (T, bool, bool)) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var extracteds []T
	doc.Find(selector).EachWithBreak(func(_ int, s *goquery.Selection) bool {
		extracted, ok, next := extractor(s)
		if ok {
			extracteds = append(extracteds, extracted)
		}
		return next
	})

	return extracteds, nil
}

func selectionTextExtractor(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}
