package genflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenflagTagParser(t *testing.T) {
	testCases := []struct {
		testName    string
		tagInput    string
		errExpected bool
		expected    genflagTagParser
	}{
		{
			testName: "Just name override",
			tagInput: "overriddenname",
			expected: genflagTagParser{
				name:    "overriddenname",
				setOpts: getSetOpts([]genflagTagOpt{}),
			},
		},
		{
			testName: "Keyword only",
			tagInput: "quoted",
			expected: genflagTagParser{
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Override with keyword",
			tagInput: "newname,single",
			expected: genflagTagParser{
				name:    "newname",
				setOpts: getSetOpts([]genflagTagOpt{genFlagSingleOpt}),
			},
		},
		{
			testName: "Title casing",
			tagInput: "titlecase",
			expected: genflagTagParser{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitBoolTitleCase}),
			},
		},
		{
			testName: "Uppercasing",
			tagInput: "uppercase",
			expected: genflagTagParser{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitBoolUppercase}),
			},
		},
		{
			testName: "Explicit boolean opt",
			tagInput: "explicit",
			expected: genflagTagParser{
				setOpts: getSetOpts([]genflagTagOpt{genFlagExplicitOpt}),
			},
		},
		{
			testName: "Quoted opt",
			tagInput: "quoted",
			expected: genflagTagParser{
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Quoted opt with quoted override",
			tagInput: "quoted,quoted",
			expected: genflagTagParser{
				name:    "quoted",
				setOpts: getSetOpts([]genflagTagOpt{genFlagQuoted}),
			},
		},
		{
			testName: "Multiple keywords set",
			tagInput: "quoted,explicit,titlecase,single",
			expected: genflagTagParser{
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
			expected: genflagTagParser{
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

			gfo, err := newGenflagTagParser(fieldName, testCase.tagInput)
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

func TestGenflagTagParserBoolean(t *testing.T) {
	t.Run("Implicit true", func(t *testing.T) {
		gfo, err := newGenflagTagParser("switch", "")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(true)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch", content)
	})

	t.Run("Implicit false", func(t *testing.T) {
		gfo, err := newGenflagTagParser("switch", "")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(false)
		assert.NoError(t, err)
		assert.Len(t, f, 0)
	})

	t.Run("Explicit true", func(t *testing.T) {
		gfo, err := newGenflagTagParser("switch", "explicit")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(true)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch true", content)
	})

	t.Run("Explicit false", func(t *testing.T) {
		gfo, err := newGenflagTagParser("switch", "explicit")
		assert.NoError(t, err)

		f, err := gfo.newBoolFlagWithName(false)
		assert.NoError(t, err)
		assert.Len(t, f, 1)

		content, err := f[0].String()
		assert.NoError(t, err)

		assert.Equal(t, "--switch false", content)
	})
}
