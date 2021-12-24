// Package store provides service data structures and methods to operate them.
package store

import (
	"fmt"
	"sort"
	"time"
)

const defaultRangeFmt = "2006-01-02 15:04:05.999999999 -0700 MST"

// Option is an adapter over date ranges.
type Option func(r *DateRange)

// In sets the time range in the given location.
func In(loc *time.Location) Option {
	return func(r *DateRange) { r.st = r.st.In(loc) }
}

// Range returns the new DateRange in the given time bounds.
func Range(start, end time.Time, opts ...Option) DateRange {
	if start.After(end) {
		panic("start is after the end")
	}

	res := DateRange{st: start, dur: end.Sub(start)}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

// DateRange represents time slot with its own start and end time boundaries
type DateRange struct {
	st  time.Time
	dur time.Duration
}

// String implements fmt.Stringer to print and log DateRange properly
func (r DateRange) String() string { return r.UTC().Format(defaultRangeFmt) }

// UTC returns the date range with boundaries in UTC.
func (r DateRange) UTC() DateRange { return r.In(time.UTC) }

// Duration returns the duration of the date range.
func (r DateRange) Duration() time.Duration { return r.dur }

// Start returns the start time of the date range.
func (r DateRange) Start() time.Time { return r.st }

// End returns the end time of the date range.
func (r DateRange) End() time.Time { return r.st.Add(r.dur) }

// Format returns the string representation of the time range with the given format.
func (r DateRange) Format(layout string) string {
	return fmt.Sprintf("[%s, %s]", r.st.Format(layout), r.End().Format(layout))
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
	rangeStart := r.st.Add(offset)

	for r.End().Sub(rangeStart.Add(duration)) >= 0 {
		res = append(res, DateRange{st: rangeStart, dur: duration})
		rangeStart = rangeStart.Add(interval)
	}

	return res
}

// Contains returns true if the other date range is within this date range.
func (r DateRange) Contains(other DateRange) bool {
	if (r.st.Before(other.st) || r.st.Equal(other.st)) &&
		(r.End().After(other.End()) || r.End().Equal(other.End())) {
		return true
	}
	return false
}

// Empty returns true if the date range is empty.
func (r DateRange) Empty() bool { return r.st.IsZero() && r.dur == 0 }

// Truncate returns the date range bounded to the *bounds*, i.e. it cuts
// the start and the end of *r* to fit into the *bounds*.
func (r DateRange) Truncate(bounds DateRange) DateRange {
	switch {
	case r.st.Before(bounds.st) && r.End().Before(bounds.st):
		// -XXX-----
		// -----YYY-
		return DateRange{}
	case r.st.After(bounds.End()) && r.End().After(bounds.End()):
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
	case r.st.Before(bounds.st) && r.End().Before(bounds.End()):
		// ---XXX---
		// ----YYY--
		return DateRange{st: bounds.st, dur: r.End().Sub(bounds.st)}
	case r.st.After(bounds.st) && r.End().After(bounds.End()):
		// ---XXX---
		// --YYY----
		return DateRange{st: r.st, dur: bounds.End().Sub(r.st)}
	default:
		return DateRange{}
	}
}

// In returns the date range with boundaries in the provided location's time zone.
func (r DateRange) In(loc *time.Location) DateRange {
	return DateRange{st: r.st.In(loc), dur: r.dur}
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

// FlipDateRanges within the given period.
//
// Requirements for correct working:
// - all ranges must be within the given time period
// - ranges must be distinct (there must not be any overlapping ranges or ranges with equal start/end boundaries)
// - ranges must be sorted by the start date
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
	if !r.st.Equal(ranges[0].st) {
		res = append(res, DateRange{st: r.st, dur: ranges[0].st.Sub(r.st)})
	}

	// skip first range
	for i := 1; i < len(ranges); i++ {
		res = append(res, DateRange{st: ranges[i-1].End(), dur: ranges[i].st.Sub(ranges[i-1].End())})
	}

	// add the gap between the end of the last range and end of the period
	if !r.End().Equal(ranges[len(ranges)-1].End()) {
		res = append(res, DateRange{st: ranges[len(ranges)-1].End(), dur: r.End().Sub(ranges[len(ranges)-1].End())})
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
			res = append(res, DateRange{st: rangeStartTm, dur: boundary.tm.Sub(rangeStartTm)})
		}
	}

	// process the last boundary, it must be the end boundary anyway
	unfinishedBoundariesCnt--
	if unfinishedBoundariesCnt == 0 {
		res = append(res, DateRange{st: rangeStartTm, dur: boundaries[len(boundaries)-1].tm.Sub(rangeStartTm)})
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
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].st.Before(ranges[j].st) })
	return ranges
}

func rangesToBoundaries(ranges []DateRange) []*timeRangeBoundary {
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
