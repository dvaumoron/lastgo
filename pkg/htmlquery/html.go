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

func NonEmptyFilter(basicExtractor func(*goquery.Selection) string) func(*goquery.Selection) (string, bool) {
	return func(s *goquery.Selection) (string, bool) {
		str := basicExtractor(s)
		return str, str != ""
	}
}

func Request[T any](callURL string, selector string, extractor func(*goquery.Selection) (T, bool)) ([]T, error) {
	response, err := http.Get(callURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return extractList(response.Body, selector, extractor)
}

func extractList[T any](reader io.Reader, selector string, extractor func(*goquery.Selection) (T, bool)) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var extracteds []T
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if extracted, ok := extractor(s); ok {
			extracteds = append(extracteds, extracted)
		}
	})

	return extracteds, nil
}

func selectionTextExtractor(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}
