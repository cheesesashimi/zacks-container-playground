package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	args []Arg
}

func NewCommand(args []Arg) (*Command, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no args found")
	}

	if _, ok := args[0].(CommandName); !ok {
		return nil, fmt.Errorf("first arg must be CommandName, got: %T", args[0])
	}

	return &Command{
		args: args,
	}, nil
}

func (c *Command) String() string {
	return strings.Join(renderArgs(c.args), " ")
}

func (c *Command) Cmd() *exec.Cmd {
	return exec.Command(c.args[0].Arg()[0], renderArgs(c.args[1:])...)
}

type CommandName string

func (c CommandName) Arg() []string {
	return []string{string(c)}
}

type PositionalArg string

func (a PositionalArg) Arg() []string {
	return []string{string(a)}
}

type Arg interface {
	Arg() []string
}

type Subcommand struct {
	Name  string
	Flags []Flag
}

func (s *Subcommand) Arg() []string {
	out := []string{s.Name}
	return append(out, renderFlags(s.Flags)...)
}

type Flag interface {
	Arg() []string
}

type SingleSwitchFlag string

func (s SingleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", string(s))}
}

type DoubleSwitchFlag string

func (d DoubleSwitchFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", string(d))}
}

type SingleValueFlag struct {
	Name  string
	Value string
}

func (s *SingleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s", s.Name), s.Value}
}

type SingleEqualValueFlag struct {
	Name  string
	Value string
}

func (s *SingleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("-%s=%s", s.Name, s.Value)}
}

type DoubleEqualValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleEqualValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s=%s", d.Name, d.Value)}
}

type DoubleValueFlag struct {
	Name  string
	Value string
}

func (d *DoubleValueFlag) Arg() []string {
	return []string{fmt.Sprintf("--%s", d.Name), d.Value}
}

func CombineArgs(arg []Arg) []string {
	return renderArgs(arg)
}

func renderArgs(args []Arg) []string {
	out := []string{}

	for _, arg := range args {
		out = append(out, arg.Arg()...)
	}

	return out
}

func renderFlags(flags []Flag) []string {
	out := []string{}

	for _, flag := range flags {
		out = append(out, flag.Arg()...)
	}

	return out
}

func podmanrun() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	containerVolume := "/src"

	containerVolumes := []Volume{
		{
			HostPath:      cwd,
			ContainerPath: containerVolume,
			Opts:          "z",
		},
	}

	for i := 0; i <= 10; i++ {
		containerVolumes = append(containerVolumes, Volume{
			HostPath:      cwd,
			ContainerPath: fmt.Sprintf("%s%d", containerVolume, i),
			Opts:          "z",
		})
	}

	pr := &PodmanRun{
		Name:        "my-container",
		Interactive: true,
		Tty:         true,
		Remove:      true,
		Image:       "registry.fedoraproject.org/fedora:41",
		Volumes:     containerVolumes,
		ImageOpts: []Arg{
			&Subcommand{
				Name: "ls",
				Flags: []Flag{
					SingleSwitchFlag("l"),
					SingleSwitchFlag("a"),
				},
			},
			PositionalArg(containerVolume),
		},
	}

	dumpArgs(pr.Args())
}

func podmanbuild() {
	pb := &PodmanBuild{
		Tag: "quay.io/zzlotnik/something:latest",
		Labels: []Label{
			{
				Name:  "label1",
				Value: "value1",
			},
			{
				Name:  "label2",
				Value: "value2",
			},
		},
	}

	dumpArgs(pb.Args())
}

func podmanpush() {
	pp := &PodmanPush{
		Authfile: "/path/to/authfile",
		Image:    "quay.io/zzlotnik/something:latest",
	}

	dumpArgs(pp.Args())

}

func podmantag() {
	pt := &PodmanTag{
		Image:      "localhost/image:latest",
		TargetName: "quay.io/zzlotnik/image:latest",
	}

	dumpArgs(pt.Args())
}

func dumpArgs(args []Arg) {
	cmd, err := NewCommand(args)
	if err != nil {
		panic(err)
	}

	fmt.Println(cmd)
}

func buildahbuild() {
	b := &BuildahBuild{
		Authfile: "/path/to/authfile",
		File:     "/path/to/containerfile",
		Proxy: &Proxy{
			Http:    "http://path.to.proxy",
			Https:   "https://path.to.https.proxy",
			NoProxy: "",
		},
		StorageDriver: "vfs",
		Tag:           "quay.io/zzlotnik/something:latest",
	}

	dumpArgs(b.Args())
}

func buildahpush() {
	b := &BuildahPush{
		Authfile:      "/path/to/authfile",
		CertDir:       "/path/to/cert-dir",
		Digestfile:    "/path/to/digestfile",
		StorageDriver: "vfs",
		Tag:           "quay.io/zzlotnik/something:latest",
	}

	dumpArgs(b.Args())
}

func emulateBuildahBuildAndPush() {
	destRoot := "/etc/pki/ca-trust/extracted"

	authfile := "/path/to/authfile"

	image := "quay.io/zzlotnik/something:latest"

	mountOpts := "z,rw"

	items := []interface{ Args() []Arg }{
		&P11KitExtract{
			Format:    "openssl-bundle",
			Filter:    "certificates",
			Overwrite: true,
			Comment:   true,
			DestPath:  filepath.Join(destRoot, "/openssl/ca-bundle.trust.crt"),
		},
		&P11KitExtract{
			Format:    "pem-bundle",
			Filter:    "ca-anchors",
			Overwrite: true,
			Comment:   true,
			Purpose:   "server-auth",
			DestPath:  filepath.Join(destRoot, "/pem/tls-ca-bundle.pem"),
		},
		&BuildahBuild{
			Authfile: authfile,
			LogLevel: "DEBUG",
			File:     "/path/to/containerfile",
			Proxy: &Proxy{
				Http:    "http://proxy.host.com",
				Https:   "https://proxy.host.com",
				NoProxy: "",
			},
			StorageDriver: "vfs",
			Tag:           image,
			Volumes: []Volume{
				{
					HostPath:      "$ETC_PKI_RPM_GPG_MOUNTPOINT",
					ContainerPath: "$ETC_PKI_RPM_GPG_MOUNTPOINT",
					Opts:          mountOpts,
				},
				{
					HostPath:      "$ETC_YUM_REPOS_D_MOUNTPOINT",
					ContainerPath: "$ETC_YUM_REPOS_D_MOUNTPOINT",
					Opts:          mountOpts,
				},
			},
		},
		&BuildahPush{
			Authfile:   authfile,
			LogLevel:   "DEBUG",
			Image:      image,
			Digestfile: "/path/to/digestfile",
			CertDir:    "/var/run/secrets/kubernetes.io/serviceaccount",
		},
	}

	for _, item := range items {
		dumpArgs(item.Args())
	}
}

func main() {
	emulateBuildahBuildAndPush()
}
