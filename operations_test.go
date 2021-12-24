package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeOverlappingRanges(t *testing.T) {
	tests := []struct {
		name string
		args []DateRange
		want []DateRange
	}{
		{
			name: "ranges don't overlap",
			args: []DateRange{
				Range(tm(13, 0), tm(14, 0)),
				Range(tm(15, 0), tm(16, 0)),
			},
			want: []DateRange{
				Range(tm(13, 0), tm(14, 0)),
				Range(tm(15, 0), tm(16, 0)),
			},
		},
		{
			name: "ranges intersect",
			args: []DateRange{
				Range(tm(13, 0), tm(14, 0)),
				Range(tm(13, 30), tm(15, 0)),
			},
			want: []DateRange{
				Range(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "one range eternally overlaps the other",
			args: []DateRange{
				Range(tm(13, 0), tm(15, 0)),
				Range(tm(13, 30), tm(14, 30)),
			},
			want: []DateRange{
				Range(tm(13, 0), tm(15, 0)),
			},
		},
		{
			name: "boundaries of two ranges are equal",
			args: []DateRange{
				Range(tm(13, 0), tm(13, 15)),
				Range(tm(13, 15), tm(13, 30)),
			},
			want: []DateRange{
				Range(tm(13, 0), tm(13, 30)),
			},
		},
		{
			name: "complex test",
			args: []DateRange{
				// next three ranges must be merged (last two are within the first one)
				Range(tm(19, 0), tm(19, 30)),
				Range(tm(19, 1), tm(19, 15)),
				Range(tm(19, 8), tm(19, 17)),
				// next second range must be removed (end of first = end of second)
				Range(tm(15, 0), tm(15, 30)),
				Range(tm(15, 16), tm(15, 30)),
				// next two ranges must NOT be merged
				Range(tm(12, 0), tm(12, 15)),
				Range(tm(12, 30), tm(12, 45)),
				// next two ranges must be merged (end of the first = start of the second)
				Range(tm(13, 0), tm(13, 15)),
				Range(tm(13, 15), tm(13, 30)),
				// next two ranges must be merged
				Range(tm(14, 0), tm(14, 16)),
				Range(tm(14, 15), tm(14, 30)),
				// next second range must be removed (start of first = start of second)
				Range(tm(16, 0), tm(16, 30)),
				Range(tm(16, 0), tm(16, 16)),
				// next second range must be removed (ranges are equal)
				Range(tm(17, 0), tm(17, 30)),
				Range(tm(17, 0), tm(17, 30)),
				// next second range must be removed
				Range(tm(18, 0), tm(18, 30)),
				Range(tm(18, 1), tm(18, 15)),
			},
			want: []DateRange{
				Range(tm(12, 0), tm(12, 15)),
				Range(tm(12, 30), tm(12, 45)),
				Range(tm(13, 0), tm(13, 30)),
				Range(tm(14, 0), tm(14, 30)),
				Range(tm(15, 0), tm(15, 30)),
				Range(tm(16, 0), tm(16, 30)),
				Range(tm(17, 0), tm(17, 30)),
				Range(tm(18, 0), tm(18, 30)),
				Range(tm(19, 0), tm(19, 30)),
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
		args []DateRange
		want DateRange
	}{
		{name: "empty list", args: []DateRange{}, want: DateRange{}},
		{name: "nil list", args: nil, want: DateRange{}},
		{
			name: "intersection",
			args: []DateRange{
				Range(tm(13, 0), tm(19, 0)),
				Range(tm(15, 0), tm(17, 0)),
				Range(tm(16, 0), tm(21, 0)),
			},
			want: Range(tm(16, 0), tm(17, 0)),
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
