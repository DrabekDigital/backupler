package helpers

import (
	"testing"
	"time"
)

func TestCloneDate(t *testing.T) {
	var timestamp = time.Now()
	var cloned = CloneDate(timestamp)

	if timestamp == *cloned {
		t.Fatalf(`cloned and original datetimes are the same instance`)
	}
	if timestamp.Year() != (*cloned).Year() {
		t.Fatalf(`years do not match`)
	}
	if timestamp.Month() != (*cloned).Month() {
		t.Fatalf(`months do not match`)
	}
	if timestamp.Day() != (*cloned).Day() {
		t.Fatalf(`days do not match`)
	}
	if timestamp.Hour() != (*cloned).Hour() {
		t.Fatalf(`hours do not match`)
	}
	if timestamp.Minute() != (*cloned).Minute() {
		t.Fatalf(`minutes do not match`)
	}
	if timestamp.Second() != (*cloned).Second() {
		t.Fatalf(`seconds do not match`)
	}
	if timestamp.Nanosecond() != (*cloned).Nanosecond() {
		t.Fatalf(`nanoseconds do not match`)
	}
	if timestamp.Location() != (*cloned).Location() {
		t.Fatalf(`locations do not match`)
	}
}
