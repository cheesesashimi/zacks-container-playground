package examples

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type DNFPackage interface {
	Package() (string, error)
}

type DNFPackageName string

func (d DNFPackageName) Package() (string, error) {
	return string(d), nil
}

// @PackageEnvGroup
type DNFPackageEnvGroup string

func (d DNFPackageEnvGroup) Package() (string, error) {
	return fmt.Sprintf("@%s", d), nil
}

// name
// name.arch
// name-[epoch:]version
// name-[epoch:]version-release
// name-[epoch:]version-release.arch
type NEVRA struct {
	Name    string
	Epoch   string
	Version string
	Release string
	Arch    string
}

func (n NEVRA) CommandSuffix() string {
	evraPopulated := n.evraPopulated()
	nameEmpty := isEmptyString(n.Name)

	if evraPopulated && nameEmpty {
		return "-nevra"
	}

	if !evraPopulated && nameEmpty {
		return "-n"
	}

	if n.Arch != "" && nameEmpty {
		return "-na"
	}

	return ""
}

func (n NEVRA) evraPopulated() bool {
	return n.Epoch != "" && n.Version != "" && n.Release != "" && n.Arch != ""
}

func (n NEVRA) validate() error {
	values := map[string]string{
		"N": n.Name,
		"E": n.Epoch,
		"V": n.Version,
		"R": n.Release,
		"A": n.Arch,
	}

	valueNames := map[string]string{
		"N": "name",
		"E": "epoch",
		"V": "version",
		"R": "release",
		"A": "arch",
	}

	valueCombo := mapset.NewSet[string]()
	for letter, value := range values {
		if !isEmptyString(value) {
			valueCombo.Add(letter)
		}
	}

	validCombos := []mapset.Set[string]{
		mapset.NewSet[string]("N", "E", "V", "R", "A"),
		mapset.NewSet[string]("N", "E", "V", "R"),
		mapset.NewSet[string]("N", "E", "V"),
		mapset.NewSet[string]("N", "A"),
		mapset.NewSet[string]("N"),
	}

	for _, combo := range validCombos {
		if combo.Equal(valueCombo) {
			return nil
		}
	}

	out := []string{}
	for letter, name := range valueNames {
		if !valueCombo.Contains(letter) {
			out = append(out, name)
		}
	}

	return fmt.Errorf("one or more required values missing: %v", out)
}

func (n NEVRA) Package() (string, error) {
	if err := n.validate(); err != nil {
		return "", err
	}

	switch {
	case !isEmptyString(n.Epoch) && !isEmptyString(n.Version) && !isEmptyString(n.Release) && !isEmptyString(n.Arch):
		return fmt.Sprintf("%s-%s:%s-%s.%s", n.Name, n.Epoch, n.Version, n.Release, n.Arch), nil
	case !isEmptyString(n.Epoch) && !isEmptyString(n.Version) && !isEmptyString(n.Release):
		return fmt.Sprintf("%s-%s:%s-%s", n.Name, n.Epoch, n.Version, n.Release), nil
	case !isEmptyString(n.Epoch) && !isEmptyString(n.Version):
		return fmt.Sprintf("%s-%s:%s", n.Name, n.Epoch, n.Version), nil
	case !isEmptyString(n.Arch):
		return fmt.Sprintf("%s.%s", n.Name, n.Arch), nil
	default:
		return n.Name, nil
	}
}

func isEmptyString(in string) bool {
	if in == "" {
		return true
	}

	if in == " " {
		return true
	}

	if in == strings.Repeat(" ", len(in)) {
		return true
	}

	return false
}
