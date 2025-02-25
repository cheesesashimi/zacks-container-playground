package containerfile

import (
	"fmt"
	"strings"

	"github.com/cheesesashimi/zacks-openshift-playground/internal/command"
)

type Stage struct {
	Name  string
	Image string
	Steps []ContainerfileStep
}

type Containerfile struct {
	Stages []*Stage
	Tag    string
}

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

type Mount struct {
	From            string
	Source          string
	Target          string
	BindPropagation string
	Type            string
	Opts            string
}

func (m *Mount) flag() string {
	out := []string{}

	ordered := []string{
		"type",
		"from",
		"source",
		"target",
		"bind-propagation",
	}

	mountOpts := map[string]string{
		"type":             m.Type,
		"from":             m.From,
		"source":           m.Source,
		"target":           m.Target,
		"bind-propagation": m.BindPropagation,
	}

	for _, key := range ordered {
		val := mountOpts[key]
		if val != "" {
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		}
	}

	splitOpts := strings.Split(m.Opts, ",")
	for _, item := range splitOpts {
		out = append(out, item)
	}

	return fmt.Sprintf("--mount=%s", strings.Join(out, ","))
}

type Command interface {
	Command() *command.Command
}

type MultiCommandRunStep struct {
	Flags    []string
	Mounts   []*Mount
	Commands []Command
}

func (m *MultiCommandRunStep) Line() string {
	cmds := []string{}
	for _, cmd := range m.Commands {
		cmds = append(cmds, cmd.Command().String())
	}

	s := &MultiRunStep{
		Flags:    m.Flags,
		Mounts:   m.Mounts,
		Commands: cmds,
	}

	return s.Line()
}

type CommandRunStep struct {
	Flags   []string
	Mounts  []*Mount
	Command Command
}

func (c *CommandRunStep) Line() string {
	r := &RunStep{
		Flags:   c.Flags,
		Mounts:  c.Mounts,
		Command: c.Command.Command().String(),
	}

	return r.Line()
}

type RunStep struct {
	Flags   []string
	Mounts  []*Mount
	Command string
}

func (r *RunStep) Line() string {
	if len(r.Flags) == 0 && len(r.Mounts) == 0 {
		return fmt.Sprintf("RUN %s", r.Command)
	}

	out := []string{"RUN"}

	for _, mount := range r.Mounts {
		out = append(out, mount.flag())
	}

	out = append(out, r.Flags...)
	out = append(out, r.Command)

	return strings.Join(out, " ")
}

type MultiRunStep struct {
	Flags    []string
	Mounts   []*Mount
	Commands []string
}

func (m *MultiRunStep) Line() string {
	r := RunStep{
		Mounts:  m.Mounts,
		Flags:   m.Flags,
		Command: strings.Join(m.Commands, " && "),
	}

	return r.Line()
}

type LabelStep struct {
	Key   string
	Value string
}

func (l *LabelStep) Line() string {
	return fmt.Sprintf("LABEL %s=%s", l.Key, l.Value)
}

type CopyStep struct {
	From  string
	Src   string
	Dest  string
	Flags []string
}

func (c *CopyStep) Line() string {
	if c.From == "" && len(c.Flags) == 0 {
		return fmt.Sprintf("COPY %s %s", c.Src, c.Dest)
	}

	out := []string{"COPY"}
	out = append(out, fmt.Sprintf("--from=%s", c.From))
	out = append(out, c.Flags...)
	out = append(out, []string{c.Src, c.Dest}...)
	return strings.Join(out, " ")
}
