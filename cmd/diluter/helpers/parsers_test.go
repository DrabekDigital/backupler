package helpers

import (
	"math"
	"testing"
)

func TestParseDaysValid(t *testing.T) {
	days, err := ParseDays("0 day", false)
	if days != 0 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 0, days, err)
	}

	days, err = ParseDays("0 days", false)
	if days != 0 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 0, days, err)
	}

	days, err = ParseDays("1 day", false)
	if days != 1 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 1, days, err)
	}

	days, err = ParseDays("1 days", false)
	if days != 1 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 1, days, err)
	}

	days, err = ParseDays("2 day", false)
	if days != 2 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 2, days, err)
	}

	days, err = ParseDays("2 days", false)
	if days != 2 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 2, days, err)
	}

	days, err = ParseDays("1250 days", false)
	if days != 1250 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, 1250, days, err)
	}

	days, err = ParseDays("infinity", true)
	if days != math.MaxInt64 || err != nil {
		t.Fatalf(`expected %d to be parsed: %d with error %v`, math.MaxInt64, days, err)
	}
}

func TestParseDaysInvalid(t *testing.T) {

	days, err := ParseDays(" 0 day", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("0 day ", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("1 dayz", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("-1 day", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("-1 days", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("-2 days", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("+1 day", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("+2 days", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}

	days, err = ParseDays("infinity", false)
	if err == nil {
		t.Fatalf(`expected to not parsed, instead %d returned`, days)
	}
}

func TestValidParseBackupDate(t *testing.T) {
	var compareFormat = "2006-01-02 15:04:05"
	time, err := ParseDate("2022-02-28_12-33-40", "yyyy-MM-dd_HH-mm-ss")
	formatted := time.UTC().Format(compareFormat)
	if formatted != "2022-02-28 12:33:40" || err != nil {
		t.Fatalf(`date expected to be obtained from backup, instead value %v and error %v returned`, formatted, err)
	}

	time, err = ParseDate("20220228123340", "yyyyMMddHHmmss")
	formatted = time.UTC().Format(compareFormat)
	if formatted != "2022-02-28 12:33:40" || err != nil {
		t.Fatalf(`date expected to be obtained from backup, instead value %v and error %v returned`, formatted, err)
	}

	time, err = ParseDate("12.33.40-28.02.2022", "HH.mm.ss-dd.MM.yyyy")
	formatted = time.UTC().Format(compareFormat)
	if formatted != "2022-02-28 12:33:40" || err != nil {
		t.Fatalf(`date expected to be obtained from backup, instead value %v and error %v returned`, formatted, err)
	}
}

func TestInvalidParseBackupDate(t *testing.T) {
	time, _ := ParseDate("2022-02-28_12-33-40_copy", "yyyy-MM-dd_HH-mm-ss")

	if time != nil {
		t.Fatalf(`date expected to be not obtained from backup`)
	}

	time, _ = ParseDate("20220228 123340", "yyyyMMddHHmmss")

	if time != nil {
		t.Fatalf(`date expected to be not obtained from backup`)
	}
}
