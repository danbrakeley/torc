package main

import (
	"strings"
)

func DiffSortedSlices(a, b []string) (aOnly, bOnly []string) {
	iA, iB := 0, 0
	cap := max(len(a), len(b))
	aOnly = make([]string, 0, cap)
	bOnly = make([]string, 0, cap)
	for {
		// if we've run out of strings in 'a', then the rest of the strings in 'b' are only in 'b'
		if iA >= len(a) {
			bOnly = append(bOnly, b[iB:]...)
			break
		}

		// if we've run out of strings in 'b', then the rest of the strings in 'a' are only in 'a'
		if iB >= len(b) {
			aOnly = append(aOnly, a[iA:]...)
			break
		}

		// otherwise, advance through the list looking for differences
		var cmp = strings.Compare(a[iA], b[iB])
		if cmp == 0 {
			iA++
			iB++
		} else if cmp < 0 {
			aOnly = append(aOnly, a[iA])
			iA++
		} else {
			bOnly = append(bOnly, b[iB])
			iB++
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
