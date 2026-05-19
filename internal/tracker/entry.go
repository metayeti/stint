package tracker

import "time"

type Entry struct {
	Date        string
	TimeSpan    string
	Description string
	Duration    time.Duration
}

// max width for description column before wrapping
const DescriptionWidth = 35