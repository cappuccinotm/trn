// Package store provides service data structures and methods to operate them.
// All times, returned by methods in this package are in UTC.
package store

import (
	"fmt"
	"sort"
	"time"
)

const defaultRangeFmt = "2006-01-02 15:04:05.999999999 -0700 MST"

// DateRange represents time slot with its own start and end time boundaries
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// GoString implements fmt.GoStringer to use DateRange in %#v formats
func (r DateRange) GoString() string { return r.UTC().Format(defaultRangeFmt) }

// String implements fmt.Stringer to print and log DateRange properly
func (r DateRange) String() string { return r.UTC().Format(defaultRangeFmt) }

// UTC returns the date range with boundaries in UTC.
func (r DateRange) UTC() DateRange { return r.In(time.UTC) }

// Format returns the string representation of the time range with the given format.
func (r DateRange) Format(layout string) string {
	return fmt.Sprintf("[%s, %s]", r.Start.Format(layout), r.End.Format(layout))
}

// Split the date range into smaller ranges, starting from the given offset,
// with fixed duration and with the given interval between the *end* of the
// one range and *start* of next range.
func (r DateRange) Split(offset time.Duration, duration time.Duration, interval time.Duration) []DateRange {
	if duration == 0 {
		panic("cannot split with zero duration")
	}
	return r.Stratify(offset, duration, duration+interval)
}

// Stratify the date range into smaller ranges, starting from the given offset,
// with fixed duration and with the given interval between the *starts* of the
// resulting ranges.
func (r DateRange) Stratify(offset time.Duration, duration time.Duration, interval time.Duration) []DateRange {
	if interval == 0 || duration == 0 {
		panic("cannot stratify with zero duration or zero interval")
	}

	var res []DateRange
	rangeStart := r.Start.Add(offset)

	for r.End.Sub(rangeStart.Add(duration)) >= 0 {
		res = append(res, DateRange{Start: rangeStart, End: rangeStart.Add(duration)})
		rangeStart = rangeStart.Add(interval)
	}

	return res
}

// Duration returns the duration of the date range
func (r DateRange) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// Contains returns true if the other date range is within this date range.
func (r DateRange) Contains(other DateRange) bool {
	if (r.Start.Before(other.Start) || r.Start.Equal(other.Start)) &&
		(r.End.After(other.End) || r.End.Equal(other.End)) {
		return true
	}
	return false
}

// Empty returns true if the date range is empty.
func (r DateRange) Empty() bool {
	return r.Start.IsZero() && r.End.IsZero()
}

// Truncate returns the date range bounded to the *bounds*, i.e. it cuts
// the start and the end of *r* to fit into the *bounds*.
func (r DateRange) Truncate(bounds DateRange) DateRange {
	switch {
	case r.Start.Before(bounds.Start) && r.End.Before(bounds.Start):
		// -XXX-----
		// -----YYY-
		return DateRange{}
	case r.Start.After(bounds.End) && r.End.After(bounds.End):
		// -----XXX-
		// -YYY-----
		return DateRange{}
	case r.Contains(bounds):
		// -XXXXXXX-
		// ---YYY---
		return bounds
	case bounds.Contains(r):
		// ---XXX---
		// -YYYYYYY-
		return r
	case r.Start.Before(bounds.Start) && r.End.Before(bounds.End):
		// ---XXX---
		// ----YYY--
		return DateRange{Start: bounds.Start, End: r.End}
	case r.Start.After(bounds.Start) && r.End.After(bounds.End):
		// ---XXX---
		// --YYY----
		return DateRange{Start: r.Start, End: bounds.End}
	default:
		return DateRange{}
	}
}

// Copy the given DateRange.
func (r DateRange) Copy() DateRange {
	return DateRange{Start: r.Start, End: r.End}
}

// In returns the date range with boundaries in the provided location's time zone.
func (r DateRange) In(loc *time.Location) DateRange {
	return DateRange{Start: r.Start.In(loc), End: r.End.In(loc)}
}

// DateRangesIn converts time zones of the provided date ranges into provided time zone.
func DateRangesIn(rngs []DateRange, loc *time.Location) []DateRange {
	res := make([]DateRange, len(rngs))
	for i := range rngs {
		res[i] = rngs[i].In(loc)
	}
	return res
}

// DateRangesToUTC converts time zones of the provided date ranges into UTC time zone.
func DateRangesToUTC(rngs []DateRange) []DateRange {
	return DateRangesIn(rngs, time.UTC)
}

// SplitToRangesPerDay splits the multiday ranges into smaller ranges, such that
// the resulting list of ranges will never took more than a millisecond from the
// next day.
//
// Requirements for correct working:
// - ranges must be distinct (there must not be any overlapping ranges or ranges with equal start/end boundaries)
// - ranges must sorted by the start date
//
// The inner last ranges in a day will always end at the next day at 00:00
func SplitToRangesPerDay(ranges []DateRange) map[Date][]DateRange {
	res := map[Date][]DateRange{}

	for _, rng := range ranges {
		startTime := rng.Start
		startDate := DateFromTime(startTime)
		endDate := DateFromTime(rng.End)

		for startDate.Before(endDate) {
			dayEnd := startDate.Add(0, 0, 1).Time(NewClock(24, 0, 0, 0, startTime.Location()))
			res[startDate] = append(res[startDate], DateRange{Start: startTime, End: dayEnd})

			startTime = dayEnd
			startDate = DateFromTime(startTime.Add(1 * time.Nanosecond))
		}

		if rng.End.Sub(startTime) > 1*time.Nanosecond {
			res[startDate] = append(res[startDate], DateRange{Start: startTime, End: rng.End})
		}
	}

	return res
}

// FlipDateRanges within the given period.
//
// Requirements for correct working:
// - all ranges must be within the given time period
// - ranges must be distinct (there must not be any overlapping ranges or ranges with equal start/end boundaries)
// - ranges must sorted by the start date
//
// The boundaries of the given ranges are considered to be inclusive, means
// that the flipped ranges will start or end at the exact nanosecond where
// the boundary from the input starts or ends.
//
// Complexity: O(n)
func (r DateRange) FlipDateRanges(ranges []DateRange) []DateRange {
	var res []DateRange

	// if the list of ranges is empty - just return the whole period
	if len(ranges) == 0 {
		return []DateRange{r}
	}

	// add the gap between the start of the period and start of the first range
	if !r.Start.Equal(ranges[0].Start) {
		res = append(res, DateRange{Start: r.Start, End: ranges[0].Start})
	}

	// skip first range
	for i := 1; i < len(ranges); i++ {
		res = append(res, DateRange{Start: ranges[i-1].End, End: ranges[i].Start})
	}

	// add the gap between the end of the last range and end of the period
	if !r.End.Equal(ranges[len(ranges)-1].End) {
		res = append(res, DateRange{Start: ranges[len(ranges)-1].End, End: r.End})
	}

	return res
}

// MergeOverlappingRanges looks in the ranges slice, seeks for overlapping ranges and
// merges such ranges into the one range.
// Complexity: O(n * log(n))
func MergeOverlappingRanges(ranges []DateRange) []DateRange {
	var res []DateRange

	boundaries := rangesToBoundaries(ranges)
	// sorting boundaries by time
	sort.Slice(boundaries, func(i, j int) bool { return boundaries[i].tm.Before(boundaries[j].tm) })

	// add first boundary
	var rangeStartTm time.Time
	unfinishedBoundariesCnt := 0

	// skip last boundary to allow to look ahead
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
			res = append(res, DateRange{Start: rangeStartTm, End: boundary.tm})
		}
	}

	// process the last boundary, it must be the end boundary anyway
	unfinishedBoundariesCnt--
	if unfinishedBoundariesCnt == 0 {
		res = append(res, DateRange{Start: rangeStartTm, End: boundaries[len(boundaries)-1].tm})
	}

	return res
}

// Intersection returns the intersections between the date ranges.
func Intersection(ranges []DateRange) DateRange {
	if len(ranges) < 1 {
		return DateRange{}
	}

	resRange := ranges[0]

	for _, rng := range ranges[1:] {
		resRange = resRange.Truncate(rng)
	}

	return resRange
}

// SortRanges sorts the given ranges by the start time.
func SortRanges(ranges []DateRange) []DateRange {
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].Start.Before(ranges[j].Start) })
	return ranges
}

func rangesToBoundaries(ranges []DateRange) []*timeRangeBoundary {
	res := make([]*timeRangeBoundary, len(ranges)*2)
	for i, rng := range ranges {
		res[i*2] = &timeRangeBoundary{tm: rng.Start, typ: boundaryStart}
		res[i*2+1] = &timeRangeBoundary{tm: rng.End, typ: boundaryEnd}
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
