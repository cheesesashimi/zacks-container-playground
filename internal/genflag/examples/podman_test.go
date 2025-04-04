package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodmanStructs(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{ PodmanCmd() (string, error) }
		expected    []string
		notExpected []string
	}{
		{
			name:     "No options set",
			input:    PodmanRun{},
			expected: []string{"podman", "run"},
		},
		{
			name: "With bool options set",
			input: &PodmanRun{
				Interactive: true,
				TTY:         true,
				Remove:      true,
			},
			expected: []string{"--interactive", "--tty", "--rm"},
		},
		{
			name: "With string options set",
			input: &PodmanRun{
				Interactive: true,
				TTY:         true,
				Remove:      true,
				Image:       "alpine:latest",
			},
			expected: []string{"podman", "run", "--interactive", "--tty", "--rm", "alpine:latest"},
		},
		{
			name: "With env vars set with custom flagmarshaler",
			input: &PodmanRun{
				Interactive: true,
				TTY:         true,
				Remove:      true,
				Image:       "alpine:latest",
				Env: []PodmanEnvOpt{
					{
						Name:  "HOME",
						Value: "/home/zack",
					},
				},
			},
			expected: []string{"podman", "run", "--interactive", "--tty", "--rm", "alpine:latest", `--env "HOME=/home/zack"`},
		},
		{
			name:     "Build with no args",
			input:    &PodmanBuild{},
			expected: []string{"podman", "build", "."},
		},
		{
			name: "Build with args",
			input: &PodmanBuild{
				Tag:  "final:latest",
				File: "Containerfile.dev",
			},
			expected: []string{
				"podman", "build", "--tag final:latest", "--file Containerfile.dev",
			},
		},
		{
			name: "Build with buildargs as stringers",
			input: &PodmanBuild{
				Tag:  "final:latest",
				File: "Containerfile.dev",
				BuildArg: []BuildArg{
					{
						Argument: "arg",
						Value:    "val",
					},
				},
			},
			expected: []string{
				"podman", "build", "--tag final:latest", `--build-arg "arg=val"`, "--file Containerfile.dev",
			},
		},
		{
			name: "Build with buildargs and env vars",
			input: &PodmanBuild{
				Tag:  "final:latest",
				File: "Containerfile.dev",
				Env: []PodmanEnvOpt{
					{
						Name:  "HOME",
						Value: "/home/zack",
					},
				},
				BuildArg: []BuildArg{
					{
						Argument: "arg",
						Value:    "val",
					},
				},
			},
			expected: []string{
				"podman", "build", "--tag final:latest", `--env "HOME=/home/zack"`, `--build-arg "arg=val"`, "--file Containerfile.dev",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := testCase.input.PodmanCmd()
			assert.NoError(t, err)

			if len(testCase.expected) == 0 && len(output) != 0 {
				t.Fatalf("expected is empty, but got: %q", output)
			}

			for _, item := range testCase.expected {
				assert.Contains(t, output, item)
			}

			for _, item := range testCase.notExpected {
				assert.NotContains(t, output, item)
			}
		})
	}
}
