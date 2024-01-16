package main

import (
	"errors"
	"path/filepath"
)

type linkFlagSet struct {
	ImportCfg string `sqflag:"-importcfg"`
	Output    string `sqflag:"-o"`
}

type linkCommand struct {
	command
}

func (cmd *linkCommand) Type() CommandType {
	return CommandTypeLink
}

func (cmd *linkCommand) Stage() string {
	return filepath.Base(filepath.Dir(filepath.Dir(cmd.flags.Output)))
}

func parseLinkCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, errors.New("unexpected number of command arguments")
	}
	flags := &linkFlagSet{}
	parseFlags(flags, args[1:])
	return makeLinkCommandExecutionFunc(flags, args), nil
}

func makeLinkCommandExecutionFunc(flags *linkFlagSet, args []string) Command {
	//buildDir := filepath.Dir(flags.Output)

	return &linkCommand{command: NewCommand(args)}
}
