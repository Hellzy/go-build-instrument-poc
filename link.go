package main

import (
	"errors"
	"path/filepath"
)

type linkFlagSet struct {
	ImportCfg string `ddflag:"-importcfg"`
	Output    string `ddflag:"-o"`
}

type linkCommand struct {
	command
	flags linkFlagSet
}

func (cmd *linkCommand) Inject(i Injector) {
	i.InjectLink(cmd)
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
	return &linkCommand{command: NewCommand(args), flags: *flags}, nil
}
