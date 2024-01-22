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

	appsecInjector := PackageInjector{
		importPath: "gopkg.in/DataDog/dd-trace-go.v1/internal/appsec",
		pkgDir:     "/Users/francois.mazeau/go/src/github.com/DataDog/dd-trace-go/internal/appsec",
	}

	appsecMainSwapper := GoFileSwapper{
		swapMap: map[string]string{"/Users/francois.mazeau/go/src/github.com/DataDog/instrumentation/main.go": "/Users/francois.mazeau/go/src/github.com/DataDog/instrumentation/.customMain/main.go"},
	}

	cmd.Inject(&appsecMainSwapper)
	cmd.Inject(&appsecInjector)

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

func exitIfError(err error) {
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

type GoFileSwapper struct {
	// Key: file to replace
	// Value: file to replace with
	swapMap map[string]string
}

func (s *GoFileSwapper) InjectCompile(cmd *compileCommand) {
	if cmd.Stage() != "b001" {
		return
	}
	log.Printf("[%s] Replacing Go files", cmd.Stage())

	for old, new := range s.swapMap {
		if err := cmd.ReplaceParam(old, new); err != nil {
			log.Printf("couldn't replace param: %v", err)
		} else {
			log.Printf("====> Replacing %s by %s", old, new)
		}
	}
}
func (s *GoFileSwapper) InjectLink(*linkCommand) {} // nothing to do at link time
