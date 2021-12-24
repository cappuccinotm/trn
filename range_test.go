package timerange

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type formattedRange struct {
	rng Range
	fmt string
}

func (r formattedRange) Format(layout string) string { return r.rng.Format(layout) }

func (r formattedRange) GoString() string { return r.Format(r.fmt) }

func (r formattedRange) String() string { return r.Format(r.fmt) }

func formattedRanges(rngs []Range, fmt string) []formattedRange {
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
		rng    Range
		bounds Range
		want   Range
	}{
		{
			name:   "doesn't intersect (earlier)",
			rng:    New(tm(13, 0), tm(14, 0)), // -XXX-----
			bounds: New(tm(15, 0), tm(16, 0)), // -----YYY-
			want:   Range{},
		},
		{
			name:   "doesn't intersect (later)",
			rng:    New(tm(15, 0), tm(16, 0)), // -----XXX-
			bounds: New(tm(13, 0), tm(14, 0)), // -YYY-----
			want:   Range{},
		},
		{
			name:   "overlaps the bounds",
			rng:    New(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			bounds: New(tm(14, 0), tm(15, 0)), // ---YYY---
			want:   New(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap",
			rng:    New(tm(14, 0), tm(15, 0)), // ---XXX---
			bounds: New(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:   New(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "intersect, bound end later",
			rng:    New(tm(13, 0), tm(15, 0)), // ---XXX---
			bounds: New(tm(14, 0), tm(16, 0)), // ----YYY--
			want:   New(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, starts are equal",
			rng:    New(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: New(tm(13, 0), tm(15, 0)), // --YYY----
			want:   New(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap, starts are equal",
			rng:    New(tm(13, 0), tm(15, 0)), // --XXX----
			bounds: New(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   New(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, ends are equal",
			rng:    New(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: New(tm(14, 0), tm(16, 0)), // ---YYY---
			want:   New(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "bounds overlap, ends are equal",
			rng:    New(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: New(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   New(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "intersect, bound end earlier",
			rng:    New(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: New(tm(13, 0), tm(15, 0)), // --YYY----
			want:   New(tm(14, 0), tm(15, 0)),
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
		rng   Range
		other Range
		want  bool
	}{
		{
			name:  "doesn't intersect (earlier)",
			rng:   New(tm(13, 0), tm(14, 0)), // -XXX-----
			other: New(tm(15, 0), tm(16, 0)), // -----YYY-
			want:  false,
		},
		{
			name:  "doesn't intersect (later)",
			rng:   New(tm(15, 0), tm(16, 0)), // -----XXX-
			other: New(tm(13, 0), tm(14, 0)), // -YYY-----
			want:  false,
		},
		{
			name:  "overlaps the bounds",
			rng:   New(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			other: New(tm(14, 0), tm(15, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap",
			rng:   New(tm(14, 0), tm(15, 0)), // ---XXX---
			other: New(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:  false,
		},
		{
			name:  "intersect, bound end later",
			rng:   New(tm(13, 0), tm(15, 0)), // ---XXX---
			other: New(tm(14, 0), tm(16, 0)), // ----YYY--
			want:  false,
		},
		{
			name:  "overlaps, starts are equal",
			rng:   New(tm(13, 0), tm(16, 0)), // --XXXX---
			other: New(tm(13, 0), tm(15, 0)), // --YYY----
			want:  true,
		},
		{
			name:  "bounds overlap, starts are equal",
			rng:   New(tm(13, 0), tm(15, 0)), // --XXX----
			other: New(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "overlaps, ends are equal",
			rng:   New(tm(13, 0), tm(16, 0)), // --XXXX---
			other: New(tm(14, 0), tm(16, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap, ends are equal",
			rng:   New(tm(14, 0), tm(16, 0)), // ---XXX---
			other: New(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "intersect, bound end earlier",
			rng:   New(tm(14, 0), tm(16, 0)), // ---XXX---
			other: New(tm(13, 0), tm(15, 0)), // --YYY----
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
		rng  Range
		args args
		want []Range
	}{
		{
			name: "space left at end",
			rng:  New(tm(1, 40), tm(2, 44)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				New(tm(1, 40), tm(2, 10)),
				New(tm(1, 45), tm(2, 15)),
				New(tm(1, 50), tm(2, 20)),
				New(tm(1, 55), tm(2, 25)),
				New(tm(2, 00), tm(2, 30)),
				New(tm(2, 05), tm(2, 35)),
				New(tm(2, 10), tm(2, 40)),
			},
		},
		{
			name: "without space left at end",
			rng:  New(tm(1, 40), tm(2, 40)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				New(tm(1, 40), tm(2, 10)),
				New(tm(1, 45), tm(2, 15)),
				New(tm(1, 50), tm(2, 20)),
				New(tm(1, 55), tm(2, 25)),
				New(tm(2, 00), tm(2, 30)),
				New(tm(2, 05), tm(2, 35)),
				New(tm(2, 10), tm(2, 40)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rng.Stratify(tt.args.duration, tt.args.interval)
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
		rng  Range
		args args
		want []Range
	}{
		{
			name: "space left at end",
			rng:  New(tm(1, 40), tm(3, 0)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				New(tm(1, 40), tm(2, 10)),
				New(tm(2, 15), tm(2, 45)),
			},
		},
		{
			name: "without space left at end",
			rng:  New(tm(1, 40), tm(3, 20)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				New(tm(1, 40), tm(2, 10)),
				New(tm(2, 15), tm(2, 45)),
				New(tm(2, 50), tm(3, 20)),
			},
		},
		{
			name: "zero interval",
			rng:  New(tm(1, 30), tm(2, 0)),
			args: args{duration: 5 * time.Minute},
			want: []Range{
				New(tm(1, 30), tm(1, 35)),
				New(tm(1, 35), tm(1, 40)),
				New(tm(1, 40), tm(1, 45)),
				New(tm(1, 45), tm(1, 50)),
				New(tm(1, 50), tm(1, 55)),
				New(tm(1, 55), tm(2, 0)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rng.Split(tt.args.duration, tt.args.interval)
			assert.Equal(t, formattedRanges(tt.want, "15:04"), formattedRanges(got, "15:04"))
		})
	}
}

func TestFlipDateRanges(t *testing.T) {
	type args struct {
		period Range
		ranges []Range
	}
	tests := []struct {
		name string
		fmt  string
		args args
		want []Range
	}{
		{
			name: "flip within a day", fmt: "15:04",
			args: args{
				period: New(tm(0, 0), tm(23, 59)),
				ranges: []Range{
					New(tm(13, 0), tm(14, 0)),
					New(tm(14, 1), tm(15, 0)),
					New(tm(16, 0), tm(20, 0)),
				},
			},
			want: []Range{
				New(tm(0, 0), tm(13, 0)),
				New(tm(14, 0), tm(14, 1)),
				New(tm(15, 0), tm(16, 0)),
				New(tm(20, 0), tm(23, 59)),
			},
		},
		{
			name: "no gap between the period and first, last range boundaries", fmt: "15:04",
			args: args{
				period: New(tm(0, 0), tm(23, 59)),
				ranges: []Range{
					New(tm(0, 0), tm(14, 0)),
					New(tm(14, 1), tm(15, 0)),
					New(tm(16, 0), tm(20, 0)),
					New(tm(20, 1), tm(23, 59)),
				},
			},
			want: []Range{
				New(tm(14, 0), tm(14, 1)),
				New(tm(15, 0), tm(16, 0)),
				New(tm(20, 0), tm(20, 1)),
			},
		},
		{
			name: "flip within several days", fmt: "02 15:04",
			args: args{
				period: New(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []Range{
					New(dhm(12, 13, 0), dhm(12, 14, 0)),
					New(dhm(12, 14, 1), dhm(12, 15, 0)),
					New(dhm(12, 16, 0), dhm(12, 20, 0)),
					New(dhm(12, 23, 0), dhm(13, 6, 59)),
					New(dhm(13, 8, 0), dhm(13, 23, 0)),
					New(dhm(14, 1, 59), dhm(14, 14, 59)),
				},
			},
			want: []Range{
				New(dhm(12, 14, 0), dhm(12, 14, 1)),
				New(dhm(12, 15, 0), dhm(12, 16, 0)),
				New(dhm(12, 20, 0), dhm(12, 23, 0)),
				New(dhm(13, 6, 59), dhm(13, 8, 0)),
				New(dhm(13, 23, 0), dhm(14, 1, 59)),
				New(dhm(14, 14, 59), dhm(14, 16, 59)),
			},
		},
		{name: "empty range list", fmt: "02 15:04",
			args: args{
				period: New(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []Range{},
			},
			want: []Range{New(dhm(12, 13, 0), dhm(14, 16, 59))},
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

func TestDateRange_Format(t *testing.T) {
	assert.Equal(t,
		"[2021-06-12T00:00:00, 2021-06-12T03:05:00]",
		Range{st: dt, dur: 3*time.Hour + 5*time.Minute}.Format("2006-01-02T15:04:05"),
	)
}

func TestDateRange_UTC(t *testing.T) {
	// won't have effect on machine in UTC ¯\_(ツ)_/¯
	assert.Equal(t, Range{st: dt.In(time.Local), dur: 0}, New(dt, dt).In(time.Local))
}

func TestDateRange_String(t *testing.T) {
	assert.Equal(t,
		"[2021-06-12 00:00:00 +0000 UTC, 2021-06-12 03:05:00 +0000 UTC]",
		Range{st: dt, dur: 3*time.Hour + 5*time.Minute}.String(),
	)
}

func TestDateRange_Duration(t *testing.T) {
	dur := 3*time.Hour + 5*time.Minute
	assert.Equal(t, dur, Range{dur: dur}.Duration())
}

func TestDateRange_Start(t *testing.T) {
	assert.Equal(t, dt, Range{st: dt}.Start())
}

func TestDateRange_Empty(t *testing.T) {
	tests := []struct {
		name string
		arg  Range
		want bool
	}{
		{name: "duration and start ts empty", arg: Range{}, want: true},
		{name: "duration empty", arg: Range{st: dt}, want: false},
		{name: "ts empty", arg: Range{dur: 3 * time.Hour}, want: false},
		{name: "duration and start ts filled", arg: Range{st: dt, dur: 3 * time.Hour}, want: false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.arg.Empty())
	}
}

func TestRange(t *testing.T) {
	t.Run("start after end", func(t *testing.T) {
		assert.Panics(t, func() { New(dt.Add(3*time.Hour), dt) })
	})

	t.Run("without options", func(t *testing.T) {
		assert.Equal(t, Range{st: dt, dur: 3 * time.Hour}, New(dt, dt.Add(3*time.Hour)))
	})

	t.Run("with location", func(t *testing.T) {
		dr := New(dt.In(time.Local), dt.Add(3*time.Hour).In(time.Local), In(time.UTC))

		// won't have effect on machine in UTC ¯\_(ツ)_/¯
		assert.Equal(t, Range{st: dt, dur: 3 * time.Hour}, dr)
	})
}
