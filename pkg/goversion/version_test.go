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

func TestLessTrue(t *testing.T) {
	t.Parallel()

	if !goversion.Less("go1.21.10", "go1.22.3") {
		t.Error("Unexpected result, get false")
	}
}

func TestLessFalse(t *testing.T) {
	t.Parallel()

	if goversion.Less("go2", "go1.22.3") {
		t.Error("Unexpected result, get true")
	}
}

func TestLast(t *testing.T) {
	t.Parallel()

	if version := goversion.Last([]string{"go1.0.1", "go1.22.3", "go1.21", ""}); version != "go1.22.3" {
		t.Error("Unexpected result, get :", version)
	}
}
