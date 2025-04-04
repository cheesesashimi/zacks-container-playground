package examples

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func newDnfPackages(names []string) []DNFPackage {
	out := []DNFPackage{}

	for _, name := range names {
		name := name
		out = append(out, DNFPackageName(name))
	}

	return out
}

type CSVString []string

func (c CSVString) String() string {
	return strings.Join([]string(c), ",")
}

type DNF struct {
	GlobalOpts *DNFGlobalOpts
	Install    *DNFInstall
	// To Implement:
	// Upgrade   *DNFUpgrade
	// Remove    *DNFRemove
	// RepoQuery *DNFRepoQuery
	// Search    *DNFSearch
	// List      *DNFList
	// Info      *DNFInfo
}

func (d DNF) Command() ([]string, error) {
	out := []string{"dnf"}

	if d.GlobalOpts != nil {
		globalOpts, err := d.GlobalOpts.Command()
		if err != nil {
			return nil, err
		}

		out = append(out, globalOpts...)
	}

	if d.Install != nil {
		installOpts, err := d.Install.Command()
		if err != nil {
			return nil, err
		}

		out = append(out, installOpts...)
	}

	return out, nil
}

const (
	DNFAdvisorySeverityCritical  string = "critical"
	DNFAdvisorySeverityImportant string = "important"
	DNFAdvisorySeverityModerate  string = "moderate"
	DNFAdvisorySeverityLow       string = "low"
	DNFAdvisorySeverityNone      string = "none"
)

type DNFInstall struct {
	Advisories         CSVString `genflag:"equaled"`
	AdvisorySeverities CSVString `genflag:"advisory-severities,equaled"`
	AllowDowngrade     bool      `genflag:""`
	AllowErasing       bool      `genflag:""`
	Bugfix             bool      `genflag:""`
	BugzillaIDs        CSVString `genflag:"bzs,equaled"`
	CVEIDs             CSVString `genflag:"cves,equaled"`
	DownloadOnly       bool      `genflag:""`
	Enhancement        bool      `genflag:""`
	NewPackage         bool      `genflag:""`
	NoAllowDowngrade   bool      `genflag:"no-allow-downgrade"`
	Offline            bool      `genflag:""`
	Security           bool      `genflag:""`
	SkipBroken         bool      `genflag:"skip-broken"`
	SkipUnavailable    bool      `genflag:"skip-unavailable"`
	Store              string    `genflag:"equaled"`
	Packages           []DNFPackage
}

func (d DNFInstall) Command() ([]string, error) {
	out := []string{"install"}

	if err := d.validateAdvisorySeverity(); err != nil {
		return nil, err
	}

	s, err := marshalFlagsToString(d)
	if err != nil {
		return nil, err
	}

	out = append(out, s...)
	for _, pkg := range d.Packages {
		pkgStr, err := pkg.Package()
		if err != nil {
			return nil, err
		}
		out = append(out, pkgStr)
	}

	return out, nil
}

func (d DNFInstall) validateAdvisorySeverity() error {
	if len(d.AdvisorySeverities) == 0 {
		return nil
	}

	severities := mapset.NewSet[string](
		DNFAdvisorySeverityCritical,
		DNFAdvisorySeverityImportant,
		DNFAdvisorySeverityModerate,
		DNFAdvisorySeverityLow,
		DNFAdvisorySeverityNone,
	)

	in := mapset.NewSet[string](d.AdvisorySeverities...)

	if in.IsProperSubset(severities) {
		return nil
	}

	return fmt.Errorf("invalid severities: %v", severities.Difference(in).ToSlice())
}

type DNFSetOpt struct {
	Key     string
	Value   string
	Enabled bool
}

func (d DNFSetOpt) String() string {
	if d.Enabled {
		return fmt.Sprintf("%s=%v", d.Key, d.Enabled)
	}

	return fmt.Sprintf("%s=%s", d.Key, d.Value)
}

type DNFGlobalOpts struct {
	Config      string      `genflag:"equaled"`
	EnableRepo  []string    `genflag:"equaled"`
	DisableRepo []string    `genflag:"equaled"`
	No          bool        `genflag:""`
	NoDocs      bool        `genflag:""`
	Quiet       bool        `genflag:""`
	Refresh     bool        `genflag:""`
	SetOpt      []DNFSetOpt `genflag:"equaled"`
	Yes         bool        `genflag:""`
}

func (d DNFGlobalOpts) Command() ([]string, error) {
	return marshalFlagsToString(d)
}
