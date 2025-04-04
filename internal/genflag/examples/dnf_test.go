package examples

import (
	"testing"

	"github.com/cheesesashimi/zacks-container-playground/internal/genflag"
	"github.com/stretchr/testify/assert"
)

func TestNEVRA(t *testing.T) {
	testCases := []struct {
		name           string
		pkg            NEVRA
		expected       string
		expectedSuffix string
		errExpected    bool
	}{
		{
			name:           "Name only",
			pkg:            NEVRA{Name: "pkg"},
			expected:       "pkg",
			expectedSuffix: "-n",
		},
		{
			name:           "Name and arch",
			pkg:            NEVRA{Name: "pkg", Arch: "x86_64"},
			expected:       "pkg.x86_64",
			expectedSuffix: "-na",
		},
		{
			name:     "Name epoch and version",
			pkg:      NEVRA{Name: "pkg", Epoch: "0", Version: "0.10.1"},
			expected: "pkg-0:0.10.1",
		},
		{
			name:     "Name epoch version and release",
			pkg:      NEVRA{Name: "pkg", Epoch: "0", Version: "0.20.1", Release: "1.fc41"},
			expected: "pkg-0:0.20.1-1.fc41",
		},
		{
			name:           "Name epoch version release and arch",
			pkg:            NEVRA{Name: "pkg", Epoch: "0", Version: "0.30.1", Release: "1.fc41", Arch: "x86_64"},
			expected:       "pkg-0:0.30.1-1.fc41.x86_64",
			expectedSuffix: "-nevra",
		},
		{
			name:        "missing name",
			pkg:         NEVRA{Name: ""},
			errExpected: true,
		},
		{
			name:        "Epoch but no version",
			pkg:         NEVRA{Name: "pkg", Epoch: "0"},
			errExpected: true,
		},
		{
			name:        "Release but no version",
			pkg:         NEVRA{Name: "pkg", Release: "1.fc41"},
			errExpected: true,
		},
		{
			name:        "Version but no epoch",
			pkg:         NEVRA{Name: "pkg", Version: "0.30.1"},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := testCase.pkg.Package()
			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}

func TestDNF(t *testing.T) {
	testCases := []struct {
		name        string
		dnf         DNF
		expectedCLI []string
		errExpected bool
	}{
		{
			name: "DNF Install - No global opts",
			dnf: DNF{
				Install: &DNFInstall{
					Packages: newDnfPackages([]string{"neovim", "sqlite3", "python3"}),
				},
			},
			expectedCLI: []string{"dnf", "install", "neovim", "sqlite3", "python3"},
		},
		{
			name: "DNF Install - With global opts",
			dnf: DNF{
				GlobalOpts: &DNFGlobalOpts{
					Yes: true,
				},
				Install: &DNFInstall{
					Packages: newDnfPackages([]string{"neovim", "sqlite3", "python3"}),
				},
			},
			expectedCLI: []string{"dnf", "--yes", "install", "neovim", "sqlite3", "python3"},
		},
		{
			name: "DNF Install - With CVE IDs",
			dnf: DNF{
				GlobalOpts: &DNFGlobalOpts{
					Yes: true,
				},
				Install: &DNFInstall{
					CVEIDs: CSVString{"cve-1", "cve-2", "cve-3"},
				},
			},
			expectedCLI: []string{"dnf", "--yes", "install", "--cves=cve-1,cve-2,cve-3"},
		},
		{
			name: "DNF Install - With Setopt",
			dnf: DNF{
				GlobalOpts: &DNFGlobalOpts{
					SetOpt: []DNFSetOpt{
						{
							Key:     "keepcache",
							Enabled: true,
						},
						{
							Key:   "repodir",
							Value: "/local",
						},
					},
					Yes: true,
				},
				Install: &DNFInstall{
					CVEIDs: CSVString{"cve-1", "cve-2", "cve-3"},
				},
			},
			expectedCLI: []string{"dnf", "--setopt=keepcache=true", "--setopt=repodir=/local", "--yes", "install", "--cves=cve-1,cve-2,cve-3"},
		},
		{
			name: "DNF Install - With Advisory Severity",
			dnf: DNF{
				GlobalOpts: &DNFGlobalOpts{
					Yes: true,
				},
				Install: &DNFInstall{
					AdvisorySeverities: CSVString{DNFAdvisorySeverityCritical},
				},
			},
			expectedCLI: []string{"dnf", "--yes", "install", "--advisory-severities=critical"},
		},
		{
			name: "DNF Install - With Invalid Advisory Severity",
			dnf: DNF{
				GlobalOpts: &DNFGlobalOpts{
					Yes: true,
				},
				Install: &DNFInstall{
					AdvisorySeverities: CSVString{"invalid-severity"},
				},
			},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			results, err := testCase.dnf.Command()
			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCLI, results)
		})
	}
}

func marshalFlagsToString(v interface{}) ([]string, error) {
	f, err := genflag.Marshal(v)
	if err != nil {
		return nil, err
	}

	return flagsToStrings(f)
}

func flagsToStrings(flags []genflag.Flag) ([]string, error) {
	out := make([]string, len(flags))

	for i, flag := range flags {
		s, err := flag.String()
		if err != nil {
			return nil, err
		}

		out[i] = s
	}

	return out, nil
}
