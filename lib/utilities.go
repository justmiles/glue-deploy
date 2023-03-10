package lib

import (
	"fmt"
	"regexp"
)

func isUpdatedVersion(txt, ver string) (bool, string) {
	const regexSemver = `(v?([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?)|(latest)|(HEAD[-_]SNAPSHOT)`
	const regexVPrefx = `^(v)`
	regexSemverR, _ := regexp.Compile(regexSemver)
	regexVPrefxR, _ := regexp.Compile(regexVPrefx)
	// Strip the v prefix out of the target version string. We'll add it back later if it's required.
	if regexVPrefxR.MatchString(ver) {
		ver = regexVPrefxR.ReplaceAllString(ver, "")
	}

	// Check if this version matches
	if regexSemverR.MatchString(txt) {
		prev := regexSemverR.FindString(txt)
		// make sure we keep the prefixed `v`
		if regexVPrefxR.MatchString(prev) && !regexVPrefxR.MatchString(ver) {
			ver = fmt.Sprintf("v%s", ver)
		}

		// ignore sets where there are not changes
		if prev == ver {
			return false, ""
		}

		return true, regexSemverR.ReplaceAllString(txt, ver)
	}
	return false, ""
}
