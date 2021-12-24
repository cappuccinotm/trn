package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dt = Date{Year: 2021, Month: 6, Day: 12}

type formattedRange struct {
	rng DateRange
	fmt string
}

func (r formattedRange) Format(layout string) string { return r.rng.Format(layout) }

func (r formattedRange) GoString() string { return r.Format(r.fmt) }

func (r formattedRange) String() string { return r.Format(r.fmt) }

func formattedRanges(rngs []DateRange, fmt string) []formattedRange {
	res := make([]formattedRange, len(rngs))
	for i := range rngs {
		res[i] = formattedRange{
			rng: rngs[i],
			fmt: fmt,
		}
	}
	return res
}

func formattedRangeMap(m map[Date][]DateRange, fmt string) map[Date][]formattedRange {
	res := map[Date][]formattedRange{}
	for k, v := range m {
		res[k] = formattedRanges(v, fmt)
	}
	return res
}

func tm(h, m int) time.Time {
	return tmd(dt.Day, h, m)
}

func tmd(d, h, m int) time.Time {
	return tmns(d, h, m, 0, 0)
}

func tmns(d, h, m, s, ns int) time.Time {
	return time.Date(dt.Year, dt.Month, d, h, m, s, ns, time.UTC)
}

func TestDateRange_Truncate(t *testing.T) {
	tests := []struct {
		name   string
		rng    DateRange
		bounds DateRange
		want   DateRange
	}{
		{
			name:   "doesn't intersect (earlier)",
			rng:    DateRange{Start: tm(13, 0), End: tm(14, 0)}, // -XXX-----
			bounds: DateRange{Start: tm(15, 0), End: tm(16, 0)}, // -----YYY-
			want:   DateRange{},
		},
		{
			name:   "doesn't intersect (later)",
			rng:    DateRange{Start: tm(15, 0), End: tm(16, 0)}, // -----XXX-
			bounds: DateRange{Start: tm(13, 0), End: tm(14, 0)}, // -YYY-----
			want:   DateRange{},
		},
		{
			name:   "overlaps the bounds",
			rng:    DateRange{Start: tm(13, 0), End: tm(16, 0)}, // -XXXXXXX-
			bounds: DateRange{Start: tm(14, 0), End: tm(15, 0)}, // ---YYY---
			want:   DateRange{Start: tm(14, 0), End: tm(15, 0)},
		},
		{
			name:   "bounds overlap",
			rng:    DateRange{Start: tm(14, 0), End: tm(15, 0)}, // ---XXX---
			bounds: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // -YYYYYYY-
			want:   DateRange{Start: tm(14, 0), End: tm(15, 0)},
		},
		{
			name:   "intersect, bound end later",
			rng:    DateRange{Start: tm(13, 0), End: tm(15, 0)}, // ---XXX---
			bounds: DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ----YYY--
			want:   DateRange{Start: tm(14, 0), End: tm(15, 0)},
		},
		{
			name:   "overlaps, starts are equal",
			rng:    DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --XXXX---
			bounds: DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --YYY----
			want:   DateRange{Start: tm(13, 0), End: tm(15, 0)},
		},
		{
			name:   "bounds overlap, starts are equal",
			rng:    DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --XXX----
			bounds: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --YYYY---
			want:   DateRange{Start: tm(13, 0), End: tm(15, 0)},
		},
		{
			name:   "overlaps, ends are equal",
			rng:    DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --XXXX---
			bounds: DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---YYY---
			want:   DateRange{Start: tm(14, 0), End: tm(16, 0)},
		},
		{
			name:   "bounds overlap, ends are equal",
			rng:    DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---XXX---
			bounds: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --YYYY---
			want:   DateRange{Start: tm(14, 0), End: tm(16, 0)},
		},
		{
			name:   "intersect, bound end earlier",
			rng:    DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---XXX---
			bounds: DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --YYY----
			want:   DateRange{Start: tm(14, 0), End: tm(15, 0)},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rng := tt.rng.Truncate(tt.bounds)
			assert.Equal(t,
				formattedRange{rng: tt.want, fmt: "15:04"},
				formattedRange{rng: rng, fmt: "15:04"},
			)
		})
	}
}

func TestDateRange_Contains(t *testing.T) {
	tests := []struct {
		name  string
		rng   DateRange
		other DateRange
		want  bool
	}{
		{
			name:  "doesn't intersect (earlier)",
			rng:   DateRange{Start: tm(13, 0), End: tm(14, 0)}, // -XXX-----
			other: DateRange{Start: tm(15, 0), End: tm(16, 0)}, // -----YYY-
			want:  false,
		},
		{
			name:  "doesn't intersect (later)",
			rng:   DateRange{Start: tm(15, 0), End: tm(16, 0)}, // -----XXX-
			other: DateRange{Start: tm(13, 0), End: tm(14, 0)}, // -YYY-----
			want:  false,
		},
		{
			name:  "overlaps the bounds",
			rng:   DateRange{Start: tm(13, 0), End: tm(16, 0)}, // -XXXXXXX-
			other: DateRange{Start: tm(14, 0), End: tm(15, 0)}, // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap",
			rng:   DateRange{Start: tm(14, 0), End: tm(15, 0)}, // ---XXX---
			other: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // -YYYYYYY-
			want:  false,
		},
		{
			name:  "intersect, bound end later",
			rng:   DateRange{Start: tm(13, 0), End: tm(15, 0)}, // ---XXX---
			other: DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ----YYY--
			want:  false,
		},
		{
			name:  "overlaps, starts are equal",
			rng:   DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --XXXX---
			other: DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --YYY----
			want:  true,
		},
		{
			name:  "bounds overlap, starts are equal",
			rng:   DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --XXX----
			other: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --YYYY---
			want:  false,
		},
		{
			name:  "overlaps, ends are equal",
			rng:   DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --XXXX---
			other: DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap, ends are equal",
			rng:   DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---XXX---
			other: DateRange{Start: tm(13, 0), End: tm(16, 0)}, // --YYYY---
			want:  false,
		},
		{
			name:  "intersect, bound end earlier",
			rng:   DateRange{Start: tm(14, 0), End: tm(16, 0)}, // ---XXX---
			other: DateRange{Start: tm(13, 0), End: tm(15, 0)}, // --YYY----
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rng.Contains(tt.other))
		})
	}
}

func TestDateRange_Stratify(t *testing.T) {
	type args struct {
		offset   time.Duration
		duration time.Duration
		interval time.Duration
	}
	tests := []struct {
		name string
		rng  DateRange
		args args
		want []DateRange
	}{
		{
			name: "space left at end",
			rng:  DateRange{Start: tm(1, 34), End: tm(2, 44)},
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				{Start: tm(1, 40), End: tm(2, 10)},
				{Start: tm(1, 45), End: tm(2, 15)},
				{Start: tm(1, 50), End: tm(2, 20)},
				{Start: tm(1, 55), End: tm(2, 25)},
				{Start: tm(2, 00), End: tm(2, 30)},
				{Start: tm(2, 05), End: tm(2, 35)},
				{Start: tm(2, 10), End: tm(2, 40)},
			},
		},
		{
			name: "without space left at end",
			rng:  DateRange{Start: tm(1, 34), End: tm(2, 40)},
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				{Start: tm(1, 40), End: tm(2, 10)},
				{Start: tm(1, 45), End: tm(2, 15)},
				{Start: tm(1, 50), End: tm(2, 20)},
				{Start: tm(1, 55), End: tm(2, 25)},
				{Start: tm(2, 00), End: tm(2, 30)},
				{Start: tm(2, 05), End: tm(2, 35)},
				{Start: tm(2, 10), End: tm(2, 40)},
			},
		},
		{
			name: "zero offset",
			rng:  DateRange{Start: tm(1, 30), End: tm(2, 0)},
			args: args{duration: 10 * time.Minute, interval: 5 * time.Minute},
			want: []DateRange{
				{Start: tm(1, 30), End: tm(1, 40)},
				{Start: tm(1, 35), End: tm(1, 45)},
				{Start: tm(1, 40), End: tm(1, 50)},
				{Start: tm(1, 45), End: tm(1, 55)},
				{Start: tm(1, 50), End: tm(2, 0)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rng.Stratify(tt.args.offset, tt.args.duration, tt.args.interval)
			assert.Equal(t, formattedRanges(tt.want, "15:04"), formattedRanges(got, "15:04"))
		})
	}
}

func TestDateRange_Split(t *testing.T) {
	type args struct {
		offset   time.Duration
		duration time.Duration
		interval time.Duration
	}
	tests := []struct {
		name string
		rng  DateRange
		args args
		want []DateRange
	}{
		{
			name: "space left at end",
			rng:  DateRange{Start: tm(1, 34), End: tm(3, 0)},
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				{Start: tm(1, 40), End: tm(2, 10)},
				{Start: tm(2, 15), End: tm(2, 45)},
			},
		},
		{
			name: "without space left at end",
			rng:  DateRange{Start: tm(1, 34), End: tm(3, 20)},
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				{Start: tm(1, 40), End: tm(2, 10)},
				{Start: tm(2, 15), End: tm(2, 45)},
				{Start: tm(2, 50), End: tm(3, 20)},
			},
		},
		{
			name: "zero offset and interval",
			rng:  DateRange{Start: tm(1, 30), End: tm(2, 0)},
			args: args{duration: 5 * time.Minute},
			want: []DateRange{
				{Start: tm(1, 30), End: tm(1, 35)},
				{Start: tm(1, 35), End: tm(1, 40)},
				{Start: tm(1, 40), End: tm(1, 45)},
				{Start: tm(1, 45), End: tm(1, 50)},
				{Start: tm(1, 50), End: tm(1, 55)},
				{Start: tm(1, 55), End: tm(2, 0)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rng.Split(tt.args.offset, tt.args.duration, tt.args.interval)
			assert.Equal(t, formattedRanges(tt.want, "15:04"), formattedRanges(got, "15:04"))
		})
	}
}

func TestMergeOverlappingRanges(t *testing.T) {
	tests := []struct {
		name string
		args []DateRange
		want []DateRange
	}{
		{
			name: "ranges don't overlap",
			args: []DateRange{
				{Start: tm(13, 0), End: tm(14, 0)},
				{Start: tm(15, 0), End: tm(16, 0)},
			},
			want: []DateRange{
				{Start: tm(13, 0), End: tm(14, 0)},
				{Start: tm(15, 0), End: tm(16, 0)},
			},
		},
		{
			name: "ranges intersect",
			args: []DateRange{
				{Start: tm(13, 0), End: tm(14, 0)},
				{Start: tm(13, 30), End: tm(15, 0)},
			},
			want: []DateRange{
				{Start: tm(13, 0), End: tm(15, 0)},
			},
		},
		{
			name: "one range eternally overlaps the other",
			args: []DateRange{
				{Start: tm(13, 0), End: tm(15, 0)},
				{Start: tm(13, 30), End: tm(14, 30)},
			},
			want: []DateRange{
				{Start: tm(13, 0), End: tm(15, 0)},
			},
		},
		{
			name: "boundaries of two ranges are equal",
			args: []DateRange{
				{Start: tm(13, 0), End: tm(13, 15)},
				{Start: tm(13, 15), End: tm(13, 30)},
			},
			want: []DateRange{
				{Start: tm(13, 0), End: tm(13, 30)},
			},
		},
		{
			name: "complex test",
			args: []DateRange{
				// next three ranges must be merged (last two are within the first one)
				{Start: tm(19, 0), End: tm(19, 30)},
				{Start: tm(19, 1), End: tm(19, 15)},
				{Start: tm(19, 8), End: tm(19, 17)},
				// next second range must be removed (end of first = end of second)
				{Start: tm(15, 0), End: tm(15, 30)},
				{Start: tm(15, 16), End: tm(15, 30)},
				// next two ranges must NOT be merged
				{Start: tm(12, 0), End: tm(12, 15)},
				{Start: tm(12, 30), End: tm(12, 45)},
				// next two ranges must be merged (end of the first = start of the second)
				{Start: tm(13, 0), End: tm(13, 15)},
				{Start: tm(13, 15), End: tm(13, 30)},
				// next two ranges must be merged
				{Start: tm(14, 0), End: tm(14, 16)},
				{Start: tm(14, 15), End: tm(14, 30)},
				// next second range must be removed (start of first = start of second)
				{Start: tm(16, 0), End: tm(16, 30)},
				{Start: tm(16, 0), End: tm(16, 16)},
				// next second range must be removed (ranges are equal)
				{Start: tm(17, 0), End: tm(17, 30)},
				{Start: tm(17, 0), End: tm(17, 30)},
				// next second range must be removed
				{Start: tm(18, 0), End: tm(18, 30)},
				{Start: tm(18, 1), End: tm(18, 15)},
			},
			want: []DateRange{
				{Start: tm(12, 0), End: tm(12, 15)},
				{Start: tm(12, 30), End: tm(12, 45)},
				{Start: tm(13, 0), End: tm(13, 30)},
				{Start: tm(14, 0), End: tm(14, 30)},
				{Start: tm(15, 0), End: tm(15, 30)},
				{Start: tm(16, 0), End: tm(16, 30)},
				{Start: tm(17, 0), End: tm(17, 30)},
				{Start: tm(18, 0), End: tm(18, 30)},
				{Start: tm(19, 0), End: tm(19, 30)},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ranges := MergeOverlappingRanges(tt.args)
			assert.Equal(t,
				formattedRanges(tt.want, "15:04"),
				formattedRanges(ranges, "15:04"),
			)
		})
	}
}

func TestFlipDateRanges(t *testing.T) {
	type args struct {
		period DateRange
		ranges []DateRange
	}
	tests := []struct {
		name string
		fmt  string
		args args
		want []DateRange
	}{
		{
			name: "flip within a day", fmt: "15:04",
			args: args{
				period: DateRange{Start: tm(0, 0), End: tm(23, 59)},
				ranges: []DateRange{
					{Start: tm(13, 0), End: tm(14, 0)},
					{Start: tm(14, 1), End: tm(15, 0)},
					{Start: tm(16, 0), End: tm(20, 0)},
				},
			},
			want: []DateRange{
				{Start: tm(0, 0), End: tm(13, 0)},
				{Start: tm(14, 0), End: tm(14, 1)},
				{Start: tm(15, 0), End: tm(16, 0)},
				{Start: tm(20, 0), End: tm(23, 59)},
			},
		},
		{
			name: "no gap between the period and first, last range boundaries", fmt: "15:04",
			args: args{
				period: DateRange{Start: tm(0, 0), End: tm(23, 59)},
				ranges: []DateRange{
					{Start: tm(0, 0), End: tm(14, 0)},
					{Start: tm(14, 1), End: tm(15, 0)},
					{Start: tm(16, 0), End: tm(20, 0)},
					{Start: tm(20, 1), End: tm(23, 59)},
				},
			},
			want: []DateRange{
				{Start: tm(14, 0), End: tm(14, 1)},
				{Start: tm(15, 0), End: tm(16, 0)},
				{Start: tm(20, 0), End: tm(20, 1)},
			},
		},
		{
			name: "flip within several days", fmt: "02 15:04",
			args: args{
				period: DateRange{Start: tmd(12, 13, 0), End: tmd(14, 16, 59)},
				ranges: []DateRange{
					{Start: tmd(12, 13, 0), End: tmd(12, 14, 0)},
					{Start: tmd(12, 14, 1), End: tmd(12, 15, 0)},
					{Start: tmd(12, 16, 0), End: tmd(12, 20, 0)},
					{Start: tmd(12, 23, 0), End: tmd(13, 6, 59)},
					{Start: tmd(13, 8, 0), End: tmd(13, 23, 0)},
					{Start: tmd(14, 1, 59), End: tmd(14, 14, 59)},
				},
			},
			want: []DateRange{
				{Start: tmd(12, 14, 0), End: tmd(12, 14, 1)},
				{Start: tmd(12, 15, 0), End: tmd(12, 16, 0)},
				{Start: tmd(12, 20, 0), End: tmd(12, 23, 0)},
				{Start: tmd(13, 6, 59), End: tmd(13, 8, 0)},
				{Start: tmd(13, 23, 0), End: tmd(14, 1, 59)},
				{Start: tmd(14, 14, 59), End: tmd(14, 16, 59)},
			},
		},
		{name: "empty range list", fmt: "02 15:04",
			args: args{
				period: DateRange{Start: tmd(12, 13, 0), End: tmd(14, 16, 59)},
				ranges: []DateRange{},
			},
			want: []DateRange{{Start: tmd(12, 13, 0), End: tmd(14, 16, 59)}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ranges := tt.args.period.FlipDateRanges(tt.args.ranges)
			assert.Equal(t,
				formattedRanges(tt.want, tt.fmt),
				formattedRanges(ranges, tt.fmt),
			)
		})
	}
}

func TestSplitToRangesPerDay(t *testing.T) {
	const fmt = "2006-01-02 15:04:05.999999999"

	dt := func(d int) Date { return Date{Year: 2021, Month: 6, Day: d} }

	tests := []struct {
		name   string
		ranges []DateRange
		want   map[Date][]DateRange
	}{
		{
			name:   "two days, without boundaries of day",
			ranges: []DateRange{{Start: tmd(12, 13, 0), End: tmd(13, 14, 0)}},
			want: map[Date][]DateRange{
				dt(12): {{Start: tmd(12, 13, 0), End: tmd(13, 0, 0)}},
				dt(13): {{Start: tmd(13, 0, 0), End: tmd(13, 14, 0)}},
			},
		},
		{
			name:   "range with boundary on 00:00",
			ranges: []DateRange{{Start: tmd(1, 10, 0), End: tmd(3, 0, 0)}},
			want: map[Date][]DateRange{
				dt(1): {{Start: tmd(1, 10, 0), End: tmd(1, 24, 0)}},
				dt(2): {{Start: tmd(2, 0, 0), End: tmd(2, 24, 0)}},
			},
		},
		{
			name: "several ranges per day",
			ranges: []DateRange{
				{Start: tmd(1, 10, 0), End: tmd(1, 13, 0)},
				{Start: tmd(1, 14, 0), End: tmd(2, 8, 0)},
				{Start: tmd(2, 8, 0), End: tmd(3, 0, 0)},
			},
			want: map[Date][]DateRange{
				dt(1): {
					{Start: tmd(1, 10, 0), End: tmd(1, 13, 0)},
					{Start: tmd(1, 14, 0), End: tmd(2, 0, 0)},
				},
				dt(2): {
					{Start: tmd(2, 0, 0), End: tmd(2, 8, 0)},
					{Start: tmd(2, 8, 0), End: tmd(3, 0, 0)},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ranges := SplitToRangesPerDay(tt.ranges)
			assert.Equal(t,
				formattedRangeMap(tt.want, fmt),
				formattedRangeMap(ranges, fmt),
			)
		})
	}
}
