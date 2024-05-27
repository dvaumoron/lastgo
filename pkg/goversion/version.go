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

package goversion

import (
	"regexp"
	"strconv"
	"strings"
)

var versionRegexp *regexp.Regexp

type version struct {
	major int
	minor int
	patch int
}

func init() {
	versionRegexp = regexp.MustCompilePOSIX("(go[0-9]+(\\.[0-9]+)*)")
}

// return a go version (like "go1.22.3").
func Find(versionStr string) string {
	return versionRegexp.FindString(versionStr)
}

func Last(versions []string) string {
	last, lastStr := version{}, ""

	for _, version := range versions {
		if parsed := parse(version); less(last, parsed) {
			last, lastStr = parsed, version
		}
	}

	return lastStr
}

func Less(v1 string, v2 string) bool {
	return less(parse(v1), parse(v2))
}

func less(v1 version, v2 version) bool {
	if v1.major < v2.major {
		return true
	}
	if v1.major > v2.major {
		return false
	}
	if v1.minor < v2.minor {
		return true
	}
	if v1.minor > v2.minor {
		return false
	}
	return v1.patch < v2.patch
}

func parse(versionStr string) version {
	if len(versionStr) < 2 || versionStr[0] != 'g' || versionStr[1] != 'o' {
		return version{}
	}

	var err error
	var major, minor, patch int
	switch splitted := strings.Split(versionStr[2:], "."); len(splitted) {
	default:
		fallthrough
	case 3:
		patch, err = strconv.Atoi(splitted[2])
		if err != nil {
			return version{}
		}
		fallthrough
	case 2:
		minor, err = strconv.Atoi(splitted[1])
		if err != nil {
			return version{}
		}
		fallthrough
	case 1:
		major, err = strconv.Atoi(splitted[0])
		if err != nil {
			return version{}
		}
	case 0:
		return version{}
	}

	return version{
		major: major,
		minor: minor,
		patch: patch,
	}
}
