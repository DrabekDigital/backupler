package helpers

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"drabek.cz/cli-utils/backupler/cmd/diluter/definitions"
)

func TestListbackups(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf(`current working directory could not been obtained due to %v`, err)
	}
	basePath := filepath.Join(pwd, "/../../../fixtures/tests/listing-backups")
	actualBackups, err2 := ListBackups(basePath, "yyyyMMdd_HHmmss")
	if err2 != nil {
		t.Fatalf(`list of backups in certain directory could not be obtained due to %v`, err2)
	}
	if len(actualBackups) != 5 {
		t.Fatalf(`5 backups expected to be listed, instead %v found`, len(actualBackups))
	}

	if actualBackups[0].Path != filepath.Join(basePath, "20180101_070103") {
		t.Fatalf(`Backup 20180101_070103 expected`)
	}
	if actualBackups[1].Path != filepath.Join(basePath, "20200513_125539") {
		t.Fatalf(`Backup 20200513_125539 expected`)
	}
	if actualBackups[2].Path != filepath.Join(basePath, "20210711_202101") {
		t.Fatalf(`Backup 20210711_202101 expected`)
	}
	if actualBackups[3].Path != filepath.Join(basePath, "20220822_045533") {
		t.Fatalf(`Backup 20220822_045533 expected`)
	}
	if actualBackups[4].Path != filepath.Join(basePath, "20220823_121429") {
		t.Fatalf(`Backup 20220823_121429 expected`)
	}

	var compareFormat = "2006-01-02 15:04:05"
	if actualBackups[0].Creation.UTC().Format(compareFormat) != "2018-01-01 07:01:03" {
		t.Fatalf(`Backup 20180101_070103 expected to match date 2018-01-01 07:01:03`)
	}
	if actualBackups[1].Creation.UTC().Format(compareFormat) != "2020-05-13 12:55:39" {
		t.Fatalf(`Backup 20200513_125539 expected to match date 2020-05-13 12:55:39`)
	}
	if actualBackups[2].Creation.UTC().Format(compareFormat) != "2021-07-11 20:21:01" {
		t.Fatalf(`Backup 20210711_202101 expected to match date 2021-07-11 20:21:01`)
	}
	if actualBackups[3].Creation.UTC().Format(compareFormat) != "2022-08-22 04:55:33" {
		t.Fatalf(`Backup 20220822_045533 expected to match date 2022-08-22 04:55:33`)
	}
	if actualBackups[4].Creation.UTC().Format(compareFormat) != "2022-08-23 12:14:29" {
		t.Fatalf(`Backup 20220823_121429 expected to match date 2022-08-23 12:14:29`)
	}

	if actualBackups[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20180101_070103 expected to be set as KEEP`)
	}
	if actualBackups[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20200513_125539 expected to be set as KEEP`)
	}
	if actualBackups[2].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20210711_202101 expected to be set as KEEP`)
	}
	if actualBackups[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20220822_045533 expected to be set as KEEP`)
	}
	if actualBackups[4].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20220823_121429 expected to be set as KEEP`)
	}
}

func TestApplyPoliciesDeletion(t *testing.T) {
	fixedPoint := time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name: "delete",
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20180101_070103",
			Creation: time.Date(2018, 1, 1, 7, 1, 3, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	if b[0].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup 20180101_070103 expected to be set as DELETE`)
	}
	if b[1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup 20200513_125539 expected to be set as DELETE`)
	}
}

func TestApplyPoliciesKeeping(t *testing.T) {
	fixedPoint := time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name: "keep",
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20180101_070103",
			Creation: time.Date(2018, 1, 1, 7, 1, 3, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	if b[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20180101_070103 expected to be set as KEEP`)
	}
	if b[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20200513_125539 expected to be set as KEEP`)
	}
}

func TestApplyPoliciesKeepingOfFuture(t *testing.T) {
	fixedPoint := time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name: "delete",
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20180101_070103",
			Creation: time.Date(2018, 1, 1, 7, 1, 3, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20221011_125539",
			Creation: time.Date(2022, 10, 11, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20221012_125539",
			Creation: time.Date(2022, 10, 12, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	if b[0].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup 20180101_070103 expected to be set as DELETE`)
	}
	if b[1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup 20200513_125539 expected to be set as DELETE`)
	}
	if b[2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup 20221011_125539 expected to be set as DELETE`)
	}
	if b[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup 20221012_125539 expected to be set as KEEP`)
	}
}

func TestApplyPoliciesDilute(t *testing.T) {
	fixedPoint := time.Date(2022, 10, 11, 0, 0, 0, 0, time.UTC)
	window := "4 days"
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window,
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20200501_125539",
			Creation: time.Date(2020, 5, 1, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200502_125539",
			Creation: time.Date(2020, 5, 2, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200503_125539",
			Creation: time.Date(2020, 5, 3, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200504_125539",
			Creation: time.Date(2020, 5, 4, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200505_125539",
			Creation: time.Date(2020, 5, 5, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200506_125539",
			Creation: time.Date(2020, 5, 6, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200507_125539",
			Creation: time.Date(2020, 5, 7, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200508_125539",
			Creation: time.Date(2020, 5, 8, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200509_125539",
			Creation: time.Date(2020, 5, 9, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200510_125539",
			Creation: time.Date(2020, 5, 10, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200511_125539",
			Creation: time.Date(2020, 5, 11, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200512_125539",
			Creation: time.Date(2020, 5, 12, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200514_125539",
			Creation: time.Date(2020, 5, 14, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200515_125539",
			Creation: time.Date(2020, 5, 15, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	var offset int = 0
	var s int = 4
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	offset = 1 * s
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	offset = 2 * s
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
}

func TestApplyPoliciesDiluteRepeatedlyOneDayAfterAnother(t *testing.T) {
	fixedPoint := time.Date(2022, 5, 16, 0, 0, 0, 0, time.UTC)
	fixedPointPlusOne := time.Date(2022, 5, 17, 0, 0, 0, 0, time.UTC)
	fixedPointPlusTwo := time.Date(2022, 5, 18, 0, 0, 0, 0, time.UTC)
	fixedPointPlusThree := time.Date(2022, 5, 19, 0, 0, 0, 0, time.UTC)
	fixedPointPlusFour := time.Date(2022, 5, 20, 0, 0, 0, 0, time.UTC)
	window := "4 days"
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window,
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20200501_125539",
			Creation: time.Date(2020, 5, 1, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200502_125539",
			Creation: time.Date(2020, 5, 2, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200503_125539",
			Creation: time.Date(2020, 5, 3, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200504_125539",
			Creation: time.Date(2020, 5, 4, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200505_125539",
			Creation: time.Date(2020, 5, 5, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200506_125539",
			Creation: time.Date(2020, 5, 6, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200507_125539",
			Creation: time.Date(2020, 5, 7, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200508_125539",
			Creation: time.Date(2020, 5, 8, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200509_125539",
			Creation: time.Date(2020, 5, 9, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200510_125539",
			Creation: time.Date(2020, 5, 10, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200511_125539",
			Creation: time.Date(2020, 5, 11, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200512_125539",
			Creation: time.Date(2020, 5, 12, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200514_125539",
			Creation: time.Date(2020, 5, 14, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200515_125539",
			Creation: time.Date(2020, 5, 15, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	var offset int = 0
	var s int = 4
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	offset = 1 * s
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	offset = 2 * s
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}

	// One day after
	b = *deleteDeletedBackupEntries(&b)
	ApplyPolicies(fixedPointPlusOne, p, &b)
	if b[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[0].Path)
	}
	if b[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[1].Path)
	}
	if b[2].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[2].Path)
	}
	if b[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[3].Path)
	}

	// Two days after
	b = *deleteDeletedBackupEntries(&b)
	ApplyPolicies(fixedPointPlusTwo, p, &b)
	if b[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[0].Path)
	}
	if b[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[1].Path)
	}
	if b[2].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[2].Path)
	}
	if b[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[3].Path)
	}

	// Three days after
	b = *deleteDeletedBackupEntries(&b)
	ApplyPolicies(fixedPointPlusThree, p, &b)
	if b[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[0].Path)
	}
	if b[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[1].Path)
	}
	if b[2].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[2].Path)
	}
	if b[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[3].Path)
	}

	// Four days after
	b = *deleteDeletedBackupEntries(&b)
	ApplyPolicies(fixedPointPlusFour, p, &b)
	if b[0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[0].Path)
	}
	if b[1].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[1].Path)
	}
	if b[2].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[2].Path)
	}
	if b[3].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[3].Path)
	}
}

func TestApplyPoliciesDiluteTwoDifferentWindows(t *testing.T) {
	fixedPoint := time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC)
	window := "3 days"
	window2 := "5 days"
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "6 days",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window,
			},
		},
		{
			From: "6 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window2,
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20200501_125539",
			Creation: time.Date(2020, 5, 1, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200502_125539",
			Creation: time.Date(2020, 5, 2, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200503_125539",
			Creation: time.Date(2020, 5, 3, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200504_125539",
			Creation: time.Date(2020, 5, 4, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200505_125539",
			Creation: time.Date(2020, 5, 5, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200506_125539",
			Creation: time.Date(2020, 5, 6, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200507_125539",
			Creation: time.Date(2020, 5, 7, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200508_125539",
			Creation: time.Date(2020, 5, 8, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200509_125539",
			Creation: time.Date(2020, 5, 9, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200510_125539",
			Creation: time.Date(2020, 5, 10, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200511_125539",
			Creation: time.Date(2020, 5, 11, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200512_125539",
			Creation: time.Date(2020, 5, 12, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200514_125539",
			Creation: time.Date(2020, 5, 14, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200515_125539",
			Creation: time.Date(2020, 5, 15, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	var offset int = 0
	var s int = 5
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	if b[offset+4].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+4].Path)
	}
	offset = 1 * s
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+3].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+3].Path)
	}
	offset = offset + 4
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	offset = offset + 3
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
}

func TestApplyPoliciesTwoDifferentDilutesAfterSomeTime(t *testing.T) {
	fixedPoint := time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC)
	fixedPoint2 := time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC)
	window := "3 days"
	window2 := "5 days"
	p := []definitions.Policy{
		{
			From: "0 days",
			To:   "14 days",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window,
			},
		},
		{
			From: "14 days",
			To:   "infinity",
			Strategy: definitions.Strategy{
				Name:   "dilute",
				Window: &window2,
			},
		},
	}
	b := []definitions.BackupDir{
		{
			Path:     "20200501_125539",
			Creation: time.Date(2020, 5, 1, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200502_125539",
			Creation: time.Date(2020, 5, 2, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200503_125539",
			Creation: time.Date(2020, 5, 3, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200504_125539",
			Creation: time.Date(2020, 5, 4, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200505_125539",
			Creation: time.Date(2020, 5, 5, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200506_125539",
			Creation: time.Date(2020, 5, 6, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200507_125539",
			Creation: time.Date(2020, 5, 7, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200508_125539",
			Creation: time.Date(2020, 5, 8, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200509_125539",
			Creation: time.Date(2020, 5, 9, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200510_125539",
			Creation: time.Date(2020, 5, 10, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200511_125539",
			Creation: time.Date(2020, 5, 11, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200512_125539",
			Creation: time.Date(2020, 5, 12, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200513_125539",
			Creation: time.Date(2020, 5, 13, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200514_125539",
			Creation: time.Date(2020, 5, 14, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
		{
			Path:     "20200515_125539",
			Creation: time.Date(2020, 5, 15, 12, 55, 39, 0, time.UTC),
			Outcome:  definitions.KeepBackup,
		},
	}
	ApplyPolicies(fixedPoint, p, &b)

	var offset int = 0
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}

	offset = 1
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+1].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	offset = offset + 3
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+1].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	offset = offset + 3
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+1].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	offset = offset + 3
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+1].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	offset = offset + 3
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+1].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}

	b = *deleteDeletedBackupEntries(&b)
	ApplyPolicies(fixedPoint2, p, &b)
	offset = 0
	if b[offset+0].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+0].Path)
	}
	if b[offset+1].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+1].Path)
	}
	if b[offset+2].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as DELETE`, b[offset+2].Path)
	}
	if b[offset+4].Outcome != definitions.DeleteBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+4].Path)
	}
	if b[offset+5].Outcome != definitions.KeepBackup {
		t.Fatalf(`Backup %s expected to be set as KEEP`, b[offset+5].Path)
	}
}

func deleteDeletedBackupEntries(backups *[]definitions.BackupDir) *[]definitions.BackupDir {

	var s []definitions.BackupDir
	for _, backup := range *backups {
		if backup.Outcome == definitions.KeepBackup {
			s = append(s, backup)
		}
	}
	return &s
}
