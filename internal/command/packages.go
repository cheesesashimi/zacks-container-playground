package command

// Represents a call to dnf install.
type DnfInstall struct {
	Yes      bool
	Packages []string
}

func (d *DnfInstall) Command() *Command {
	return getGenericInstallerCommand("dnf", d.Yes, d.Packages)
}

// Represents a call to yum install.
type YumInstall struct {
	Yes      bool
	Packages []string
}

func (y *YumInstall) Command() *Command {
	return getGenericInstallerCommand("yum", y.Yes, y.Packages)
}

// Represents a call to apt-get install
type AptGetInstall struct {
	Yes      bool
	Packages []string
}

func (a *AptGetInstall) Command() *Command {
	return getGenericInstallerCommand("apt-get", a.Yes, a.Packages)
}

// Represents a call to apt-get upgrade
type AptGetUpdate struct{}

func (a *AptGetUpdate) Command() *Command {
	cmd := &CommandLiteral{"apt-get", "update"}
	return cmd.Command()
}

func getGenericInstallerCommand(installer string, yes bool, packages []string) *Command {
	args := []Arg{
		&Subcommand{
			Name: "install",
		},
	}

	if yes {
		y := SingleSwitchFlag("y")
		args = append(args, &y)
	}

	args = append(args, itemsToPositionalArgs(packages)...)

	return NewCommand(installer, args)
}
