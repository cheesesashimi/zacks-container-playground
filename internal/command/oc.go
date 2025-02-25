package command

import (
	"fmt"
)

// Represents an oc command. Note: This is not well fleshed-out.
type OcCommand interface {
	Command(string) *Command
}

type Login struct {
	Token  string
	Server string
}

func newOcCommand(kubeconfig string, args []Arg) *Command {
	return NewCommandWithEnv("oc", args, map[string]string{"KUBECONFIG": kubeconfig})
}

func (l *Login) Command(kubeconfig string) *Command {
	args := ArgLiterals{
		"login",
		"--token", l.Token,
		"--server", l.Server,
	}

	return newOcCommand(kubeconfig, args.Arg())
}

type RegistryLogin struct {
	To string
}

func (r *RegistryLogin) Command(kubeconfig string) *Command {
	return newOcCommand(kubeconfig, []Arg{
		PositionalArg("registry"),
		&Subcommand{
			Name: "login",
			Flags: mapToValFlags(map[string]string{
				"to": r.To,
			}),
		},
	})
}

type ImageExtract struct {
	Pullspec       string
	Path           string
	RegistryConfig string
}

func (i *ImageExtract) Command(kubeconfig string) *Command {
	args := []Arg{
		PositionalArg("image"),
		PositionalArg("extract"),
		PositionalArg(i.Pullspec),
	}

	args = append(args, flagsToArgs(mapToValFlags(map[string]string{
		"path":            i.Path,
		"registry-config": i.RegistryConfig,
	}))...)

	return newOcCommand(kubeconfig, args)
}

type ReleaseInfo struct {
	ReleasePullspec string
	Template        string
	JSON            bool
}

func (r *ReleaseInfo) Command(kubeconfig string) *Command {
	args := []Arg{
		PositionalArg("adm"),
		PositionalArg("release"),
		PositionalArg("info"),
	}

	if r.Template != "" {
		args = append(args, PositionalArg("-o=template="+r.Template))
	} else if r.JSON {
		args = append(args, &DoubleValueFlag{
			Name:  "output",
			Value: "json",
		})
	}

	args = append(args, PositionalArg(r.ReleasePullspec))

	return newOcCommand(kubeconfig, args)
}

type ReleaseExtract struct {
	RegistryConfig   string
	CommandToExtract string
	ReleasePullspec  string
	To               string
}

func (r *ReleaseExtract) Command(kubeconfig string) *Command {
	args := []Arg{
		PositionalArg("adm"),
		PositionalArg("release"),
		PositionalArg("extract"),
	}

	args = append(args, flagsToArgs(mapToValFlags(map[string]string{
		"registry-config": r.RegistryConfig,
		"command":         r.CommandToExtract,
		"to":              r.To,
	}))...)

	return newOcCommand(kubeconfig, args)
}

type DebugNode struct {
	Node             string
	WithChroot       bool
	ToNamespace      string
	CommandToExecute string
}

func (d *DebugNode) Command(kubeconfig string) *Command {
	args := []Arg{
		PositionalArg("debug"),
	}

	if d.ToNamespace != "" {
		args = append(args, &DoubleValueFlag{
			Name:  "to-namespace",
			Value: d.ToNamespace,
		})
	}

	args = append(args, PositionalArg(fmt.Sprintf("node/%s", d.Node)))
	args = append(args, PositionalArg("--"))

	if d.WithChroot {
		args = append(args, []Arg{PositionalArg("chroot"), PositionalArg("/host")}...)
	}

	args = append(args, []Arg{PositionalArg("/bin/bash"), PositionalArg("-c"), PositionalArg(d.CommandToExecute)}...)

	return newOcCommand(kubeconfig, args)
}
