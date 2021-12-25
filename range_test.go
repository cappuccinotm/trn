package trn

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

func TestRange_Truncate(t *testing.T) {
	tests := []struct {
		name   string
		rng    Range
		bounds Range
		want   Range
	}{
		{
			name:   "doesn't intersect (earlier)",
			rng:    MustBetween(tm(13, 0), tm(14, 0)), // -XXX-----
			bounds: MustBetween(tm(15, 0), tm(16, 0)), // -----YYY-
			want:   Range{},
		},
		{
			name:   "doesn't intersect (later)",
			rng:    MustBetween(tm(15, 0), tm(16, 0)), // -----XXX-
			bounds: MustBetween(tm(13, 0), tm(14, 0)), // -YYY-----
			want:   Range{},
		},
		{
			name:   "overlaps the bounds",
			rng:    MustBetween(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			bounds: MustBetween(tm(14, 0), tm(15, 0)), // ---YYY---
			want:   MustBetween(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap",
			rng:    MustBetween(tm(14, 0), tm(15, 0)), // ---XXX---
			bounds: MustBetween(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:   MustBetween(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "intersect, bound end later",
			rng:    MustBetween(tm(13, 0), tm(15, 0)), // ---XXX---
			bounds: MustBetween(tm(14, 0), tm(16, 0)), // ----YYY--
			want:   MustBetween(tm(14, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, starts are equal",
			rng:    MustBetween(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: MustBetween(tm(13, 0), tm(15, 0)), // --YYY----
			want:   MustBetween(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "bounds overlap, starts are equal",
			rng:    MustBetween(tm(13, 0), tm(15, 0)), // --XXX----
			bounds: MustBetween(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   MustBetween(tm(13, 0), tm(15, 0)),
		},
		{
			name:   "overlaps, ends are equal",
			rng:    MustBetween(tm(13, 0), tm(16, 0)), // --XXXX---
			bounds: MustBetween(tm(14, 0), tm(16, 0)), // ---YYY---
			want:   MustBetween(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "bounds overlap, ends are equal",
			rng:    MustBetween(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: MustBetween(tm(13, 0), tm(16, 0)), // --YYYY---
			want:   MustBetween(tm(14, 0), tm(16, 0)),
		},
		{
			name:   "intersect, bound end earlier",
			rng:    MustBetween(tm(14, 0), tm(16, 0)), // ---XXX---
			bounds: MustBetween(tm(13, 0), tm(15, 0)), // --YYY----
			want:   MustBetween(tm(14, 0), tm(15, 0)),
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

func TestRange_Contains(t *testing.T) {
	tests := []struct {
		name  string
		rng   Range
		other Range
		want  bool
	}{
		{
			name:  "doesn't intersect (earlier)",
			rng:   MustBetween(tm(13, 0), tm(14, 0)), // -XXX-----
			other: MustBetween(tm(15, 0), tm(16, 0)), // -----YYY-
			want:  false,
		},
		{
			name:  "doesn't intersect (later)",
			rng:   MustBetween(tm(15, 0), tm(16, 0)), // -----XXX-
			other: MustBetween(tm(13, 0), tm(14, 0)), // -YYY-----
			want:  false,
		},
		{
			name:  "overlaps the bounds",
			rng:   MustBetween(tm(13, 0), tm(16, 0)), // -XXXXXXX-
			other: MustBetween(tm(14, 0), tm(15, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap",
			rng:   MustBetween(tm(14, 0), tm(15, 0)), // ---XXX---
			other: MustBetween(tm(13, 0), tm(16, 0)), // -YYYYYYY-
			want:  false,
		},
		{
			name:  "intersect, bound end later",
			rng:   MustBetween(tm(13, 0), tm(15, 0)), // ---XXX---
			other: MustBetween(tm(14, 0), tm(16, 0)), // ----YYY--
			want:  false,
		},
		{
			name:  "overlaps, starts are equal",
			rng:   MustBetween(tm(13, 0), tm(16, 0)), // --XXXX---
			other: MustBetween(tm(13, 0), tm(15, 0)), // --YYY----
			want:  true,
		},
		{
			name:  "bounds overlap, starts are equal",
			rng:   MustBetween(tm(13, 0), tm(15, 0)), // --XXX----
			other: MustBetween(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "overlaps, ends are equal",
			rng:   MustBetween(tm(13, 0), tm(16, 0)), // --XXXX---
			other: MustBetween(tm(14, 0), tm(16, 0)), // ---YYY---
			want:  true,
		},
		{
			name:  "bounds overlap, ends are equal",
			rng:   MustBetween(tm(14, 0), tm(16, 0)), // ---XXX---
			other: MustBetween(tm(13, 0), tm(16, 0)), // --YYYY---
			want:  false,
		},
		{
			name:  "intersect, bound end earlier",
			rng:   MustBetween(tm(14, 0), tm(16, 0)), // ---XXX---
			other: MustBetween(tm(13, 0), tm(15, 0)), // --YYY----
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rng.Contains(tt.other))
		})
	}
}

func TestRange_Stratify(t *testing.T) {
	type args struct {
		duration time.Duration
		interval time.Duration
	}
	tests := []struct {
		name    string
		rng     Range
		args    args
		want    []Range
		wantErr error
	}{
		{
			name: "space left at end",
			rng:  MustBetween(tm(1, 40), tm(2, 44)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				MustBetween(tm(1, 40), tm(2, 10)),
				MustBetween(tm(1, 45), tm(2, 15)),
				MustBetween(tm(1, 50), tm(2, 20)),
				MustBetween(tm(1, 55), tm(2, 25)),
				MustBetween(tm(2, 00), tm(2, 30)),
				MustBetween(tm(2, 05), tm(2, 35)),
				MustBetween(tm(2, 10), tm(2, 40)),
			},
		},
		{
			name: "without space left at end",
			rng:  MustBetween(tm(1, 40), tm(2, 40)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				MustBetween(tm(1, 40), tm(2, 10)),
				MustBetween(tm(1, 45), tm(2, 15)),
				MustBetween(tm(1, 50), tm(2, 20)),
				MustBetween(tm(1, 55), tm(2, 25)),
				MustBetween(tm(2, 00), tm(2, 30)),
				MustBetween(tm(2, 05), tm(2, 35)),
				MustBetween(tm(2, 10), tm(2, 40)),
			},
		},
		{name: "zero interval", args: args{interval: 0}, wantErr: ErrZeroDurationInterval},
		{name: "zero duration", args: args{duration: 0}, wantErr: ErrZeroDurationInterval},
		{name: "negative duration", args: args{duration: -10}, wantErr: ErrZeroDurationInterval},
		{name: "negative interval", args: args{interval: -10}, wantErr: ErrZeroDurationInterval},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rng.Stratify(tt.args.duration, tt.args.interval)
			assert.Equal(t, formattedRanges(tt.want, "15:04"), formattedRanges(got, "15:04"))
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRange_Split(t *testing.T) {
	type args struct {
		duration time.Duration
		interval time.Duration
	}
	tests := []struct {
		name    string
		rng     Range
		args    args
		want    []Range
		wantErr error
	}{
		{
			name: "space left at end",
			rng:  MustBetween(tm(1, 40), tm(3, 0)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				MustBetween(tm(1, 40), tm(2, 10)),
				MustBetween(tm(2, 15), tm(2, 45)),
			},
		},
		{
			name: "without space left at end",
			rng:  MustBetween(tm(1, 40), tm(3, 20)),
			args: args{duration: 30 * time.Minute, interval: 5 * time.Minute},
			want: []Range{
				MustBetween(tm(1, 40), tm(2, 10)),
				MustBetween(tm(2, 15), tm(2, 45)),
				MustBetween(tm(2, 50), tm(3, 20)),
			},
		},
		{
			name: "zero interval",
			rng:  MustBetween(tm(1, 30), tm(2, 0)),
			args: args{duration: 5 * time.Minute},
			want: []Range{
				MustBetween(tm(1, 30), tm(1, 35)),
				MustBetween(tm(1, 35), tm(1, 40)),
				MustBetween(tm(1, 40), tm(1, 45)),
				MustBetween(tm(1, 45), tm(1, 50)),
				MustBetween(tm(1, 50), tm(1, 55)),
				MustBetween(tm(1, 55), tm(2, 0)),
			},
		},
		{
			name:    "zero duration",
			rng:     MustBetween(tm(1, 30), tm(2, 0)),
			args:    args{duration: 0, interval: 5 * time.Minute},
			wantErr: ErrZeroDurationInterval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rng.Split(tt.args.duration, tt.args.interval)
			assert.Equal(t, formattedRanges(tt.want, "15:04"), formattedRanges(got, "15:04"))
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestFlipRanges(t *testing.T) {
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
				period: MustBetween(tm(0, 0), tm(23, 59)),
				ranges: []Range{
					MustBetween(tm(13, 0), tm(14, 0)),
					MustBetween(tm(14, 1), tm(15, 0)),
					MustBetween(tm(16, 0), tm(20, 0)),
				},
			},
			want: []Range{
				MustBetween(tm(0, 0), tm(13, 0)),
				MustBetween(tm(14, 0), tm(14, 1)),
				MustBetween(tm(15, 0), tm(16, 0)),
				MustBetween(tm(20, 0), tm(23, 59)),
			},
		},
		{
			name: "no gap between the period and first, last range boundaries", fmt: "15:04",
			args: args{
				period: MustBetween(tm(0, 0), tm(23, 59)),
				ranges: []Range{
					MustBetween(tm(0, 0), tm(14, 0)),
					MustBetween(tm(14, 1), tm(15, 0)),
					MustBetween(tm(16, 0), tm(20, 0)),
					MustBetween(tm(20, 1), tm(23, 59)),
				},
			},
			want: []Range{
				MustBetween(tm(14, 0), tm(14, 1)),
				MustBetween(tm(15, 0), tm(16, 0)),
				MustBetween(tm(20, 0), tm(20, 1)),
			},
		},
		{
			name: "flip within several days", fmt: "02 15:04",
			args: args{
				period: MustBetween(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []Range{
					MustBetween(dhm(12, 13, 0), dhm(12, 14, 0)),
					MustBetween(dhm(12, 14, 1), dhm(12, 15, 0)),
					MustBetween(dhm(12, 16, 0), dhm(12, 20, 0)),
					MustBetween(dhm(12, 23, 0), dhm(13, 6, 59)),
					MustBetween(dhm(13, 8, 0), dhm(13, 23, 0)),
					MustBetween(dhm(14, 1, 59), dhm(14, 14, 59)),
				},
			},
			want: []Range{
				MustBetween(dhm(12, 14, 0), dhm(12, 14, 1)),
				MustBetween(dhm(12, 15, 0), dhm(12, 16, 0)),
				MustBetween(dhm(12, 20, 0), dhm(12, 23, 0)),
				MustBetween(dhm(13, 6, 59), dhm(13, 8, 0)),
				MustBetween(dhm(13, 23, 0), dhm(14, 1, 59)),
				MustBetween(dhm(14, 14, 59), dhm(14, 16, 59)),
			},
		},
		{name: "empty range list", fmt: "02 15:04",
			args: args{
				period: MustBetween(dhm(12, 13, 0), dhm(14, 16, 59)),
				ranges: []Range{},
			},
			want: []Range{MustBetween(dhm(12, 13, 0), dhm(14, 16, 59))},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ranges := tt.args.period.Flip(tt.args.ranges)
			assert.Equal(t,
				formattedRanges(tt.want, tt.fmt),
				formattedRanges(ranges, tt.fmt),
			)
		})
	}
}

func TestRange_Format(t *testing.T) {
	assert.Equal(t,
		"[2021-06-12T00:00:00, 2021-06-12T03:05:00]",
		Range{st: dt, dur: 3*time.Hour + 5*time.Minute}.Format("2006-01-02T15:04:05"),
	)
}

func TestRange_UTC(t *testing.T) {
	// won't have effect on machine in UTC ¯\_(ツ)_/¯
	assert.Equal(t, Range{st: dt.In(time.Local), dur: 0}, MustBetween(dt, dt).In(time.Local))
}

func TestRange_String(t *testing.T) {
	assert.Equal(t,
		"[2021-06-12 00:00:00 +0000 UTC, 2021-06-12 03:05:00 +0000 UTC]",
		Range{st: dt, dur: 3*time.Hour + 5*time.Minute}.String(),
	)
}

func TestRange_Duration(t *testing.T) {
	dur := 3*time.Hour + 5*time.Minute
	assert.Equal(t, dur, Range{dur: dur}.Duration())
}

func TestRange_Start(t *testing.T) {
	assert.Equal(t, dt, Range{st: dt}.Start())
}

func TestRange_Empty(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.Empty())
		})
	}
}

func TestBetween(t *testing.T) {
	type args struct {
		start, end time.Time
		opts       []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Range
		wantErr error
	}{
		{
			name:    "start after end",
			args:    args{start: dt.Add(3 * time.Hour), end: dt},
			wantErr: ErrStartAfterEnd,
		},
		{
			name: "without options",
			args: args{start: dt, end: dt.Add(3 * time.Hour)},
			want: Range{st: dt, dur: 3 * time.Hour},
		},
		{
			name: "with location",
			args: args{start: dt, end: dt.Add(3 * time.Hour), opts: []Option{In(time.Local)}},
			want: Range{st: dt.In(time.Local), dur: 3 * time.Hour},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng, err := Between(tt.args.start, tt.args.end, tt.args.opts...)
			assert.Equal(t, tt.want, rng)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestMustSplit(t *testing.T) {
	assert.Panics(t, func() {
		New(tm(1, 30), 30*time.Minute).MustSplit(0, 15*time.Minute)
	})

	assert.NotPanics(t, func() {
		rngs := New(tm(1, 40), 1*time.Hour+20*time.Minute).
			MustSplit(30*time.Minute, 5*time.Minute)
		assert.Equal(t, []Range{
			MustBetween(tm(1, 40), tm(2, 10)),
			MustBetween(tm(2, 15), tm(2, 45)),
		}, rngs)
	})
}

func TestMustStratify(t *testing.T) {
	assert.Panics(t, func() {
		New(tm(1, 30), 30*time.Minute).MustStratify(0, 15*time.Minute)
	})

	assert.NotPanics(t, func() {
		rngs := New(tm(1, 40), 20*time.Minute).
			MustStratify(10*time.Minute, 5*time.Minute)
		assert.Equal(t, []Range{
			New(tm(1, 40), 10*time.Minute),
			New(tm(1, 45), 10*time.Minute),
			New(tm(1, 50), 10*time.Minute),
		}, rngs)
	})
}

func TestNew(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		assert.Equal(t, Range{st: dt, dur: 3 * time.Hour}, New(dt, 3*time.Hour))
	})

	t.Run("with location", func(t *testing.T) {
		dr := New(dt.In(time.Local), 3*time.Hour, In(time.UTC))

		// won't have effect on machine in UTC ¯\_(ツ)_/¯
		assert.Equal(t, Range{st: dt, dur: 3 * time.Hour}, dr)
	})
}

func TestRange_GoString(t *testing.T) {
	assert.Equal(t,
		"trn.New(time.Date(2021, time.December, 25, 18, 34, 30, 0, time.UTC), 900000000000)",
		Range{
			st:  time.Date(2021, time.December, 25, 18, 34, 30, 0, time.UTC),
			dur: 900000000000, // 15 * time.Minute
		}.GoString(),
	)
}

func TestError_Error(t *testing.T) {
	assert.Equal(t, "blah", Error("blah").Error())
}

func TestMustBetween(t *testing.T) {
	assert.Panics(t, func() { MustBetween(dt.Add(3*time.Hour), dt) })
	assert.NotPanics(t, func() {
		assert.Equal(t,
			Range{st: dt, dur: 3 * time.Hour,},
			MustBetween(dt, dt.Add(3*time.Hour)),
		)
	})
}
