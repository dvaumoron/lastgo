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

package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"

	"github.com/dvaumoron/lastgo/pkg/sanitize"
)

func UnzipToDir(dataZip []byte, dirPath string) error {
	dataReader := bytes.NewReader(dataZip)
	zipReader, err := zip.NewReader(dataReader, int64(len(dataZip)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath); err != nil {
			return err
		}
	}

	return nil
}

// a separate function allows deferred Close to execute earlier.
func copyZipFileToDir(zipFile *zip.File, dirPath string) error {
	destPath, err := sanitize.ArchivePath(dirPath, zipFile.Name)
	if err != nil {
		return err
	}

	if destPath[len(destPath)-1] == '/' {
		// trailing slash indicates a directory
		return os.MkdirAll(destPath, 0755)
	}

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, zipFile.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	return err
}
