package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheesesashimi/zacks-openshift-playground/internal/command"
	"github.com/cheesesashimi/zacks-openshift-playground/internal/containerfile"
	"github.com/davecgh/go-spew/spew"
	"github.com/openshift/imagebuilder"
	"github.com/openshift/imagebuilder/dockerfile/parser"
)

func getExtensionsRepoCommands(extensionPkgs []string) []containerfile.Command {
	extRepoFilePath := "/etc/yum.repos.d/coreos-extensions.repo"
	repoFile := `[coreos-extensions]
enabled=1
metadata_expire=1m
baseurl=/tmp/mco-extensions/os-extensions-content/usr/share/rpm-ostree/extensions/
gpgcheck=0
skip_if_unavailable=False`

	return []containerfile.Command{
		&command.Echo{
			Content:    repoFile,
			Escape:     true,
			RedirectTo: extRepoFilePath,
		},
		&command.Chmod{
			Path:      extRepoFilePath,
			Mode:      "644",
			Recursive: true,
		},
		&command.RpmOstreeInstall{
			Packages: extensionPkgs,
		},
		&command.Delete{
			Path: extRepoFilePath,
		},
	}
}

func generatingApproaches() error {
	mounts := []*containerfile.Mount{
		{
			Type:            "bind",
			From:            "registry.ci.openshift.org/ocp/4.18-2024-11-26-081001@sha256:cf24f7665d164b2e3072ddb22a3410ba802bf19d5c67638a84bc5faf4452f1fa",
			Source:          "/",
			Target:          "/tmp/mco-extensions/os-extensions-content",
			BindPropagation: "rshared",
			Opts:            "rw,z",
		},
		{
			Type:   "cache",
			Target: "/var/cache/dnf",
			Opts:   "z",
		},
		{
			Type:   "cache",
			Target: "/go/rhel9/.cache",
			Opts:   "z",
		},
		{
			Type:   "cache",
			Target: "/go/rhel9/pkg/mod",
			Opts:   "z",
		},
	}

	steps := containerfile.ContainerfileSteps{
		&containerfile.FromStep{
			Image: "registry.ci.openshift.org/ocp/4.18-2024-11-26-081001@sha256:7e2831f10dec594e0ee569918078b34c7c09df6f34a49b55d2271216c0ba6edc",
		},
		&containerfile.MultiCommandRunStep{
			Mounts:   mounts,
			Commands: getExtensionsRepoCommands([]string{"usbguard", "kata-containers", "kernel-rt"}),
		},
		&containerfile.RunStep{
			Command: "ostree container commit",
		},
	}

	fmt.Println("Building:")
	fmt.Println(steps.Containerfile())

	tempdir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempdir)

	containerfilePath := filepath.Join(tempdir, "Containerfile")

	if err := os.WriteFile(containerfilePath, []byte(steps.Containerfile()), 0o755); err != nil {
		return err
	}

	build := &command.PodmanBuild{
		Tag:  "localhost/something:latest",
		File: containerfilePath,
	}

	cmd := build.Command().Cmd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

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
	//	parsingApproaches()
	if err := generatingApproaches(); err != nil {
		panic(err)
	}
}
