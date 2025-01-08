package examples

import "time"

// parse is a helper function to parse a time string in RFC3339 format. It will panic if the string is not a valid time.
func parse(s string) time.Time {
	if tm, err := time.Parse(time.RFC3339, s); err == nil {
		return tm
	} else {
		panic(err)
	}
}
