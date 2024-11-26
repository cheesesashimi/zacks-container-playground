package main

import (
	"fmt"
	"strings"

	"github.com/openshift/imagebuilder/dockerfile/parser"
)

type ContainerfileSteps []ContainerfileStep

func (c ContainerfileSteps) Node() *parser.Node {
	root := &parser.Node{}

	for _, step := range c {
		root.Children = append(root.Children, step.Node())
	}

	return root
}

func (c ContainerfileSteps) Line() string {
	out := &strings.Builder{}

	for _, step := range c {
		fmt.Fprintln(out, step.Line())
	}

	return out.String()
}

type ContainerfileStep interface {
	Node() *parser.Node
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

func (f *FromStep) Node() *parser.Node {
	as := &parser.Node{
		Value: "as",
		Next: &parser.Node{
			Value: f.As,
		},
	}

	from := &parser.Node{
		Value: "from",
		Next: &parser.Node{
			Value: f.Image,
		},
	}

	if f.As != "" {
		from.Next.Next = as
	}

	return nestIntoParent(from)
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

func (r *RunStep) Node() *parser.Node {
	return nestIntoParent(&parser.Node{
		Value: "run",
		Next: &parser.Node{
			Value: r.Command,
		},
		Flags: r.Flags,
	})
}

func nestIntoParent(n *parser.Node) *parser.Node {
	return &parser.Node{
		Next: n,
	}
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

func (m *MultiRunStep) Node() *parser.Node {
	r := RunStep{
		Flags:   m.Flags,
		Command: strings.Join(m.Commands, " && "),
	}

	return r.Node()
}
