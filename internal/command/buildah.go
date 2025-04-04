package command

// Runs the p11kitextract command.
type P11KitExtract struct {
	Format    string
	Filter    string
	Overwrite bool
	Comment   bool
	Purpose   string
	DestPath  string
}

func (p *P11KitExtract) Command() *Command {
	extractFlags := mapToValFlags(map[string]string{
		"format":  p.Format,
		"filter":  p.Filter,
		"purpose": p.Purpose,
	})

	extractFlags = append(extractFlags, mapToSwitchFlags(map[string]bool{
		"comment":   p.Comment,
		"overwrite": p.Overwrite,
	})...)

	return NewCommand("p11-kit", []Arg{
		&Subcommand{
			Name:  "extract",
			Flags: extractFlags,
		},
		PositionalArg(p.DestPath),
	})
}

// Represents proxy configuration.
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

// Represents a buildah build command.
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

func (b *BuildahBuild) Command() *Command {
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

	return NewCommand("buildah", []Arg{
		&Subcommand{
			Name:  "build",
			Flags: buildFlags,
		},
		PositionalArg(buildCtx),
	})
}

// Represents a buildah push command.
type BuildahPush struct {
	Authfile      string
	CertDir       string
	Digestfile    string
	Image         string
	LogLevel      string
	StorageDriver string
	Tag           string
}

func (b *BuildahPush) Command() *Command {
	buildFlags := mapToValFlags(map[string]string{
		"authfile":       b.Authfile,
		"cert-dir":       b.CertDir,
		"digestfile":     b.Digestfile,
		"log-level":      b.LogLevel,
		"storage-driver": b.StorageDriver,
	})

	return NewCommand("buildah", []Arg{
		&Subcommand{
			Name:  "push",
			Flags: buildFlags,
		},
		PositionalArg(b.Image),
	})
}
