package mocker

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"drabek.digital/cli-utils/backupler/cmd/diluter/helpers"
	"github.com/viant/toolbox"
)

type Mocker struct {
	directoryPath      string
	mocksConfiguration string
	from               *time.Time
	to                 *time.Time
	format             *string
}

func (mocker *Mocker) GetDirectoryPath() string {
	return mocker.directoryPath
}

func NewMocker(directoryPath string, mocksConfiguration string) *Mocker {
	mocker := Mocker{}

	mocker.directoryPath = directoryPath
	mocker.mocksConfiguration = mocksConfiguration

	return &mocker
}

func (mocker *Mocker) Execute() (bool, error) {
	if mocker.from == nil || mocker.to == nil {
		panic("expected interval to be populated")
	}

	now := *mocker.from
	layout := toolbox.DateFormatToLayout(*mocker.format)
	for now.Before(*mocker.to) {
		dir := now.Format(layout)
		fmt.Fprintf(os.Stdout, "Creating directory `%s/%s`\n", mocker.directoryPath, dir)
		if err := os.Mkdir(mocker.directoryPath+"/"+dir, os.ModePerm); err != nil {
			return false, err
		}

		now = now.Add(24 * time.Hour)
	}

	return true, nil
}

func (mocker *Mocker) ParseAndValidateConfig() (bool, error) {
	fileInfo, err := os.Stat(mocker.directoryPath)
	if os.IsNotExist(err) {
		return false, err
	}
	if !fileInfo.IsDir() {
		return false, errors.New("given directory is not a directory")
	}
	parts := strings.Split(mocker.mocksConfiguration, ":")
	if len(parts) != 3 {
		return false, errors.New("mocking interval is invalid")
	}
	dateFormat := "yyyy-MM-dd"
	from, err := helpers.ParseDate(parts[0], dateFormat)
	if err != nil {
		return false, err
	}
	to, err := helpers.ParseDate(parts[1], dateFormat)
	if err != nil {
		return false, err
	}
	mocker.from = from
	mocker.to = to
	mocker.format = &parts[2]

	return true, nil
}
