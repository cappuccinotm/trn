package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClock_Sub(t *testing.T) {
	dur := NewClock(13, 12, 11, 10, time.UTC).Sub(NewClock(9, 8, 7, 6, time.UTC))
	res := 4*time.Nanosecond + 4*time.Second + 4*time.Minute + 4*time.Hour
	assert.Equal(t, res, dur)
}
