package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	CommandType uint32
	Command     interface {
		Exec() error
		ReplaceParam(param string, val ...string) error
		Stage() string
		Type() CommandType
	}

	commandFlagSet struct {
		Output string `sqflag:"-o"`
	}

	command struct {
		args []string
		// paramPos is the index in args of the *value* provided for the parameter stored in the key
		paramPos map[string]int
		flags    commandFlagSet
	}
)

const (
	CommandTypeUnknown CommandType = iota
	CommandTypeCompile
	CommandTypeLink
)

func (t CommandType) String() string {
	switch t {
	case CommandTypeCompile:
		return "Compile"
	case CommandTypeLink:
		return "Link"
	default:
		return "Unknown"
	}
}

func NewCommand(args []string) command {
	cmd := command{
		args:     args,
		paramPos: make(map[string]int),
	}
	for pos, v := range args[1:] {
		cmd.paramPos[v] = pos + 1
	}

	parseFlags(&cmd.flags, args)

	return cmd
}

func (cmd *command) ReplaceParam(param string, val ...string) error {
	i, ok := cmd.paramPos[param]
	if !ok {
		return fmt.Errorf("param %s not found", param)
	}
	for _, v := range val {
		cmd.args[i] = v
		i++
	}
	return nil
}

func (cmd *command) Exec() error {
	c := exec.Command(cmd.args[0], cmd.args[1:]...)
	if c == nil {
		return fmt.Errorf("command couldn't build")
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (cmd *command) Stage() string {
	return filepath.Base(filepath.Dir(cmd.flags.Output))
}

func (cmd *command) Type() CommandType {
	return CommandTypeUnknown
}
