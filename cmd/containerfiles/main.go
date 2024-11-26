package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/openshift/imagebuilder"
	"github.com/openshift/imagebuilder/dockerfile/parser"
)

func generatingApproaches() error {
	steps := ContainerfileSteps{
		&FromStep{
			Image: "registry.fedoraproject.org/fedora:41",
			As:    "fedora",
		},
		&RunStep{
			Command: "echo 'hello world'",
		},
		&MultiRunStep{
			Commands: []string{
				"cp -r -v /etc /etc/something",
				"go build .",
				"go test -count=1 -v ./...",
			},
		},
	}

	spew.Dump(steps.Node())

	fmt.Println(steps.Node().Dump())
	fmt.Println(steps.Line())
	return nil
}

func parsingApproaches() {
	path := "./Containerfile.generated-with-extensions"
	// path := "/home/zzlotnik/Repos/machine-config-operator/Dockerfile"
	if err := approach2(path); err != nil {
		panic(err)
	}

	if err := approach3(path); err != nil {
		panic(err)
	}
}

func parseFile(path string) (*parser.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	parsed, err := imagebuilder.ParseDockerfile(file)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func traverseNodes(node *parser.Node, visit func(*parser.Node)) {
	visit(node)

	for _, child := range node.Children {
		traverseNodes(child, visit)
	}

	for n := node.Next; n != nil; n = n.Next {
		traverseNodes(n, visit)
	}
}

func approach1(path string) error {
	parsed, err := parseFile(path)
	if err != nil {
		return err
	}

	fromNodes := []*parser.Node{}

	baseImages := []string{}
	aliases := []string{}

	traverseNodes(parsed, func(n *parser.Node) {
		if n.Value == "from" {
			fromNodes = append(fromNodes, n)
		}
	})

	// This actually works.
	for _, fromNode := range fromNodes {
		// Images
		if fromNode.Next != nil {
			baseImages = append(baseImages, fromNode.Next.Value)
		}

		// Aliases
		if fromNode.Next.Next != nil && fromNode.Next.Next.Value == "AS" && fromNode.Next.Next.Next != nil {
			aliases = append(aliases, fromNode.Next.Next.Next.Value)
		}
	}

	spew.Dump(baseImages)
	spew.Dump(aliases)

	return nil
}

func approach2(path string) error {
	parsed, err := parseFile(path)
	if err != nil {
		return err
	}

	builder := imagebuilder.NewBuilder(map[string]string{})

	stages, err := imagebuilder.NewStages(parsed, builder)
	if err != nil {
		return err
	}

	for _, stage := range stages {
		fmt.Println(stage.Builder.From(stage.Node))
		spew.Dump(stage)
	}

	return nil
}

func approach3(path string) error {
	return nil
}

func main() {
	parsingApproaches()
	if err := generatingApproaches(); err != nil {
		panic(err)
	}
}
