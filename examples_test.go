package trn

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestMergeRanges(t *testing.T) {
	now := time.Now()

	rng := New(now, time.Hour+3*time.Minute)
	ranges := rng.Stratify(15*time.Minute, 5*time.Minute)
	assert.Equal(t, []Range{
		New(now, 15*time.Minute),
		New(now.Add(5*time.Minute), 15*time.Minute),
		New(now.Add(10*time.Minute), 15*time.Minute),
		New(now.Add(15*time.Minute), 15*time.Minute),
		New(now.Add(20*time.Minute), 15*time.Minute),
		New(now.Add(25*time.Minute), 15*time.Minute),
		New(now.Add(30*time.Minute), 15*time.Minute),
		New(now.Add(35*time.Minute), 15*time.Minute),
		New(now.Add(40*time.Minute), 15*time.Minute),
		New(now.Add(45*time.Minute), 15*time.Minute),
	}, ranges)
}

func TestTruncate(t *testing.T) {
	now := time.Now()

	rng := New(now, time.Hour)
	truncateRange := New(now.Add(15*time.Minute), time.Hour)
	assert.Equal(t,
		New(now.Add(15*time.Minute), 45*time.Minute),
		rng.Truncate(truncateRange),
	)
}

func TestFlip(t *testing.T) {
	now := time.Now()

	flipRange := New(now, time.Hour)

	rngs := []Range{
		New(now.Add(25*time.Minute), 15*time.Minute),
		New(now.Add(50*time.Minute), 5*time.Minute),
	}

	assert.Equal(t,
		[]Range{
			New(now, 25*time.Minute),
			New(now.Add(40*time.Minute), 10*time.Minute),
			New(now.Add(55*time.Minute), 5*time.Minute),
		},
		flipRange.Flip(rngs),
	)

}
