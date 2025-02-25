package containerfile

import (
	"fmt"
	"strings"

	"github.com/cheesesashimi/zacks-container-playground/internal/command"
)

// Top-level Containerfile object
type Containerfile struct {
	Stages []*Stage
	Tag    string
}

func (c *Containerfile) String() string {
	sb := &strings.Builder{}

	for _, stage := range c.Stages {
		fmt.Fprintln(sb, stage.Line())
	}

	return sb.String()
}

// Represents a single stage in a Containerfile, including its base image. Each
// Containerfile must have at least one stage.
type Stage struct {
	Name  string
	Image string
	Steps []ContainerfileStep
}

func (s *Stage) Line() string {
	sb := &strings.Builder{}

	from := &FromStep{
		Image: s.Image,
		As:    s.Name,
	}

	fmt.Fprintln(sb, from.Line())

	for _, step := range s.Steps {
		fmt.Fprintln(sb, step.Line())
	}

	return sb.String()
}

// Catch-all for a given step in a Containerfile. Each *Step type must
// implement this interface. For now, this just emits a string representation
// of what the given step should do.
type ContainerfileStep interface {
	Line() string
}

// Represents a FROM statement.
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

// A Command primitive knows how to construct a CLI command given its options.
type Command interface {
	Command() *command.Command
}

// Chains multiple Commands together with && in between them.
type MultiCommandRunStep struct {
	// A flag is an option that gets passed to the RUN directive.
	Flags []string
	// A Mount is a specific --mount option that gets passed to the RUN directive.
	Mounts []*Mount
	// A list of commands to execute.
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

// Represents a RUN statement.
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

// Chains multiple command literals together with && and runs them.
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

// Represents a LABEL statement.
type LabelStep struct {
	Key   string
	Value string
}

func (l *LabelStep) Line() string {
	return fmt.Sprintf("LABEL %s=%s", l.Key, l.Value)
}

// Represents a COPY statement.
type CopyStep struct {
	From  string
	Src   string
	Dest  string
	Flags []string
}

// Shortcut for copying everything from the given build context into the
// container.
func CopyAllStep() ContainerfileStep {
	return &CopyStep{
		Src:  ".",
		Dest: ".",
	}
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

// Represents a WORKDIR statement.
type WorkDirStep string

func NewWorkDirStep(path string) ContainerfileStep {
	w := WorkDirStep(path)
	return &w
}

func (w *WorkDirStep) Line() string {
	return fmt.Sprintf("WORKDIR %s", *w)
}

// Represents a USER statement.
type UserStep string

func NewUserStep(user string) ContainerfileStep {
	u := UserStep(user)
	return &u
}

func (u *UserStep) Line() string {
	return fmt.Sprintf("USER %s", *u)
}
