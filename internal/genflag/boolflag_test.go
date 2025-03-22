package genflag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolFlag(t *testing.T) {
	testCases := []struct {
		name              string
		value             bool
		optionFuncs       []optionFunc
		expected          string
		expectedSegmented []string
		errExpected       bool
	}{
		{
			name:              "Implicit true",
			value:             true,
			optionFuncs:       []optionFunc{},
			expected:          "--name",
			expectedSegmented: []string{"--name"},
		},
		{
			name:        "Implicit false",
			errExpected: true,
		},
		{
			name:              "Implicit single optionFunc true",
			value:             true,
			optionFuncs:       []optionFunc{Single},
			expected:          "-name",
			expectedSegmented: []string{"-name"},
		},
		{
			name:        "Implicit single optionFunc false",
			errExpected: true,
		},
		{
			name:              "Quoted optionFunc true",
			value:             true,
			optionFuncs:       []optionFunc{Quoted},
			expected:          `--name "true"`,
			expectedSegmented: []string{"--name", `"true"`},
		},
		{
			name:              "Quoted optionFunc false",
			optionFuncs:       []optionFunc{Quoted},
			expected:          `--name "false"`,
			expectedSegmented: []string{"--name", `"false"`},
		},
		{
			name:              "Equal separator optionFunc true",
			value:             true,
			optionFuncs:       []optionFunc{EqualSeparator},
			expected:          "--name=true",
			expectedSegmented: []string{"--name=true"},
		},
		{
			name:              "Equal separator optionFunc false",
			optionFuncs:       []optionFunc{EqualSeparator},
			expected:          "--name=false",
			expectedSegmented: []string{"--name=false"},
		},
		{
			name:              "All optionFuncs true",
			value:             true,
			optionFuncs:       []optionFunc{Single, Quoted, EqualSeparator},
			expected:          `-name="true"`,
			expectedSegmented: []string{`-name="true"`},
		},
		{
			name:              "All optionFuncs false",
			optionFuncs:       []optionFunc{Single, Quoted, EqualSeparator},
			expected:          `-name="false"`,
			expectedSegmented: []string{`-name="false"`},
		},
		{
			name:              "Title optionFunc",
			value:             true,
			optionFuncs:       []optionFunc{Title},
			expected:          `--name True`,
			expectedSegmented: []string{"--name", "True"},
		},
		{
			name:              "Uppercase optionFunc",
			value:             true,
			optionFuncs:       []optionFunc{Uppercase},
			expected:          `--name TRUE`,
			expectedSegmented: []string{"--name", "TRUE"},
		},
		{
			name:              "Explicit optionFunc",
			value:             true,
			optionFuncs:       []optionFunc{Explicit},
			expected:          `--name true`,
			expectedSegmented: []string{"--name", "true"},
		},
		{
			name:              "Title optionFunc false",
			optionFuncs:       []optionFunc{Title},
			expected:          `--name False`,
			expectedSegmented: []string{"--name", "False"},
		},
		{
			name:              "Uppercase optionFunc false",
			optionFuncs:       []optionFunc{Uppercase},
			expected:          `--name FALSE`,
			expectedSegmented: []string{"--name", "FALSE"},
		},
		{
			name:              "Explicit optionFunc false",
			optionFuncs:       []optionFunc{Explicit},
			expected:          `--name false`,
			expectedSegmented: []string{"--name", "false"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bf, err := NewBoolFlag("name", testCase.value, testCase.optionFuncs...)
			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			}

			assert.NoError(t, err)
			s, err := bf.String()
			assert.NoError(t, err)

			assert.Equal(t, testCase.expected, s)

			sep, err := bf.Segmented()
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedSegmented, sep)

			assert.Equal(t, "name", bf.Name())
			assert.Equal(t, fmt.Sprintf("%v", testCase.value), bf.Value())
		})
	}

}
