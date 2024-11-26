package main

import (
	"fmt"
	"strings"
)

type ContainerfileSteps []ContainerfileStep

func (c ContainerfileSteps) Containerfile() string {
	out := &strings.Builder{}

	for _, step := range c {
		fmt.Fprintln(out, step.Line())
	}

	return out.String()
}

type ContainerfileStep interface {
	Line() string
}

type FromStep struct {
	Image string
	As    string
}

func (f *FromStep) Line() string {
	from := fmt.Sprintf("FROM %s", f.Image)
	if f.As == "" {
		return from
	}

	return fmt.Sprintf("%s AS %s", from, f.As)
}

type RunStep struct {
	Flags   []string
	Command string
}

func (r *RunStep) Line() string {
	if len(r.Flags) == 0 {
		return fmt.Sprintf("RUN %s", r.Command)
	}

	out := []string{"RUN"}
	out = append(out, r.Flags...)
	out = append(out, r.Command)

	return strings.Join(out, " ")
}

type MultiRunStep struct {
	Flags    []string
	Commands []string
}

func (m *MultiRunStep) Line() string {
	r := RunStep{
		Flags:   m.Flags,
		Command: strings.Join(m.Commands, " && "),
	}

	return r.Line()
}
