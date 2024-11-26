package main

import "fmt"

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

type PodmanEnv struct {
	Name  string
	Value string
}

func (p *PodmanEnv) flag() Flag {
	return &DoubleValueFlag{
		Name:  "env",
		Value: p.render(),
	}
}

func (p *PodmanEnv) render() string {
	return fmt.Sprintf("%s=%s", p.Name, p.Value)
}

type PodmanTag struct {
	Image      string
	TargetName string
}

func (p *PodmanTag) Args() []Arg {
	return []Arg{
		CommandName("podman"),
		&Subcommand{
			Name: "tag",
		},
		PositionalArg(p.Image),
		PositionalArg(p.TargetName),
	}
}

type PodmanPush struct {
	Authfile  string
	Format    string
	Image     string
	TLSVerify *bool
}

func (p *PodmanPush) Args() []Arg {
	pushFlags := mapToValFlags(map[string]string{
		"authfile": p.Authfile,
		"format":   p.Format,
	})

	pushFlags = append(pushFlags, mapToOptSwitchFlags(map[string]*bool{
		"tls-verify": p.TLSVerify,
	})...)

	return []Arg{
		CommandName("podman"),
		&Subcommand{
			Name:  "push",
			Flags: pushFlags,
		},
		PositionalArg(p.Image),
	}
}

type PodmanBuild struct {
	BuildContext string
	Tag          string
	Target       string
	BuildArgs    []BuildArg
	Labels       []Label
	File         string
}

func (p *PodmanBuild) Args() []Arg {
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

	return []Arg{
		CommandName("podman"),
		&Subcommand{
			Name:  "build",
			Flags: buildFlags,
		},
		PositionalArg(buildCtx),
	}
}

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

func (p *PodmanRun) Args() []Arg {
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

	out := []Arg{
		CommandName("podman"),
		&Subcommand{
			Name:  "run",
			Flags: runFlags,
		},
		PositionalArg(p.Image),
	}

	return append(out, p.ImageOpts...)
}

type PodmanPull struct {
	Image string
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
