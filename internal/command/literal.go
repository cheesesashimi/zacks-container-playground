package command

// Constructs a command from a string slice.
type CommandLiteral []string

func (c CommandLiteral) Command() *Command {
	args := []Arg{}

	for _, item := range c[1:] {
		args = append(args, PositionalArg(item))
	}

	return NewCommand(c[0], args)
}

// Constructs a list of args from a string slice.
type ArgLiterals []string

func (a ArgLiterals) Arg() []Arg {
	args := []Arg{}

	for _, item := range a {
		args = append(args, PositionalArg(item))
	}

	return args
}
