package command

type CommandLiteral []string

func (c CommandLiteral) Command() *Command {
	args := []Arg{}

	for _, item := range c[1:] {
		args = append(args, PositionalArg(item))
	}

	return NewCommand(c[0], args)
}

type ArgLiterals []string

func (a ArgLiterals) Arg() []Arg {
	args := []Arg{}

	for _, item := range a {
		args = append(args, PositionalArg(item))
	}

	return args
}
