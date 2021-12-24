// Package timerange provides structures and methods to perform complex
// operations over time ranges.
package timerange

import (
	"fmt"
	"time"
)

const defaultRangeFmt = "2006-01-02 15:04:05.999999999 -0700 MST"

// Option is an adapter over date ranges.
type Option func(r *Range)

// In sets the time range in the given location.
func In(loc *time.Location) Option {
	return func(r *Range) { r.st = r.st.In(loc) }
}

// New returns the new Range in the given time bounds.
// Range will get the location of the start timestamp.
func New(start, end time.Time, opts ...Option) Range {
	if start.After(end) {
		panic("start is after the end")
	}

	res := Range{st: start, dur: end.Sub(start)}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

// Range represents time slot with its own start and end time boundaries
type Range struct {
	st  time.Time
	dur time.Duration
}

// String implements fmt.Stringer to print and log Range properly
func (r Range) String() string { return r.UTC().Format(defaultRangeFmt) }

// UTC returns the date range with boundaries in UTC.
func (r Range) UTC() Range { return r.In(time.UTC) }

// Duration returns the duration of the date range.
func (r Range) Duration() time.Duration { return r.dur }

// Start returns the start time of the date range.
func (r Range) Start() time.Time { return r.st }

// End returns the end time of the date range.
func (r Range) End() time.Time { return r.st.Add(r.dur) }

// In returns the date range with boundaries in the provided location's time zone.
func (r Range) In(loc *time.Location) Range { return Range{st: r.st.In(loc), dur: r.dur} }

// Empty returns true if the date range is empty.
func (r Range) Empty() bool { return r.st.IsZero() && r.dur == 0 }

// Format returns the string representation of the time range with the given format.
func (r Range) Format(layout string) string {
	return fmt.Sprintf("[%s, %s]", r.st.Format(layout), r.End().Format(layout))
}

// Split the date range into smaller ranges, with fixed duration and with the
// given interval between the *end* of the one range and *start* of next range.
func (r Range) Split(duration time.Duration, interval time.Duration) []Range {
	if duration == 0 {
		panic("cannot split with zero duration")
	}
	return r.Stratify(duration, duration+interval)
}

// Stratify the date range into smaller ranges, with fixed duration and with the
// given interval between the *starts* of the resulting ranges.
func (r Range) Stratify(duration time.Duration, interval time.Duration) []Range {
	if interval == 0 || duration == 0 {
		panic("cannot stratify with zero duration or zero interval")
	}

	var res []Range
	rangeEnd := r.End()
	rangeStart := r.st

	for rangeEnd.Sub(rangeStart.Add(duration)) >= 0 {
		res = append(res, Range{st: rangeStart, dur: duration})
		rangeStart = rangeStart.Add(interval)
	}

	return res
}

// Contains returns true if the other date range is within this date range.
func (r Range) Contains(other Range) bool {
	if (r.st.Before(other.st) || r.st.Equal(other.st)) &&
		(r.End().After(other.End()) || r.End().Equal(other.End())) {
		return true
	}
	return false
}

// Truncate returns the date range bounded to the *bounds*, i.e. it cuts
// the start and the end of *r* to fit into the *bounds*.
func (r Range) Truncate(bounds Range) Range {
	switch {
	case r.st.Before(bounds.st) && r.End().Before(bounds.st):
		// -XXX-----
		// -----YYY-
		return Range{}
	case r.st.After(bounds.End()) && r.End().After(bounds.End()):
		// -----XXX-
		// -YYY-----
		return Range{}
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
		return Range{st: bounds.st, dur: r.End().Sub(bounds.st)}
	case r.st.After(bounds.st) && r.End().After(bounds.End()):
		// ---XXX---
		// --YYY----
		return Range{st: r.st, dur: bounds.End().Sub(r.st)}
	default:
		panic("should never happen")
	}
}

// FlipDateRanges within the given period.
//
// The boundaries of the given ranges are considered to be inclusive, means
// that the flipped ranges will start or end at the exact nanosecond where
// the boundary from the input starts or ends.
func (r Range) FlipDateRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{r}
	}

	// to exclude the case of distinct ranges, ranges not within the period
	// and unsorted list of ranges
	rngs := MergeOverlappingRanges(ranges)

	return r.flipValidRanges(rngs)
}

func (r Range) flipValidRanges(ranges []Range) []Range {
	var res []Range

	// add the gap between the start of the period and start of the first range
	if !r.st.Equal(ranges[0].st) {
		res = append(res, Range{st: r.st, dur: ranges[0].st.Sub(r.st)})
	}

	// skip first range
	for i := 1; i < len(ranges); i++ {
		res = append(res, Range{st: ranges[i-1].End(), dur: ranges[i].st.Sub(ranges[i-1].End())})
	}

	// add the gap between the end of the last range and end of the period
	if !r.End().Equal(ranges[len(ranges)-1].End()) {
		res = append(res, Range{st: ranges[len(ranges)-1].End(), dur: r.End().Sub(ranges[len(ranges)-1].End())})
	}

	return res
}
