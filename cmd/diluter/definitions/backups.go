package definitions

import "time"

type Outcome string

const (
	KeepBackup   Outcome = "KEEP"
	DeleteBackup Outcome = "DELETE"
)

type BackupDir struct {
	Path     string
	Creation time.Time
	Outcome  Outcome
}
