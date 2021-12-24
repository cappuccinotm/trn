package store

import "time"

// Date represents a single date without any information about the Clock.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// Time returns the time.Time that represents this date, with Clock information provided.
func (dt Date) Time(c Clock) time.Time {
	return time.Date(
		dt.Year,
		dt.Month,
		dt.Day,
		c.Hour(),
		c.Minute(),
		c.Second(),
		c.Nanosecond(),
		c.Location(),
	)
}

// After checks that the current date is after the other date.
func (dt Date) After(other Date) bool {
	if dt.Year == other.Year {
		if dt.Month == other.Month {
			return dt.Day > other.Day
		}
		return dt.Month > other.Month
	}
	return dt.Year > other.Year
}

// Before checks that the current date is before the given date.
func (dt Date) Before(other Date) bool {
	if dt.Year == other.Year {
		if dt.Month == other.Month {
			return dt.Day < other.Day
		}
		return dt.Month < other.Month
	}
	return dt.Year < other.Year
}

// BeforeOrEqual checks that the current date is before or equal the other date.
func (dt Date) BeforeOrEqual(other Date) bool {
	return dt.Before(other) || dt.Equal(other)
}

// AfterOrEqual checks that the current date is after or equal the other date.
func (dt Date) AfterOrEqual(other Date) bool {
	return dt.After(other) || dt.Equal(other)
}

// Equal returns true if the dates are the same.
func (dt Date) Equal(other Date) bool {
	return dt.Year == other.Year && dt.Month == other.Month && dt.Day == other.Day
}

// Add some time to the current date.
func (dt Date) Add(y int, m int, d int) Date {
	return DateFromTime(time.Date(
		dt.Year+y, dt.Month+time.Month(m), dt.Day+d,
		0, 0, 0, 0, time.UTC))
}

// DateFromTime returns the Date extracted from the given time.Time
func DateFromTime(t time.Time) Date {
	y, m, d := t.Date()
	return Date{Year: y, Month: m, Day: d}
}
