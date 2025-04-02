package genflag

import (
	"fmt"
	"strings"
)

var _ FlagMarshaller = PodmanEnvOpt{}
var _ FlagMarshaller = &PodmanEnvOpt{}

type PodmanEnvOpt struct {
	Name  string
	Value string
}

func (p PodmanEnvOpt) MarshalFlags() ([]Flag, error) {
	val := fmt.Sprintf("%s=%s", p.Name, p.Value)
	f, err := NewStringFlag("env", val, Quoted)
	return []Flag{f}, err
}

type PodmanRun struct {
	Interactive bool `genflag:""`
	TTY         bool `genflag:""`
	Remove      bool `genflag:"rm"`
	Env         []PodmanEnvOpt
	Entrypoint  string `genflag:""`
	Image       string
}

func (p PodmanRun) PodmanCmd() (string, error) {
	flags, err := MarshalFlags(p)
	if err != nil {
		return "", err
	}

	out := []string{"podman", "run"}

	for _, flag := range flags {
		segmented, err := flag.Segmented()
		if err != nil {
			return "", err
		}

		out = append(out, segmented...)
	}

	out = append(out, p.Image)

	return strings.Join(out, " "), nil
}

type BuildArg struct {
	Argument string
	Value    string
}

func (b BuildArg) String() string {
	return fmt.Sprintf("%s=%s", b.Argument, b.Value)
}

type PodmanBuild struct {
	Tag            string         `genflag:""`
	BuildArg       []BuildArg     `genflag:"build-arg,quoted"`
	File           string         `genflag:""`
	Env            []PodmanEnvOpt `genflag:""`
	DefaultContext string
}

func (p PodmanBuild) PodmanCmd() (string, error) {
	flags, err := MarshalFlags(p)
	if err != nil {
		return "", err
	}

	out := []string{"podman", "build"}

	for _, flag := range flags {
		segmented, err := flag.Segmented()
		if err != nil {
			return "", err
		}

		out = append(out, segmented...)
	}

	defCtx := p.DefaultContext
	if defCtx == "" {
		defCtx = "."
	}
	out = append(out, defCtx)

	return strings.Join(out, " "), nil
}
