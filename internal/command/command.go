package command

import (
	"fmt"
	"os/exec"
	"strings"
)

// Basic command struct which holds a list of arguments and environment
// variables.
type Command struct {
	args []Arg
	env  map[string]string
}

func NewCommand(name string, args []Arg) *Command {
	return &Command{
		args: append([]Arg{PositionalArg(name)}, args...),
	}
}

func (c *Command) envVars() []string {
	envVars := []string{}

	for name, val := range c.env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", name, val))
	}

	return envVars
}

// Emits a string representation of the command including each environment
// variable.
func (c *Command) String() string {
	out := append(c.envVars(), renderArgs(c.args)...)
	return strings.Join(out, " ")
}

// Emits an instantiated exec.Cmd instance ready for execution.
func (c *Command) Cmd() *exec.Cmd {
	cmd := exec.Command(c.args[0].Arg()[0], renderArgs(c.args[1:])...)

	if c.env == nil {
		return cmd
	}

	cmd.Env = append(cmd.Env, c.envVars()...)
	return cmd
}

func NewCommandWithEnv(name string, args []Arg, env map[string]string) *Command {
	cmd := NewCommand(name, args)
	cmd.env = env
	return cmd
}

// Represents a positional argument to give to a command.
type PositionalArg string

func (a PositionalArg) Arg() []string {
	return []string{string(a)}
}

// All arguments must implement this interfeace.
type Arg interface {
	Arg() []string
}

// Represents a subcommand such as the get in $ oc get.
type Subcommand struct {
	Name  string
	Flags []Flag
}

func (s *Subcommand) Arg() []string {
	out := []string{s.Name}
	return append(out, renderFlags(s.Flags)...)
}

// Represents a flag argument such as --full-name, -single, or -s.
type Flag interface {
	Arg() []string
}

// Represents a single-switch flag like -s.
type SingleSwitchFlag string

func (s SingleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", string(s))}
}

// Represents a double-switch flag like --hello.
type DoubleSwitchFlag string

func (d DoubleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", string(d))}
}

// Represents a single-value flag such as "-key val"
type SingleValueFlag struct {
	Name  string
	Value string
}

func (s *SingleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", s.Name), s.Value}
}

// Represents a single-value equal flag such as "-key=val"
type SingleEqualValueFlag struct {
	Name  string
	Value string
}

func (s *SingleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s=%s", s.Name, s.Value)}
}

// Represents a double-flag with an equal such as "--key=val"
type DoubleEqualValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s=%s", d.Name, d.Value)}
}

// Represents a double-flag with no equal such as "--key val"
type DoubleValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", d.Name), d.Value}
}

// Combines and renders arguments into a string array.
func CombineArgs(arg []Arg) []string {
	return renderArgs(arg)
}

// Renders a list of args into a string array.
func renderArgs(args []Arg) []string {
	out := []string{}

	for _, arg := range args {
		out = append(out, arg.Arg()...)
	}

	return out
}

// Renders a list of flags into a string array.
func renderFlags(flags []Flag) []string {
	out := []string{}

	for _, flag := range flags {
		out = append(out, flag.Arg()...)
	}

	return out
}
