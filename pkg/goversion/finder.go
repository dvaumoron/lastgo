package goversion

import (
	"regexp"
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

	// TODO

	return version{}
}
