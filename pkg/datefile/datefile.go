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

package datefile

import (
	"fmt"
	"os"
	"time"
)

func OutsideInterval(filePath string, interval time.Duration) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Unable to read previous date in file :", err)
		return true
	}
	previousCheck, err := time.Parse(time.DateTime, string(data))
	if err != nil {
		fmt.Println("Unable to parse previous date in file :", err)
		return true
	}
	return time.Now().Sub(previousCheck) > interval
}

func Write(filePath string) {
	if err := os.WriteFile(filePath, time.Now().AppendFormat(nil, time.DateTime), 0644); err != nil {
		fmt.Println("Unable to write date in file :", err)
	}
}
