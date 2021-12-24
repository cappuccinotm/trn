package store

import (
	"fmt"
	"time"
)

// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock struct{ time.Time }

// NewClock returns the Clock in the given location with given hours, minutes and secs
func NewClock(h, m, s, ns int, loc *time.Location) Clock {
	return Clock{Time: time.Date(0, time.January, 1, h, m, s, ns, loc)}
}

// ClockFromTime returns the clock extracted from the given time.Time.
func ClockFromTime(t time.Time) Clock {
	return Clock{t}
}

// Sub returns the duration between the clock at the date of the other time and current clock
func (c Clock) Sub(other Clock) time.Duration {
	return c.Time.Sub(other.Time)
}

// String implements fmt.Stringer to print and log Clock properly
func (c Clock) String() string {
	return fmt.Sprintf("%02d:%02d:%02d %s", c.Hour(), c.Minute(), c.Second(), c.Location())
}

// GoString implements fmt.GoStringer to use Clock in %#v formats
func (c Clock) GoString() string {
	return fmt.Sprintf("NewClock(%d, %d, %d, %s)", c.Hour(), c.Minute(), c.Second(), c.Location())
}
