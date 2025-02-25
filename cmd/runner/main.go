package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheesesashimi/zacks-openshift-playground/internal/command"
)

func podmanrun() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	containerVolume := "/src"

	containerVolumes := []command.Volume{
		{
			HostPath:      cwd,
			ContainerPath: containerVolume,
			Opts:          "z",
		},
	}

	for i := 0; i <= 10; i++ {
		containerVolumes = append(containerVolumes, command.Volume{
			HostPath:      cwd,
			ContainerPath: fmt.Sprintf("%s%d", containerVolume, i),
			Opts:          "z",
		})
	}

	pr := &command.PodmanRun{
		Name:        "my-container",
		Interactive: true,
		Tty:         true,
		Remove:      true,
		Image:       "registry.fedoraproject.org/fedora:41",
		Volumes:     containerVolumes,
		ImageOpts: []command.Arg{
			&command.Subcommand{
				Name: "ls",
				Flags: []command.Flag{
					command.SingleSwitchFlag("l"),
					command.SingleSwitchFlag("a"),
				},
			},
			command.PositionalArg(containerVolume),
		},
	}

	fmt.Println(pr.Command())
}

func podmanbuild() {
	pb := &command.PodmanBuild{
		Tag: "quay.io/zzlotnik/something:latest",
		Labels: []command.Label{
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

	fmt.Println(pb.Command())
}

func podmanpush() {
	pp := &command.PodmanPush{
		Authfile: "/path/to/authfile",
		Image:    "quay.io/zzlotnik/something:latest",
	}

	fmt.Println(pp.Command())
}

func podmantag() {
	pt := &command.PodmanTag{
		Image:      "localhost/image:latest",
		TargetName: "quay.io/zzlotnik/image:latest",
	}

	fmt.Println(pt.Command())
}

func buildahbuild() {
	b := &command.BuildahBuild{
		Authfile: "/path/to/authfile",
		File:     "/path/to/containerfile",
		Proxy: &command.Proxy{
			Http:    "http://path.to.proxy",
			Https:   "https://path.to.https.proxy",
			NoProxy: "",
		},
		StorageDriver: "vfs",
		Tag:           "quay.io/zzlotnik/something:latest",
	}

	fmt.Println(b.Command())
}

func buildahpush() {
	b := &command.BuildahPush{
		Authfile:      "/path/to/authfile",
		CertDir:       "/path/to/cert-dir",
		Digestfile:    "/path/to/digestfile",
		StorageDriver: "vfs",
		Tag:           "quay.io/zzlotnik/something:latest",
	}

	fmt.Println(b.Command())
}

func emulateBuildahBuildAndPush() {
	destRoot := "/etc/pki/ca-trust/extracted"

	authfile := "/path/to/authfile"

	image := "quay.io/zzlotnik/something:latest"

	mountOpts := "z,rw"

	items := []interface{ Command() *command.Command }{
		&command.P11KitExtract{
			Format:    "openssl-bundle",
			Filter:    "certificates",
			Overwrite: true,
			Comment:   true,
			DestPath:  filepath.Join(destRoot, "/openssl/ca-bundle.trust.crt"),
		},
		&command.P11KitExtract{
			Format:    "pem-bundle",
			Filter:    "ca-anchors",
			Overwrite: true,
			Comment:   true,
			Purpose:   "server-auth",
			DestPath:  filepath.Join(destRoot, "/pem/tls-ca-bundle.pem"),
		},
		&command.BuildahBuild{
			Authfile: authfile,
			LogLevel: "DEBUG",
			File:     "/path/to/containerfile",
			Proxy: &command.Proxy{
				Http:    "http://proxy.host.com",
				Https:   "https://proxy.host.com",
				NoProxy: "",
			},
			StorageDriver: "vfs",
			Tag:           image,
			Volumes: []command.Volume{
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
		&command.BuildahPush{
			Authfile:   authfile,
			LogLevel:   "DEBUG",
			Image:      image,
			Digestfile: "/path/to/digestfile",
			CertDir:    "/var/run/secrets/kubernetes.io/serviceaccount",
		},
	}

	for _, item := range items {
		fmt.Println(item.Command())
	}
}

func main() {
	emulateBuildahBuildAndPush()

	l := command.Login{
		Token:  "aosifuhjdasiojf",
		Server: "server-url",
	}

	fmt.Println(l.Command("/path/to/kubeconfig"))

	re := command.ReleaseExtract{
		RegistryConfig:   "/path/to/registry/config",
		CommandToExtract: "openshift-install",
		ReleasePullspec:  "registry.hostname.com/org/repo:tag",
		To:               "/path/on/local/disk",
	}

	fmt.Println(re.Command("/path/to/kubeconfig"))
}
