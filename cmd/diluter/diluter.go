package diluter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"drabek.cz/cli-utils/backupler/cmd/diluter/config"
	"drabek.cz/cli-utils/backupler/cmd/diluter/definitions"
	"drabek.cz/cli-utils/backupler/cmd/diluter/helpers"
	"gopkg.in/yaml.v3"
)

type Diluter struct {
	directoryPath string
	configPath    string
	testRun       bool
	approval      bool
	config        definitions.Config
}

func (diluter *Diluter) GetDirectoryPath() string {
	return diluter.directoryPath
}

func (diluter *Diluter) GetConfigPath() string {
	return diluter.configPath
}

func (diluter *Diluter) ValidateDirectory() bool {
	_, err := os.Stat(diluter.directoryPath)
	return !os.IsNotExist(err)
}

func (diluter *Diluter) ParseAndValidateConfig() (bool, error) {
	_, err := os.Stat(diluter.configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	content, err := ioutil.ReadFile(diluter.configPath)

	if err != nil {
		return false, err
	}

	var c definitions.Config
	if err := yaml.Unmarshal(content, &c); err != nil {
		return false, err
	}
	err = config.ValidateConfig(c)
	if err != nil {
		return false, err
	}

	diluter.config = c
	return true, nil
}

func NewDiluter(directoryPath string, configPath string, testRun bool, approval bool) *Diluter {
	diluter := Diluter{}

	diluter.directoryPath = directoryPath
	diluter.configPath = configPath
	diluter.testRun = testRun
	diluter.approval = approval

	return &diluter
}

func (diluter *Diluter) Execute() (bool, error) {
	// Load backups
	backups, err := helpers.ListBackups(diluter.directoryPath, diluter.config.Backup.Naming)
	if err != nil {
		return false, err
	}

	// Create fixed point for evaluation
	fixedPoint := time.Now().Round(0)
	fixedPoint = time.Date(fixedPoint.Year(), fixedPoint.Month(), fixedPoint.Day(), 0, 0, 0, 0, time.UTC)

	// Apply policies
	helpers.ApplyPolicies(fixedPoint, diluter.config.Policy, &backups)

	// Print outcome
	toBeDeleted := 0
	fmt.Printf("Outcome\n")
	fmt.Printf("---\n")
	for _, backup := range backups {
		fmt.Printf("%s\t\t%s\n", backup.Path, backup.Outcome)
		if backup.Outcome == definitions.DeleteBackup {
			toBeDeleted++
		}
	}
	fmt.Printf("---\n")
	fmt.Printf("backups for deletion: %d\n", toBeDeleted)

	if diluter.testRun {
		fmt.Printf("---\n")
		fmt.Printf("test run requested... nothing performed\n")
		return true, nil
	}
	if diluter.approval {
		confirmation := helpers.StringPrompt("do you want to proceed [y/N]?")
		if strings.ToUpper(confirmation) != "Y" {
			fmt.Printf("---\n")
			fmt.Printf("not approved... nothing performed\n")
			return true, nil
		}
		fmt.Printf("---\n")
		fmt.Printf("deletion in\t\t3... ")
		time.Sleep(1 * time.Second)
		fmt.Printf("2... ")
		time.Sleep(1 * time.Second)
		fmt.Printf("1... ")
		time.Sleep(1 * time.Second)
		fmt.Printf("NOW!")
	}

	// The real deletion
	deleted := 0
	for _, backup := range backups {
		if backup.Outcome == definitions.DeleteBackup {
			fmt.Printf("%s\t\t", backup.Path)
			err := os.RemoveAll(backup.Path)
			if err != nil {
				fmt.Printf("%s\n", "FAILED")
				continue
			}
			fmt.Printf("%s\n", backup.Outcome)
			deleted = deleted + 1

		}
	}
	fmt.Printf("---\n")
	fmt.Printf("backups for deletion: %d\n", toBeDeleted)
	fmt.Printf("backups deleted: %d\n", deleted)

	return true, nil
}
