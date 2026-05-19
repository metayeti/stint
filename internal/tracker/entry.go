package tracker

import "time"

type Entry struct {
	Date        string
	TimeSpan    string
	Description string
	Duration    time.Duration
}