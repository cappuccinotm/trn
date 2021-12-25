package timerange

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
				Between(tm(13, 0), tm(14, 0)),
				Between(tm(15, 0), tm(16, 0)),
			},
			want: []Range{
				Between(tm(13, 0), tm(14, 0)),
				Between(tm(15, 0), tm(16, 0)),
			},
		},
		{
			name: "ranges intersect",
			args: []Range{
				Between(tm(13, 0), tm(14, 0)),
				Between(tm(13, 30), tm(15, 0)),
			},
			want: []Range{
				Between(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "one range eternally overlaps the other",
			args: []Range{
				Between(tm(13, 0), tm(15, 0)),
				Between(tm(13, 30), tm(14, 30)),
			},
			want: []Range{
				Between(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "boundaries of two ranges are equal",
			args: []Range{
				Between(tm(13, 0), tm(13, 15)),
				Between(tm(13, 15), tm(13, 30)),
			},
			want: []Range{
				Between(tm(13, 0), tm(13, 30)),
			},
		},
		{
			name: "complex test",
			args: []Range{
				// next three ranges must be merged (last two are within the first one)
				Between(tm(19, 0), tm(19, 30)),
				Between(tm(19, 1), tm(19, 15)),
				Between(tm(19, 8), tm(19, 17)),
				// next second range must be removed (end of first = end of second)
				Between(tm(15, 0), tm(15, 30)),
				Between(tm(15, 16), tm(15, 30)),
				// next two ranges must NOT be merged
				Between(tm(12, 0), tm(12, 15)),
				Between(tm(12, 30), tm(12, 45)),
				// next two ranges must be merged (end of the first = start of the second)
				Between(tm(13, 0), tm(13, 15)),
				Between(tm(13, 15), tm(13, 30)),
				// next two ranges must be merged
				Between(tm(14, 0), tm(14, 16)),
				Between(tm(14, 15), tm(14, 30)),
				// next second range must be removed (start of first = start of second)
				Between(tm(16, 0), tm(16, 30)),
				Between(tm(16, 0), tm(16, 16)),
				// next second range must be removed (ranges are equal)
				Between(tm(17, 0), tm(17, 30)),
				Between(tm(17, 0), tm(17, 30)),
				// next second range must be removed
				Between(tm(18, 0), tm(18, 30)),
				Between(tm(18, 1), tm(18, 15)),
			},
			want: []Range{
				Between(tm(12, 0), tm(12, 15)),
				Between(tm(12, 30), tm(12, 45)),
				Between(tm(13, 0), tm(13, 30)),
				Between(tm(14, 0), tm(14, 30)),
				Between(tm(15, 0), tm(15, 30)),
				Between(tm(16, 0), tm(16, 30)),
				Between(tm(17, 0), tm(17, 30)),
				Between(tm(18, 0), tm(18, 30)),
				Between(tm(19, 0), tm(19, 30)),
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
				Between(tm(13, 0), tm(19, 0)),
				Between(tm(15, 0), tm(17, 0)),
				Between(tm(16, 0), tm(21, 0)),
			},
			want: Between(tm(16, 0), tm(17, 0)),
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
