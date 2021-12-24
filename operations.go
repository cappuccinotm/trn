package timerange

import (
	"sort"
	"time"
)

// Intersection returns the date range, which is common for all the given ranges.
func Intersection(ranges []Range) Range {
	if len(ranges) == 0 {
		return Range{}
	}

	resRange := ranges[0]

	for _, rng := range ranges[1:] {
		resRange = resRange.Truncate(rng)
	}

	return resRange
}

// MergeOverlappingRanges looks in the ranges slice, seeks for overlapping ranges and
// merges such ranges into the one range.
func MergeOverlappingRanges(ranges []Range) []Range {
	var res []Range

	boundaries := rangesToBoundaries(ranges)
	// sorting boundaries by time
	sort.Slice(boundaries, func(i, j int) bool { return boundaries[i].tm.Before(boundaries[j].tm) })

	// add first boundary
	var rangeStartTm time.Time
	unfinishedBoundariesCnt := 0

	// skip last boundary to allow looking ahead
	for i := 0; i < len(boundaries)-1; i++ {
		boundary := boundaries[i]

		if boundary.typ == boundaryStart {
			if unfinishedBoundariesCnt == 0 {
				rangeStartTm = boundary.tm
			}
			unfinishedBoundariesCnt++
			continue
		}

		nextBoundary := boundaries[i+1]
		// if current and previous boundaries are equal - ignore them
		if boundary.tm.Equal(nextBoundary.tm) && nextBoundary.typ == boundaryStart {
			i++
			continue
		}

		unfinishedBoundariesCnt--
		// if this is an ending boundary and there is where the merged range ends...
		if unfinishedBoundariesCnt == 0 {
			res = append(res, Range{st: rangeStartTm, dur: boundary.tm.Sub(rangeStartTm)})
		}
	}

	// process the last boundary, it must be the end boundary anyway
	unfinishedBoundariesCnt--
	if unfinishedBoundariesCnt == 0 {
		res = append(res, Range{st: rangeStartTm, dur: boundaries[len(boundaries)-1].tm.Sub(rangeStartTm)})
	}

	return res
}

func rangesToBoundaries(ranges []Range) []*timeRangeBoundary {
	res := make([]*timeRangeBoundary, len(ranges)*2)
	for i, rng := range ranges {
		res[i*2] = &timeRangeBoundary{tm: rng.st, typ: boundaryStart}
		res[i*2+1] = &timeRangeBoundary{tm: rng.End(), typ: boundaryEnd}
	}
	return res
}

type boundaryType int

const (
	boundaryStart boundaryType = 0
	boundaryEnd   boundaryType = 1
)

type timeRangeBoundary struct {
	tm  time.Time
	typ boundaryType
}
