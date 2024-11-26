package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Command struct {
	args []Arg
}

func NewCommand(args []Arg) (*Command, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no args found")
	}

	if _, ok := args[0].(CommandName); !ok {
		return nil, fmt.Errorf("first arg must be CommandName, got: %T", args[0])
	}

	return &Command{
		args: args,
	}, nil
}

func (c *Command) String() string {
	return strings.Join(renderArgs(c.args), " ")
}

func (c *Command) Cmd() *exec.Cmd {
	return exec.Command(c.args[0].Arg()[0], renderArgs(c.args[1:])...)
}

type CommandName string

func (c CommandName) Arg() []string {
	return []string{string(c)}
}

type PositionalArg string

func (a PositionalArg) Arg() []string {
	return []string{string(a)}
}

type Arg interface {
	Arg() []string
}

type Subcommand struct {
	Name  string
	Flags []Flag
}

func (s *Subcommand) Arg() []string {
	out := []string{s.Name}
	return append(out, renderFlags(s.Flags)...)
}

type Flag interface {
	Arg() []string
}

type SingleSwitchFlag string

func (s SingleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", string(s))}
}

type DoubleSwitchFlag string

func (d DoubleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", string(d))}
}

type SingleValueFlag struct {
	Name  string
	Value string
}

func (s *SingleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", s.Name), s.Value}
}

type SingleEqualValueFlag struct {
	Name  string
	Value string
}

func (s *SingleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s=%s", s.Name, s.Value)}
}

type DoubleEqualValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s=%s", d.Name, d.Value)}
}

type DoubleValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", d.Name), d.Value}
}

func CombineArgs(arg []Arg) []string {
	return renderArgs(arg)
}

func renderArgs(args []Arg) []string {
	out := []string{}

	for _, arg := range args {
		out = append(out, arg.Arg()...)
	}

	return out
}

func renderFlags(flags []Flag) []string {
	out := []string{}

	for _, flag := range flags {
		out = append(out, flag.Arg()...)
	}

	return out
}
