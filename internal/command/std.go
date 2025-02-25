package command

import "fmt"

type Echo struct {
	Content    string
	Escape     bool
	RedirectTo string
}

func (e *Echo) Command() *Command {
	args := flagsToArgs(mapToSingleSwitchFlag(map[string]bool{
		"e": e.Escape,
	}))

	args = append(args, PositionalArg(fmt.Sprintf("%q", e.Content)))

	if e.RedirectTo != "" {
		// Figure out a better way to do this because exec.Command won't run this.
		args = append(args, itemsToPositionalArgs([]string{">", e.RedirectTo})...)
	}

	return NewCommand("echo", args)
}

type Delete struct {
	Path      string
	Recursive bool
	Verbose   bool
}

func (d *Delete) Command() *Command {
	args := flagsToArgs(mapToSingleSwitchFlag(map[string]bool{
		"r": d.Recursive,
		"v": d.Verbose,
	}))

	args = append(args, PositionalArg(d.Path))

	return NewCommand("rm", args)
}

type Chmod struct {
	Path      string
	Mode      string
	Recursive bool
}

func (c *Chmod) Command() *Command {
	args := flagsToArgs(mapToSingleSwitchFlag(map[string]bool{
		"R": c.Recursive,
	}))

	args = append(args, itemsToPositionalArgs([]string{
		c.Mode,
		c.Path,
	})...)

	return NewCommand("chmod", args)
}

func mapToValFlags(valOpts map[string]string) []Flag {
	out := []Flag{}

	for name, val := range valOpts {
		if val != "" {
			out = append(out, &DoubleValueFlag{
				Name:  name,
				Value: val,
			})
		}
	}

	return out
}

func mapToSwitchFlags(switchOpts map[string]bool) []Flag {
	out := []Flag{}

	for name, isSet := range switchOpts {
		if isSet {
			out = append(out, DoubleSwitchFlag(name))
		}
	}

	return out
}

func mapToOptSwitchFlags(optSwitchOpts map[string]*bool) []Flag {
	out := []Flag{}

	for name, val := range optSwitchOpts {
		if val != nil {
			out = append(out, &DoubleEqualValueFlag{
				Name:  name,
				Value: fmt.Sprintf("%v", *val),
			})
		}
	}

	return out
}

func mapToSingleSwitchFlag(switchOpts map[string]bool) []Flag {
	out := []Flag{}

	for name, isSet := range switchOpts {
		if isSet {
			out = append(out, SingleSwitchFlag(name))
		}
	}

	return out
}

func flagsToArgs(flags []Flag) []Arg {
	args := []Arg{}

	for _, flag := range flags {
		args = append(args, flag)
	}

	return args
}

func itemsToPositionalArgs(items []string) []Arg {
	args := []Arg{}

	for _, item := range items {
		args = append(args, PositionalArg(item))
	}

	return args
}
