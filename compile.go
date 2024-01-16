// Copyright (c) 2016 - 2019 Sqreen. All Rights Reserved.
// Please refer to our terms for more information:
// https://www.sqreen.io/terms.html

package main

import (
	"errors"
	"fmt"
)

type compileFlagSet struct {
	Package   string `sqflag:"-p"`
	ImportCfg string `sqflag:"-importcfg"`
	Output    string `sqflag:"-o"`
}

type compileCommand struct {
	command
}

func (cmd *compileCommand) Type() CommandType { return CommandTypeCompile }

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
	return makeCompileCommandExecutionFunc(flags, args), nil
}

func makeCompileCommandExecutionFunc(flags *compileFlagSet, args []string) Command {
	//buildDir := filepath.Dir(flags.Output)
	//stage := filepath.Base(buildDir)

	return &compileCommand{NewCommand(args)}
	//	return func() ([]string, error) {
	//		if !flags.IsValid() {
	//			// Skip when the required set of flags is not valid.
	//			log.Printf("nothing to do (%s)\n", flags)
	//			return nil, nil
	//		}
	//
	//		return args, nil
}
