package examples

import (
	"fmt"
	"strings"

	"github.com/cheesesashimi/zacks-container-playground/internal/genflag"
)

var _ genflag.Marshaler = PodmanEnvOpt{}
var _ genflag.Marshaler = &PodmanEnvOpt{}

// Holds the name and value arguments for a type that holds an --env name=value
// flag for Podman.
type PodmanEnvOpt struct {
	Name  string
	Value string
}

// Instantiates a flag from the contents of this struct by passing it to a
// concrete flag type.
func (p PodmanEnvOpt) MarshalFlags() ([]genflag.Flag, error) {
	val := fmt.Sprintf("%s=%s", p.Name, p.Value)
	f, err := genflag.NewStringFlag("env", val, genflag.Quoted)
	return []genflag.Flag{f}, err
}

// Holds the values needed to run a container image using Podman.
type PodmanRun struct {
	Interactive bool `genflag:""`
	TTY         bool `genflag:""`
	Remove      bool `genflag:"rm"`
	Env         []PodmanEnvOpt
	Entrypoint  string `genflag:""`
	Image       string
}

// Constructs a command-line incantation for running Podman according to the
// values in the struct.
func (p PodmanRun) PodmanCmd() (string, error) {
	flags, err := genflag.Marshal(p)
	if err != nil {
		return "", err
	}

	out := []string{"podman", "run"}

	for _, flag := range flags {
		segmented, err := flag.Segmented()
		if err != nil {
			return "", err
		}

		out = append(out, segmented...)
	}

	out = append(out, p.Image)

	return strings.Join(out, " "), nil
}

// Holds a build arg and value showcasing that genflag knows what to do with a
// stringer.
type BuildArg struct {
	Argument string
	Value    string
}

// Implements the stringer interface which concatenates the arg and the value.
func (b BuildArg) String() string {
	return fmt.Sprintf("%s=%s", b.Argument, b.Value)
}

// Holds the arguments needed to run podman build.
type PodmanBuild struct {
	Tag            string         `genflag:""`
	BuildArg       []BuildArg     `genflag:"build-arg,quoted"`
	File           string         `genflag:""`
	Env            []PodmanEnvOpt `genflag:""`
	DefaultContext string
}

// Constructs a command-line incantation for running Podman according to the
// values in the struct.
func (p PodmanBuild) PodmanCmd() (string, error) {
	flags, err := genflag.Marshal(p)
	if err != nil {
		return "", err
	}

	out := []string{"podman", "build"}

	for _, flag := range flags {
		segmented, err := flag.Segmented()
		if err != nil {
			return "", err
		}

		out = append(out, segmented...)
	}

	defCtx := p.DefaultContext
	if defCtx == "" {
		defCtx = "."
	}
	out = append(out, defCtx)

	return strings.Join(out, " "), nil
}
