package genflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenFlagOpts(t *testing.T) {
	testCases := []struct {
		testName    string
		tagInput    string
		errExpected bool
		expected    genFlagOpts
	}{
		{
			testName: "Just name override",
			tagInput: "overriddenname",
			expected: genFlagOpts{
				name:    "overriddenname",
				setOpts: getSetOpts([]genflagTagOpt{}),
			},
		},
		{
			testName: "Keyword only",
			tagInput: "quoted",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Override with keyword",
			tagInput: "newname,single",
			expected: genFlagOpts{
				name:    "newname",
				setOpts: getSetOpts([]genflagTagOpt{genFlagSingleOpt}),
			},
		},
		{
			testName: "Title casing",
			tagInput: "titlecase",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitBoolTitleCase}),
			},
		},
		{
			testName: "Uppercasing",
			tagInput: "uppercase",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitBoolUppercase}),
			},
		},
		{
			testName: "Explicit boolean opt",
			tagInput: "explicit",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitOpt}),
			},
		},
		{
			testName: "Quoted opt",
			tagInput: "quoted",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Quoted opt with quoted override",
			tagInput: "quoted,quoted",
			expected: genFlagOpts{
				name:    "quoted",
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Multiple keywords set",
			tagInput: "quoted,explicit,titlecase,single",
			expected: genFlagOpts{
				setOpts: getSetOpts([]genflagTagOpt{
					genFlagQuoted,
					genFlagExplicitOpt,
					genFlagExplicitBoolTitleCase,
					genFlagSingleOpt,
				}),
			},
		},
		{
			testName: "Multiple keywords set with override",
			tagInput: "newname,quoted,explicit,titlecase,single",
			expected: genFlagOpts{
				name: "newname",
				setOpts: getSetOpts([]genflagTagOpt{
					genFlagQuoted,
					genFlagExplicitOpt,
					genFlagExplicitBoolTitleCase,
					genFlagSingleOpt,
				}),
			},
		},
		{
			testName:    "Errors on leading or trailing spaces",
			tagInput:    " newname ",
			errExpected: true,
		},
		{
			testName:    "Errors on leading comma",
			tagInput:    ",newname,single",
			errExpected: true,
		},
		{
			testName:    "Errors on trailing comma",
			tagInput:    "newname,single,",
			errExpected: true,
		},
		{
			testName:    "Errors on middle empty",
			tagInput:    "newname,,single",
			errExpected: true,
		},
		{
			testName:    "Errors on leading empty space",
			tagInput:    " ,newname,single",
			errExpected: true,
		},
		{
			testName:    "Errors on trailing empty space",
			tagInput:    "newname,single, ",
			errExpected: true,
		},
		{
			testName:    "Errors on middle space",
			tagInput:    "newname, ,single",
			errExpected: true,
		},
		{
			testName:    "Errors on multiple leading empty space",
			tagInput:    "   ,newname,single",
			errExpected: true,
		},
		{
			testName:    "Errors on multiple trailing empty space",
			tagInput:    "newname,single,   ",
			errExpected: true,
		},
		{
			testName:    "Errors on multiple middle empty space",
			tagInput:    "newname,   ,single",
			errExpected: true,
		},
		{
			testName:    "override only allowed in first position",
			tagInput:    "single,titlecase,newname",
			errExpected: true,
		},
		{
			testName:    "only one transform allowed",
			tagInput:    "single,uppercase,titlecase",
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			fieldName := "OldName"

			if testCase.expected.name == "" {
				testCase.expected.name = "oldname"
			}

			gfo, err := newGenFlagOpts(fieldName, testCase.tagInput)
			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, *gfo)
		})
	}
}

func TestGenFlagOptsBoolean(t *testing.T) {
	t.Run("Implicit true", func(t *testing.T) {
		gfo, err := newGenFlagOpts("switch", "")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(true)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch", content)
	})

	t.Run("Implicit false", func(t *testing.T) {
		gfo, err := newGenFlagOpts("switch", "")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(false)
		assert.NoError(t, err)
		assert.Len(t, f, 0)
	})

	t.Run("Explicit true", func(t *testing.T) {
		gfo, err := newGenFlagOpts("switch", "explicit")

		f, err := gfo.newBoolFlagWithName(true)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch true", content)
	})

	t.Run("Explicit false", func(t *testing.T) {
		gfo, err := newGenFlagOpts("switch", "explicit")

		f, err := gfo.newBoolFlagWithName(false)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch false", content)
	})
}
