package helpers

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"

	"drabek.cz/cli-utils/backupler/cmd/diluter/definitions"
	"drabek.cz/cli-utils/backupler/cmd/diluter/strategies"
)

const DEBUG_POLICIES_ALG = false

func ApplyPolicies(fixedPoint time.Time, p []definitions.Policy, backups *[]definitions.BackupDir) {
	// Go through all policies and apply them
	for _, policy := range p {
		// Keep policies can be skipped as they do not have any effect
		if policy.Strategy.Name == strategies.Keep {
			continue
		}
		// For dilute and delete policies iterate all backups
		var windowLastDate *time.Time = nil
		for i, backup := range *backups {
			// ... and consider only those failing into the right window
			backupCreationTrimmed := CloneDateTrimmed(backup.Creation)
			dayDiff := int64((fixedPoint.Sub(*backupCreationTrimmed).Hours() / 24))
			// Ignoring backups in future and today (after the midnight)
			if dayDiff < 0 {
				Log("KEEP %s - FUTURE\n", backup.Path)
				continue
			}
			// Ignore every backup outside the window
			policyFrom, err := ParseDays(policy.From, false)
			if err != nil {
				panic("policies should already be validate")
			}
			policyTo, err := ParseDays(policy.To, true)
			if err != nil {
				panic("policies should already be validate")
			}
			if dayDiff < policyFrom || dayDiff >= policyTo {
				Log("SKIP %s - %d OiI [%d, %d>", backup.Path, dayDiff, policyFrom, policyTo)
				continue
			}
			// Apply delete policy if relevant
			if policy.Strategy.Name == strategies.Delete {
				backup.Outcome = definitions.DeleteBackup
			}
			// Apply diluting policy if relevant
			if policy.Strategy.Name == strategies.Dilute {
				if windowLastDate == nil {
					backup.Outcome = definitions.KeepBackup
					windowLastDate = CloneDate(backup.Creation)
					Log("KEEP %s - FIRST IN WINDOW", backup.Path)
				} else {
					policyWindow, err := ParseDays(*policy.Strategy.Window, true)
					if err != nil {
						panic("policies should already be validate")
					}
					toLastKeptBackupDiff := int64((backup.Creation.Sub(*windowLastDate).Hours() / 24))
					if toLastKeptBackupDiff >= policyWindow {
						windowLastDate = CloneDate(backup.Creation)
						backup.Outcome = definitions.KeepBackup
						Log("KEEP %s - NEXT IN WINDOW", backup.Path)
					} else {
						backup.Outcome = definitions.DeleteBackup
						Log("KEEP %s - DILUTED IN WINDOW", backup.Path)
					}
				}
			}
			(*backups)[i] = backup
		}
	}
}

func ListBackups(directoryPath string, namingPattern string) ([]definitions.BackupDir, error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	output := []definitions.BackupDir{}
	for _, f := range files {
		absoluteFilePath, err := filepath.Abs(filepath.Join(directoryPath, f.Name()))
		if err != nil {
			return nil, err
		}
		creation, err := ParseDate(f.Name(), namingPattern)
		if err != nil {
			continue
		}
		backup := definitions.BackupDir{
			Path:     absoluteFilePath,
			Creation: *creation,
			Outcome:  definitions.KeepBackup,
		}
		output = append(output, backup)
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Creation.Before(output[j].Creation)
	})
	return output, nil
}
