package command

type RpmOstreeInstall struct {
	Packages []string
}

func (r *RpmOstreeInstall) Command() *Command {
	args := []Arg{
		&Subcommand{
			Name: "install",
		},
	}

	for _, pkg := range r.Packages {
		args = append(args, PositionalArg(pkg))
	}

	return NewCommand("rpm-ostree", args)
}
