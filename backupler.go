package main

import (
	"fmt"
	"os"

	"drabek.cz/cli-utils/backupler/cmd/diluter"
	"drabek.cz/cli-utils/backupler/cmd/mocker"
	"github.com/pborman/getopt/v2"
)

func main() {
	optionConfigFilePathPtr := getopt.StringLong("config", 'c', ".backupler.yaml", "filepath to the backupler.yaml file")
	optionTestRunPtr := getopt.BoolLong("test-run", 't', "changes will be only printed")
	optionApprovalPtr := getopt.BoolLong("approval", 'a', "request approval before changes are made")
	optionMockDirsPtr := getopt.StringLong("mock-dirs", 'm', "", "mocks directories between given dates `2020-01-01:2022-12-31:yyyyMMdd_HHmmss` ")
	optHelpPtr := getopt.BoolLong("help", 0, "Help")

	getopt.SetParameters("DIRECTORY")

	getopt.SetUsage(func() {
		fmt.Print("Backupler is a tool for removing some of the historic backups to keep desired history of backups and save space.\n\n\n")
		getopt.PrintUsage(os.Stderr)
	})

	getopt.Parse()

	if *optHelpPtr {
		getopt.Usage()
		os.Exit(0)
	}

	args := getopt.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "missing DIRECTORY argument")
		os.Exit(10)
	}

	if optionMockDirsPtr != nil && *optionMockDirsPtr != "" {
		fmt.Fprintf(os.Stdout, "Running mocking directories in `%s` with configuration `%s`\n\n", args[0], *optionMockDirsPtr)
		mockerPtr := mocker.NewMocker(args[0], *optionMockDirsPtr)
		success, err := mockerPtr.ParseAndValidateConfig()
		if !success {
			fmt.Fprintln(os.Stderr, "config `"+mockerPtr.GetDirectoryPath()+"` is not valid:", err)
			os.Exit(21)
		}
		success, err = mockerPtr.Execute()
		if !success || err != nil {
			fmt.Fprintln(os.Stderr, "execution has failed:", err)
			os.Exit(22)
		}
		os.Exit(0)
	} else if *optionApprovalPtr && *optionTestRunPtr {
		fmt.Fprintf(os.Stdout, "Running test run of backupler with manual approval in `%s` with config `%s`!\n\n", args[0], *optionConfigFilePathPtr)
	} else if *optionApprovalPtr {
		fmt.Fprintf(os.Stdout, "Running backupler with manual approval in `%s` with config `%s`!\n\n", args[0], *optionConfigFilePathPtr)
	} else if *optionTestRunPtr {
		fmt.Fprintf(os.Stdout, "Running test run of backupler in `%s` with config `%s`\n\n!", args[0], *optionConfigFilePathPtr)
	} else {
		fmt.Fprintf(os.Stdout, "Running backupler in `%s` with config `%s`!\n\n", args[0], *optionConfigFilePathPtr)
	}

	diluterPtr := diluter.NewDiluter(args[0], *optionConfigFilePathPtr, *optionTestRunPtr, *optionApprovalPtr)

	if !diluterPtr.ValidateDirectory() {
		fmt.Fprintln(os.Stderr, "directory `"+diluterPtr.GetDirectoryPath()+"` either does not exist or is not readable")
		os.Exit(11)
	}
	success, err := diluterPtr.ParseAndValidateConfig()
	if !success && err == nil {
		fmt.Fprintln(os.Stderr, "config `"+diluterPtr.GetConfigPath()+"` either does not exist or is not readable")
		os.Exit(12)
	}
	if !success && err != nil {
		fmt.Fprintln(os.Stderr, "config `"+diluterPtr.GetConfigPath()+"` is not valid:", err)
		os.Exit(13)
	}

	success, err = diluterPtr.Execute()
	if !success || err != nil {
		fmt.Fprintln(os.Stderr, "execution has failed:", err)
		os.Exit(20)
	}
	os.Exit(0)
}
