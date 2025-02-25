package command

import "fmt"

// Represents a label given to podman or buildah build
type Label struct {
	Name  string
	Value string
}

func (l *Label) flag() Flag {
	return &DoubleValueFlag{
		Name:  "label",
		Value: fmt.Sprintf("%s=%s", l.Name, l.Value),
	}
}

// Represents a build arg passed to podman or buildah build
type BuildArg struct {
	Name  string
	Value string
}

func (b *BuildArg) flag() Flag {
	return &DoubleValueFlag{
		Name:  "build-arg",
		Value: fmt.Sprintf("%s=%s", b.Name, b.Value),
	}
}

// Represents a volume passed to podman run
type Volume struct {
	HostPath      string
	ContainerPath string
	Opts          string
}

func (p *Volume) flag() Flag {
	return &DoubleValueFlag{
		Name:  "volume",
		Value: p.render(),
	}
}

func (p *Volume) render() string {
	out := fmt.Sprintf("%s:%s", p.HostPath, p.ContainerPath)
	if p.Opts != "" {
		out = fmt.Sprintf("%s:%s", out, p.Opts)
	}

	return out
}

// Represents an env var passed into podman run
type PodmanEnv struct {
	Name  string
	Value string
}

func (p *PodmanEnv) flag() Flag {
	return &DoubleValueFlag{
		Name:  "env",
		Value: fmt.Sprintf("%s=%s", p.Name, p.Value),
	}
}

// Represents a podman tag command
type PodmanTag struct {
	Image      string
	TargetName string
}

func (p *PodmanTag) Command() *Command {
	return NewCommand("podman", []Arg{
		&Subcommand{
			Name: "tag",
		},
		PositionalArg(p.Image),
		PositionalArg(p.TargetName),
	})
}

// Represents a podman push command
type PodmanPush struct {
	Authfile  string
	Format    string
	Image     string
	TLSVerify *bool
}

func (p *PodmanPush) Command() *Command {
	pushFlags := mapToValFlags(map[string]string{
		"authfile": p.Authfile,
		"format":   p.Format,
	})

	pushFlags = append(pushFlags, mapToOptSwitchFlags(map[string]*bool{
		"tls-verify": p.TLSVerify,
	})...)

	return NewCommand("podman", []Arg{
		&Subcommand{
			Name:  "push",
			Flags: pushFlags,
		},
		PositionalArg(p.Image),
	})
}

// Represents a podman build command
type PodmanBuild struct {
	BuildContext string
	Tag          string
	Target       string
	BuildArgs    []BuildArg
	Labels       []Label
	File         string
}

func (p *PodmanBuild) Command() *Command {
	buildFlags := mapToValFlags(map[string]string{
		"tag":    p.Tag,
		"target": p.Target,
		"file":   p.File,
	})

	for _, buildArg := range p.BuildArgs {
		buildFlags = append(buildFlags, buildArg.flag())
	}

	for _, label := range p.Labels {
		buildFlags = append(buildFlags, label.flag())
	}

	buildCtx := p.BuildContext
	if buildCtx == "" {
		buildCtx = "."
	}

	return NewCommand("podman", []Arg{
		&Subcommand{
			Name:  "build",
			Flags: buildFlags,
		},

		PositionalArg(buildCtx),
	})
}

// Represents a podman run command.
type PodmanRun struct {
	Interactive     bool
	Tty             bool
	Remove          bool
	Detach          bool
	Name            string
	AdditionalFlags []Flag
	Volumes         []Volume
	Env             []PodmanEnv
	Workdir         string
	Entrypoint      string
	ImageOpts       []Arg
	Image           string
}

func (p *PodmanRun) Command() *Command {
	runFlags := mapToSwitchFlags(map[string]bool{
		"interactive": p.Interactive,
		"tty":         p.Tty,
		"rm":          p.Remove,
		"detach":      p.Detach,
	})

	runFlags = append(runFlags, mapToValFlags(map[string]string{
		"name":       p.Name,
		"workdir":    p.Workdir,
		"entrypoint": p.Entrypoint,
	})...)

	for _, volume := range p.Volumes {
		runFlags = append(runFlags, volume.flag())
	}

	for _, env := range p.Env {
		runFlags = append(runFlags, env.flag())
	}

	runFlags = append(runFlags, p.AdditionalFlags...)

	return NewCommand("podman", append([]Arg{
		&Subcommand{
			Name:  "run",
			Flags: runFlags,
		},
		PositionalArg(p.Image),
	}, p.ImageOpts...))
}

// Represents an unimplemented podman pull command
type PodmanPull struct {
	Image string
}
