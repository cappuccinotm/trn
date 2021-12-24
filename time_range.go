package store

import "fmt"

// TimeRange represents the timeslot within a day.
type TimeRange struct {
	Start Clock `json:"start"`
	End   Clock `json:"end"`
}

// GoString implements fmt.GoStringer to use TimeRange in %#v formats
func (r TimeRange) GoString() string { return r.String() }

// String implements fmt.Stringer to print and log DateRange properly
func (r TimeRange) String() string {
	return fmt.Sprintf("[%s, %s]", r.Start.String(), r.End.String())
}
