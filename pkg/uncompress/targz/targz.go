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

package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/dvaumoron/lastgo/pkg/sanitize"
)

func UntarToDir(dataTarGz []byte, dirPath string) error {
	uncompressedStream, err := gzip.NewReader(bytes.NewReader(dataTarGz))
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		destPath, err := sanitize.ArchivePath(dirPath, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
			if err != nil {
				return err
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type during tar extraction : %c in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
