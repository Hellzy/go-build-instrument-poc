// Copyright (c) 2016 - 2019 Sqreen. All Rights Reserved.
// Please refer to our terms for more information:
// https://www.sqreen.io/terms.html

package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type compileFlagSet struct {
	Package   string `ddflag:"-p"`
	ImportCfg string `ddflag:"-importcfg"`
	Output    string `ddflag:"-o"`
}

type compileCommand struct {
	command
}

func (cmd *compileCommand) Type() CommandType { return CommandTypeCompile }

func (cmd *compileCommand) Inject(i Injector) {
	i.InjectCompile(cmd)
}

func (f *compileFlagSet) IsValid() bool {
	return f.Package != "" && f.Output != "" && f.ImportCfg != ""
}

func (f *compileFlagSet) String() string {
	return fmt.Sprintf("-p=%q -o=%q -importcfg=%q", f.Package, f.Output, f.ImportCfg)
}

func parseCompileCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, errors.New("unexpected number of command arguments")
	}
	flags := &compileFlagSet{}
	parseFlags(flags, args[1:])
	return &compileCommand{NewCommand(args)}, nil
}

// Go builds in provided dir and returns the work directory's true path
func goBuild(dir string, args ...string) (string, error) {
	args = append([]string{"build", "-work", "-a", "-p", "1"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		return "", err
	}

	// Extract work dir from output
	wDir := strings.Split(string(out), "=")[1]
	wDir = strings.TrimSuffix(wDir, "\n")

	return wDir, nil
}
