package main

type P11KitExtract struct {
	Format    string
	Filter    string
	Overwrite bool
	Comment   bool
	Purpose   string
	DestPath  string
}

func (p *P11KitExtract) Args() []Arg {
	extractFlags := mapToValFlags(map[string]string{
		"format":  p.Format,
		"filter":  p.Filter,
		"purpose": p.Purpose,
	})

	extractFlags = append(extractFlags, mapToSwitchFlags(map[string]bool{
		"comment":   p.Comment,
		"overwrite": p.Overwrite,
	})...)

	return []Arg{
		CommandName("p11-kit"),
		&Subcommand{
			Name:  "extract",
			Flags: extractFlags,
		},
		PositionalArg(p.DestPath),
	}
}

type Proxy struct {
	Http    string
	Https   string
	NoProxy string
}

func (p *Proxy) flags() []Flag {
	args := map[string]string{
		"HTTP_PROXY":  p.Http,
		"HTTPS_PROXY": p.Https,
		"NO_PROXY":    p.NoProxy,
	}

	flags := []Flag{}
	for arg, val := range args {
		buildArg := BuildArg{
			Name:  arg,
			Value: val,
		}

		flags = append(flags, buildArg.flag())
	}

	return flags
}

type BuildahBuild struct {
	Authfile      string
	BuildArgs     []BuildArg
	BuildContext  string
	File          string
	LogLevel      string
	Proxy         *Proxy
	StorageDriver string
	Tag           string
	Volumes       []Volume
}

func (b *BuildahBuild) Args() []Arg {
	buildFlags := mapToValFlags(map[string]string{
		"authfile":       b.Authfile,
		"file":           b.File,
		"log-level":      b.LogLevel,
		"storage-driver": b.StorageDriver,
		"tag":            b.Tag,
	})

	for _, buildArg := range b.BuildArgs {
		buildFlags = append(buildFlags, buildArg.flag())
	}

	for _, volume := range b.Volumes {
		buildFlags = append(buildFlags, volume.flag())
	}

	if b.Proxy != nil {
		buildFlags = append(buildFlags, b.Proxy.flags()...)
	}

	buildCtx := b.BuildContext
	if buildCtx == "" {
		buildCtx = "."
	}

	return []Arg{
		CommandName("buildah"),
		&Subcommand{
			Name:  "build",
			Flags: buildFlags,
		},
		PositionalArg(buildCtx),
	}
}

type BuildahPush struct {
	Authfile      string
	CertDir       string
	Digestfile    string
	Image         string
	LogLevel      string
	StorageDriver string
	Tag           string
}

func (b *BuildahPush) Args() []Arg {
	buildFlags := mapToValFlags(map[string]string{
		"authfile":       b.Authfile,
		"cert-dir":       b.CertDir,
		"digestfile":     b.Digestfile,
		"log-level":      b.LogLevel,
		"storage-driver": b.StorageDriver,
	})

	return []Arg{
		CommandName("buildah"),
		&Subcommand{
			Name:  "push",
			Flags: buildFlags,
		},
		PositionalArg(b.Image),
	}
}
