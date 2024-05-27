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

package sha256_test

import (
	_ "embed"
	"testing"

	"github.com/dvaumoron/lastgo/pkg/sha256"
)

//go:embed testdata/hello.txt
var data []byte

func TestSha256CheckCorrect(t *testing.T) {
	t.Parallel()

	if err := sha256.Check(data, "f39dbe5c5d496487ebf99acb7edde362bf0c2e0967473cd219f3bdaa08928e9b"); err != nil {
		t.Error("Unexpected error : ", err)
	}
}

func TestSha256CheckError(t *testing.T) {
	t.Parallel()

	if err := sha256.Check(data, "dee7dbd08c2b244e15e309b028e3856d5a48e688854cd40fbf46733c1d66ebac"); err == nil {
		t.Error("Should fail on non corresponding file and checksum")
	} else if err != sha256.ErrCheck {
		t.Error("Incorrect error reported, get :", err)
	}
}
