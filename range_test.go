package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

var dt = time.Date(2021, 6, 12, 0, 0, 0, 0, time.UTC)

func tm(h, m int) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), h, m, 0, 0, time.UTC)
}

func dhm(d, h, m int) time.Time {
	return time.Date(dt.Year(), dt.Month(), d, h, m, 0, 0, time.UTC)
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
			rng:    Range(tm(13, 0), tm(14, 0)), // -XXX-----
			bounds: Range(tm(15, 0), tm(16, 0)), // -----YYY-
			want:   DateRange{},
		},
		{
			name:   "doesn't intersect (later)",
			rng:    Range(tm(15, 0), tm(16, 0)), // -----XXX-
			bounds: Range(tm(13, 0), tm(14, 0)), // -YYY-----
			want:   DateRange{},
		},
		{
			name:   "overlaps the bounds",
			rng:    Range(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			bounds: Range(tm(14, 0), tm(15, 0)), // ---YYY---
			want:   Range(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap",
			rng:    Range(tm(14, 0), tm(15, 0)), // ---XXX---
			bounds: Range(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:   Range(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "intersect, bound end later",
			rng:    Range(tm(13, 0), tm(15, 0)), // ---XXX---
			bounds: Range(tm(14, 0), tm(16, 0)), // ----YYY--
			want:   Range(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, starts are equal",
			rng:    Range(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: Range(tm(13, 0), tm(15, 0)), // --YYY----
			want:   Range(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap, starts are equal",
			rng:    Range(tm(13, 0), tm(15, 0)), // --XXX----
			bounds: Range(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   Range(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, ends are equal",
			rng:    Range(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: Range(tm(14, 0), tm(16, 0)), // ---YYY---
			want:   Range(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "bounds overlap, ends are equal",
			rng:    Range(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: Range(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   Range(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "intersect, bound end earlier",
			rng:    Range(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: Range(tm(13, 0), tm(15, 0)), // --YYY----
			want:   Range(tm(14, 0), tm(15, 0)),
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
			rng:   Range(tm(13, 0), tm(14, 0)), // -XXX-----
			other: Range(tm(15, 0), tm(16, 0)), // -----YYY-
			want:  false,
		},
		{
			name:  "doesn't intersect (later)",
			rng:   Range(tm(15, 0), tm(16, 0)), // -----XXX-
			other: Range(tm(13, 0), tm(14, 0)), // -YYY-----
			want:  false,
		},
		{
			name:  "overlaps the bounds",
			rng:   Range(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			other: Range(tm(14, 0), tm(15, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap",
			rng:   Range(tm(14, 0), tm(15, 0)), // ---XXX---
			other: Range(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:  false,
		},
		{
			name:  "intersect, bound end later",
			rng:   Range(tm(13, 0), tm(15, 0)), // ---XXX---
			other: Range(tm(14, 0), tm(16, 0)), // ----YYY--
			want:  false,
		},
		{
			name:  "overlaps, starts are equal",
			rng:   Range(tm(13, 0), tm(16, 0)), // --XXXX---
			other: Range(tm(13, 0), tm(15, 0)), // --YYY----
			want:  true,
		},
		{
			name:  "bounds overlap, starts are equal",
			rng:   Range(tm(13, 0), tm(15, 0)), // --XXX----
			other: Range(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "overlaps, ends are equal",
			rng:   Range(tm(13, 0), tm(16, 0)), // --XXXX---
			other: Range(tm(14, 0), tm(16, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap, ends are equal",
			rng:   Range(tm(14, 0), tm(16, 0)), // ---XXX---
			other: Range(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "intersect, bound end earlier",
			rng:   Range(tm(14, 0), tm(16, 0)), // ---XXX---
			other: Range(tm(13, 0), tm(15, 0)), // --YYY----
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
			rng:  Range(tm(1, 34), tm(2, 44)),
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				Range(tm(1, 40), tm(2, 10)),
				Range(tm(1, 45), tm(2, 15)),
				Range(tm(1, 50), tm(2, 20)),
				Range(tm(1, 55), tm(2, 25)),
				Range(tm(2, 00), tm(2, 30)),
				Range(tm(2, 05), tm(2, 35)),
				Range(tm(2, 10), tm(2, 40)),
			},
		},
		{
			name: "without space left at end",
			rng:  Range(tm(1, 34), tm(2, 40)),
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				Range(tm(1, 40), tm(2, 10)),
				Range(tm(1, 45), tm(2, 15)),
				Range(tm(1, 50), tm(2, 20)),
				Range(tm(1, 55), tm(2, 25)),
				Range(tm(2, 00), tm(2, 30)),
				Range(tm(2, 05), tm(2, 35)),
				Range(tm(2, 10), tm(2, 40)),
			},
		},
		{
			name: "zero offset",
			rng:  Range(tm(1, 30), tm(2, 0)),
			args: args{duration: 10 * time.Minute, interval: 5 * time.Minute},
			want: []DateRange{
				Range(tm(1, 30), tm(1, 40)),
				Range(tm(1, 35), tm(1, 45)),
				Range(tm(1, 40), tm(1, 50)),
				Range(tm(1, 45), tm(1, 55)),
				Range(tm(1, 50), tm(2, 0)),
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
			rng:  Range(tm(1, 34), tm(3, 0)),
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				Range(tm(1, 40), tm(2, 10)),
				Range(tm(2, 15), tm(2, 45)),
			},
		},
		{
			name: "without space left at end",
			rng:  Range(tm(1, 34), tm(3, 20)),
			args: args{
				offset:   6 * time.Minute,
				duration: 30 * time.Minute,
				interval: 5 * time.Minute,
			},
			want: []DateRange{
				Range(tm(1, 40), tm(2, 10)),
				Range(tm(2, 15), tm(2, 45)),
				Range(tm(2, 50), tm(3, 20)),
			},
		},
		{
			name: "zero offset and interval",
			rng:  Range(tm(1, 30), tm(2, 0)),
			args: args{duration: 5 * time.Minute},
			want: []DateRange{
				Range(tm(1, 30), tm(1, 35)),
				Range(tm(1, 35), tm(1, 40)),
				Range(tm(1, 40), tm(1, 45)),
				Range(tm(1, 45), tm(1, 50)),
				Range(tm(1, 50), tm(1, 55)),
				Range(tm(1, 55), tm(2, 0)),
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
				period: Range(tm(0, 0), tm(23, 59)),
				ranges: []DateRange{
					Range(tm(13, 0), tm(14, 0)),
					Range(tm(14, 1), tm(15, 0)),
					Range(tm(16, 0), tm(20, 0)),
				},
			},
			want: []DateRange{
				Range(tm(0, 0), tm(13, 0)),
				Range(tm(14, 0), tm(14, 1)),
				Range(tm(15, 0), tm(16, 0)),
				Range(tm(20, 0), tm(23, 59)),
			},
		},
		{
			name: "no gap between the period and first, last range boundaries", fmt: "15:04",
			args: args{
				period: Range(tm(0, 0), tm(23, 59)),
				ranges: []DateRange{
					Range(tm(0, 0), tm(14, 0)),
					Range(tm(14, 1), tm(15, 0)),
					Range(tm(16, 0), tm(20, 0)),
					Range(tm(20, 1), tm(23, 59)),
				},
			},
			want: []DateRange{
				Range(tm(14, 0), tm(14, 1)),
				Range(tm(15, 0), tm(16, 0)),
				Range(tm(20, 0), tm(20, 1)),
			},
		},
		{
			name: "flip within several days", fmt: "02 15:04",
			args: args{
				period: Range(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []DateRange{
					Range(dhm(12, 13, 0), dhm(12, 14, 0)),
					Range(dhm(12, 14, 1), dhm(12, 15, 0)),
					Range(dhm(12, 16, 0), dhm(12, 20, 0)),
					Range(dhm(12, 23, 0), dhm(13, 6, 59)),
					Range(dhm(13, 8, 0), dhm(13, 23, 0)),
					Range(dhm(14, 1, 59), dhm(14, 14, 59)),
				},
			},
			want: []DateRange{
				Range(dhm(12, 14, 0), dhm(12, 14, 1)),
				Range(dhm(12, 15, 0), dhm(12, 16, 0)),
				Range(dhm(12, 20, 0), dhm(12, 23, 0)),
				Range(dhm(13, 6, 59), dhm(13, 8, 0)),
				Range(dhm(13, 23, 0), dhm(14, 1, 59)),
				Range(dhm(14, 14, 59), dhm(14, 16, 59)),
			},
		},
		{name: "empty range list", fmt: "02 15:04",
			args: args{
				period: Range(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []DateRange{},
			},
			want: []DateRange{Range(dhm(12, 13, 0), dhm(14, 16, 59))},
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
