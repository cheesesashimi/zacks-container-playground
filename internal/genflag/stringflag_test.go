package genflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFlag(t *testing.T) {
	testCases := []struct {
		name              string
		optionFuncs       []optionFunc
		expected          string
		expectedSeparated []string
		errExpected       bool
	}{
		{
			name:              "No optionFuncs",
			optionFuncs:       []optionFunc{},
			expected:          "--name value",
			expectedSeparated: []string{"--name", "value"},
		},
		{
			name:              "Single optionFunc",
			optionFuncs:       []optionFunc{Single},
			expected:          "-name value",
			expectedSeparated: []string{"-name", "value"},
		},
		{
			name:              "Quoted optionFunc",
			optionFuncs:       []optionFunc{Quoted},
			expected:          `--name "value"`,
			expectedSeparated: []string{"--name", `"value"`},
		},
		{
			name:              "Equal separator optionFunc",
			optionFuncs:       []optionFunc{EqualSeparator},
			expected:          "--name=value",
			expectedSeparated: []string{"--name=value"},
		},
		{
			name:              "All optionFuncs",
			optionFuncs:       []optionFunc{Single, Quoted, EqualSeparator},
			expected:          `-name="value"`,
			expectedSeparated: []string{`-name="value"`},
		},
		{
			name:        "Invalid optionFuncs",
			optionFuncs: []optionFunc{Explicit},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sf, err := NewStringFlag("name", "value", testCase.optionFuncs...)
			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			}

			assert.NoError(t, err)
			s, err := sf.String()
			assert.NoError(t, err)

			assert.Equal(t, testCase.expected, s)

			sep, err := sf.Segmented()
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedSeparated, sep)

			assert.Equal(t, "name", sf.Name())
			assert.Equal(t, "value", sf.Value())
		})
	}
}
