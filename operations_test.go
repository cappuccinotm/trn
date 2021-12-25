package trn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeOverlappingRanges(t *testing.T) {
	tests := []struct {
		name string
		args []Range
		want []Range
	}{
		{
			name: "ranges don't overlap",
			args: []Range{
				MustBetween(tm(13, 0), tm(14, 0)),
				MustBetween(tm(15, 0), tm(16, 0)),
			},
			want: []Range{
				MustBetween(tm(13, 0), tm(14, 0)),
				MustBetween(tm(15, 0), tm(16, 0)),
			},
		},
		{
			name: "ranges intersect",
			args: []Range{
				MustBetween(tm(13, 0), tm(14, 0)),
				MustBetween(tm(13, 30), tm(15, 0)),
			},
			want: []Range{
				MustBetween(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "one range eternally overlaps the other",
			args: []Range{
				MustBetween(tm(13, 0), tm(15, 0)),
				MustBetween(tm(13, 30), tm(14, 30)),
			},
			want: []Range{
				MustBetween(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "boundaries of two ranges are equal",
			args: []Range{
				MustBetween(tm(13, 0), tm(13, 15)),
				MustBetween(tm(13, 15), tm(13, 30)),
			},
			want: []Range{
				MustBetween(tm(13, 0), tm(13, 30)),
			},
		},
		{
			name: "complex test",
			args: []Range{
				// next three ranges must be merged (last two are within the first one)
				MustBetween(tm(19, 0), tm(19, 30)),
				MustBetween(tm(19, 1), tm(19, 15)),
				MustBetween(tm(19, 8), tm(19, 17)),
				// next second range must be removed (end of first = end of second)
				MustBetween(tm(15, 0), tm(15, 30)),
				MustBetween(tm(15, 16), tm(15, 30)),
				// next two ranges must NOT be merged
				MustBetween(tm(12, 0), tm(12, 15)),
				MustBetween(tm(12, 30), tm(12, 45)),
				// next two ranges must be merged (end of the first = start of the second)
				MustBetween(tm(13, 0), tm(13, 15)),
				MustBetween(tm(13, 15), tm(13, 30)),
				// next two ranges must be merged
				MustBetween(tm(14, 0), tm(14, 16)),
				MustBetween(tm(14, 15), tm(14, 30)),
				// next second range must be removed (start of first = start of second)
				MustBetween(tm(16, 0), tm(16, 30)),
				MustBetween(tm(16, 0), tm(16, 16)),
				// next second range must be removed (ranges are equal)
				MustBetween(tm(17, 0), tm(17, 30)),
				MustBetween(tm(17, 0), tm(17, 30)),
				// next second range must be removed
				MustBetween(tm(18, 0), tm(18, 30)),
				MustBetween(tm(18, 1), tm(18, 15)),
			},
			want: []Range{
				MustBetween(tm(12, 0), tm(12, 15)),
				MustBetween(tm(12, 30), tm(12, 45)),
				MustBetween(tm(13, 0), tm(13, 30)),
				MustBetween(tm(14, 0), tm(14, 30)),
				MustBetween(tm(15, 0), tm(15, 30)),
				MustBetween(tm(16, 0), tm(16, 30)),
				MustBetween(tm(17, 0), tm(17, 30)),
				MustBetween(tm(18, 0), tm(18, 30)),
				MustBetween(tm(19, 0), tm(19, 30)),
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

func TestIntersection(t *testing.T) {
	tests := []struct {
		name string
		args []Range
		want Range
	}{
		{name: "empty list", args: []Range{}, want: Range{}},
		{name: "nil list", args: nil, want: Range{}},
		{
			name: "intersection",
			args: []Range{
				MustBetween(tm(13, 0), tm(19, 0)),
				MustBetween(tm(15, 0), tm(17, 0)),
				MustBetween(tm(16, 0), tm(21, 0)),
			},
			want: MustBetween(tm(16, 0), tm(17, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intersection := Intersection(tt.args)
			assert.Equal(t,
				formattedRange{rng: tt.want, fmt: "15:04"},
				formattedRange{rng: intersection, fmt: "15:04"},
			)
		})
	}
}
