package helpers

import (
	"time"
)

func CloneDate(t time.Time) *time.Time {
	var cloned = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return &cloned
}
func CloneDateTrimmed(t time.Time) *time.Time {
	var cloned = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.UTC().Location())
	return &cloned
}
