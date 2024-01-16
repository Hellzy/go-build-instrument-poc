// Copyright (c) 2016 - 2019 Sqreen. All Rights Reserved.
// Please refer to our terms for more information:
// https://www.sqreen.io/terms.html

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var globalFlags instrumentationToolFlagSet

func main() {
	log.SetFlags(0)
	log.SetPrefix("[instrumentation] | ")
	log.SetOutput(os.Stderr)

	args := os.Args[1:]
	cmdIdPos := parseFlagsUntilFirstNonOptionArg(&globalFlags, args)

	if cmdIdPos == -1 {
		log.Println("unexpected arguments")
		return
	}
	cmd, err := parseCommand(&globalFlags, args[cmdIdPos:])
	if err != nil || globalFlags.Help {
		log.Println(err)
		printUsage()
		os.Exit(1)
	}

	log.Printf("cmd: %s\n", cmd.Type().String())
	log.Printf("stage: %s\n", cmd.Stage())

	var logs strings.Builder
	if !globalFlags.Verbose {
		// Save the logs to show them in case of instrumentation error
		log.SetOutput(&logs)
	}

	err = cmd.Exec()
	var exitErr *exec.ExitError
	if err != nil {
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalln(err)
		}
	}
	os.Exit(0)
}

func printUsage() {
	const usageFormat = `Usage: go {build,install,get,test} -a -toolexec '%s' PACKAGES...

Datadog's instrumentation tool for Go v%s. Description to be provided in the future.

Options:
        -h
                Print this usage message.
`
	_, _ = fmt.Fprintf(os.Stderr, usageFormat, os.Args[0], "v0.0.0-poc.1")
	os.Exit(2)
}
