package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheesesashimi/zacks-openshift-playground/internal/command"
	"github.com/cheesesashimi/zacks-openshift-playground/internal/containerfile"
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

func generateContainerfile() containerfile.ContainerfileSteps {
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

	return containerfile.ContainerfileSteps{
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
}

func buildContainerfile(steps containerfile.ContainerfileSteps) error {
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

func main() {
	cfile := generateContainerfile()

	fmt.Println("Containerfile:")
	fmt.Println(cfile.Containerfile())
}
