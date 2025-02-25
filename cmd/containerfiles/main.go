package main

import (
	"fmt"
	"path/filepath"

	"github.com/cheesesashimi/zacks-container-playground/internal/command"
	"github.com/cheesesashimi/zacks-container-playground/internal/containerfile"
)

type packageManager string

const (
	dnf    packageManager = "dnf"
	yum    packageManager = "yum"
	aptGet packageManager = "apt-get"
)

// Holds options applicable for all containerfiles that we generate.
type otherContainerfileOpts struct {
	// The base image to pull
	baseImage string
	// The username to create
	username string
	// The packages to install
	packages []string
	// The package manager to call for installation
	packageManager packageManager
}

// Generates the containerfile.
func (o *otherContainerfileOpts) containerfile() containerfile.Containerfile {
	return containerfile.Containerfile{
		Stages: []*containerfile.Stage{
			// This containerfile only has a single stage.
			{
				Image: o.baseImage,
				Steps: []containerfile.ContainerfileStep{
					&containerfile.MultiCommandRunStep{
						Commands: append(o.getPackageCommands(), &command.CommandLiteral{"useradd", o.username}),
					},
					containerfile.NewUserStep(o.username),
				},
			},
		},
	}
}

// Gets the appropriate package command(s) given the supplied package manager.
func (o *otherContainerfileOpts) getPackageCommands() []containerfile.Command {
	if o.packageManager == dnf {
		return []containerfile.Command{
			&command.DnfInstall{
				Yes:      true,
				Packages: o.packages,
			},
		}
	}

	if o.packageManager == aptGet {
		return []containerfile.Command{
			// Needs an update step before we can install packages.
			&command.AptGetUpdate{},
			&command.AptGetInstall{
				Yes:      true,
				Packages: o.packages,
			},
		}
	}

	if o.packageManager == yum {
		return []containerfile.Command{
			&command.YumInstall{
				Yes:      true,
				Packages: o.packages,
			},
		}
	}

	// Default to dnf if no other option exists.
	return []containerfile.Command{
		&command.DnfInstall{
			Yes:      true,
			Packages: o.packages,
		},
	}
}

// Returns the containerfiles that were templatized in the slide deck.
func generateOtherContainerfiles() []containerfile.Containerfile {
	items := map[string]packageManager{
		"registry.fedoraproject.org/fedora:latest": dnf,
		"quay.io/centos/stream:9":                  yum,
		"ubuntu:latest":                            aptGet,
	}

	out := []containerfile.Containerfile{}

	for baseImage, packageManager := range items {
		opts := otherContainerfileOpts{
			baseImage:      baseImage,
			packageManager: packageManager,
			packages:       []string{"nvim", "git", "golang"},
			username:       "zack",
		}

		out = append(out, opts.containerfile())
	}

	return out
}

func generateGolangBuildContainerfile() containerfile.Containerfile {
	baseImage := "registry.fedoraproject.org/fedora:latest"
	workdirPath := "/go/src/github.com/cheesesashimi/zacks-openshift-helpers"

	return containerfile.Containerfile{
		Stages: []*containerfile.Stage{
			{
				Name:  "builder",
				Image: baseImage,
				Steps: []containerfile.ContainerfileStep{
					containerfile.NewWorkDirStep(workdirPath),
					&containerfile.CommandRunStep{
						Command: &command.DnfInstall{
							Yes:      true,
							Packages: []string{"golang"},
						},
					},
					containerfile.CopyAllStep(),
					&containerfile.RunStep{
						Command: "make all",
					},
				},
			},
			{
				Name:  "final",
				Image: baseImage,
				Steps: []containerfile.ContainerfileStep{
					&containerfile.CopyStep{
						From: "builder",
						Src:  filepath.Join(workdirPath, "_output"),
						Dest: "/usr/local/bin/",
					},
				},
			},
		},
	}
}

func main() {
	fmt.Println("Golang Containerfile:")
	cfile := generateGolangBuildContainerfile()
	fmt.Println(cfile.String())

	for _, cfile := range generateOtherContainerfiles() {
		fmt.Println(cfile.String())
	}
}
