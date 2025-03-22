package genflag

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalFlags(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		expectedFlags []string
		errExpected   bool
		matchOrder    bool
	}{
		{
			name: "Simple struct",
			input: struct {
				SingleOption   string `genflag:"single"`
				DoubleOption   string `genflag:""`
				Override       string `genflag:"overridden"`
				ImplicitSwitch bool   `genflag:""`
				ExplicitSwitch bool   `genflag:"explicit"`
				FieldIgnored   bool
				SingleQuoted   string `genflag:"quoted,single"`
			}{
				SingleOption:   "singlevalue",
				DoubleOption:   "doublevalue",
				Override:       "overriddenvalue",
				ImplicitSwitch: true,
				ExplicitSwitch: false,
				FieldIgnored:   false,
				SingleQuoted:   "singlequotedvalue",
			},
			expectedFlags: []string{
				"-singleoption singlevalue",
				"--doubleoption doublevalue",
				"--overridden overriddenvalue",
				"--implicitswitch",
				"--explicitswitch false",
				`-singlequoted "singlequotedvalue"`,
			},
		},
		{
			name: "Struct with flag slice",
			input: struct {
				Args []string `genflag:"arg"`
			}{
				Args: []string{
					"arg1",
					"arg2",
					"arg3",
					"arg4",
					"arg5",
				},
			},
			expectedFlags: []string{
				"--arg arg1",
				"--arg arg2",
				"--arg arg3",
				"--arg arg4",
				"--arg arg5",
			},
		},
		{
			name: "Struct with flag slice and single",
			input: struct {
				Args []string `genflag:"arg,single"`
			}{
				Args: []string{
					"arg1",
					"arg2",
					"arg3",
					"arg4",
					"arg5",
				},
			},
			expectedFlags: []string{
				"-arg arg1",
				"-arg arg2",
				"-arg arg3",
				"-arg arg4",
				"-arg arg5",
			},
		},
		{
			name: "Struct with key value",
			input: struct {
				KeyValue map[string]string `genflag:""`
			}{
				KeyValue: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			expectedFlags: []string{
				"--one two",
				"--three four",
			},
			matchOrder: false,
		},
		{
			name: "Struct with bool keys",
			input: struct {
				BoolOpts map[string]bool `genflag:"explicit"`
			}{
				BoolOpts: map[string]bool{
					"arg1": true,
					"arg2": false,
				},
			},
			expectedFlags: []string{
				"--arg1 true",
				"--arg2 false",
			},
			matchOrder: false,
		},
		{
			// TODO: Determine what to do here. Should be like the JSON
			// parser.
			name: "Struct with tagged nested struct without field tag",
			input: struct {
				TopLevel     string `genflag:""`
				NestedStruct NestedStructWithTaggedField
			}{
				TopLevel: "toplevelopt",
				NestedStruct: NestedStructWithTaggedField{
					TaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
				"--taggedfield value",
			},
		},
		{
			name: "Struct with tagged nested struct with field tag",
			input: struct {
				TopLevel     string                      `genflag:""`
				NestedStruct NestedStructWithTaggedField `genflag:""`
			}{
				TopLevel: "toplevelopt",
				NestedStruct: NestedStructWithTaggedField{
					TaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
				"--taggedfield value",
			},
		},
		{
			// TODO: Determine what to do here. Should be like the JSON
			// parser.
			name: "Struct with embedded nested struct without field tag",
			input: struct {
				TopLevel string `genflag:""`
				NestedStructWithTaggedField
			}{
				TopLevel: "toplevelopt",
				NestedStructWithTaggedField: NestedStructWithTaggedField{
					TaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
				"--taggedfield value",
			},
		},
		{
			name: "Struct with embedded nested struct with field tag",
			input: struct {
				TopLevel                    string `genflag:""`
				NestedStructWithTaggedField `genflag:""`
			}{
				TopLevel: "toplevelopt",
				NestedStructWithTaggedField: NestedStructWithTaggedField{
					TaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
				"--taggedfield value",
			},
		},
		{
			name: "Struct with untagged nested struct without field tag",
			input: struct {
				TopLevel     string `genflag:""`
				NestedStruct NestedStructWithoutTaggedField
			}{
				TopLevel: "toplevelopt",
				NestedStruct: NestedStructWithoutTaggedField{
					UntaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
			},
		},
		{
			name: "Struct with untagged nested struct with field tag",
			input: struct {
				TopLevel     string                         `genflag:""`
				NestedStruct NestedStructWithoutTaggedField `genflag:""`
			}{
				TopLevel: "toplevelopt",
				NestedStruct: NestedStructWithoutTaggedField{
					UntaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
			},
		},
		{
			name: "Struct with untagged embedded struct without field tag",
			input: struct {
				TopLevel string `genflag:""`
				NestedStructWithoutTaggedField
			}{
				TopLevel: "toplevelopt",
				NestedStructWithoutTaggedField: NestedStructWithoutTaggedField{
					UntaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
			},
		},
		{
			name: "Struct with untagged nested struct with field tag",
			input: struct {
				TopLevel                       string `genflag:""`
				NestedStructWithoutTaggedField `genflag:""`
			}{
				TopLevel: "toplevelopt",
				NestedStructWithoutTaggedField: NestedStructWithoutTaggedField{
					UntaggedField: "value",
				},
			},
			expectedFlags: []string{
				"--toplevel toplevelopt",
			},
		},
		{
			name:  "Custom marshaler",
			input: newCustomMarshaler("custom", "marshaler"),
			expectedFlags: []string{
				"--custom marshaler",
			},
		},
		{
			name: "Nested custom marshaler",
			input: struct {
				Custom customMarshaler
				Other  string `genflag:""`
			}{
				Custom: newCustomMarshaler("custom", "marshaler"),
				Other:  "hello",
			},
			expectedFlags: []string{
				"--custom marshaler",
				"--other hello",
			},
		},
		{
			name: "Nested custom pointer marshaler",
			input: struct {
				Custom *customMarshaler
				Other  string `genflag:""`
			}{
				Custom: newCustomMarshalerPtr("custom", "marshaler"),
				Other:  "hello",
			},
			expectedFlags: []string{
				"--custom marshaler",
				"--other hello",
			},
		},
		{
			name: "Nil pointer on field",
			input: struct {
				Hello *string
			}{
				Hello: nil,
			},
			expectedFlags: []string{},
		},
		{
			name: "Coexists with other struct tags",
			input: struct {
				Hello string `genflag:"" json:"hi,omitempty"`
			}{
				Hello: "hi",
			},
			expectedFlags: []string{
				"--hello hi",
			},
		},
		{
			name: "Ignores empty strings",
			input: struct {
				Hello string `genflag:""`
				Other string `genflag:""`
			}{
				Hello: "hi",
				Other: "",
			},
			expectedFlags: []string{
				"--hello hi",
			},
		},
		{
			name: "Handles pointers",
			input: struct {
				Switch    *bool     `genflag:""`
				Opt       *string   `genflag:""`
				NilSwitch *bool     `genflag:""`
				NilOpt    *string   `genflag:""`
				OptList   []*string `genflag:""`
			}{
				Switch: boolToPtr(true),
				Opt:    stringToPtr("arg"),
				OptList: []*string{
					stringToPtr("opt1"),
					stringToPtr("opt2"),
				},
			},
			expectedFlags: []string{
				"--switch",
				"--opt arg",
				"--optlist opt1",
				"--optlist opt2",
			},
		},
		{
			name: "Handles empty interfaces",
			input: struct {
				SwitchPtr            interface{}            `genflag:""`
				Switch               interface{}            `genflag:""`
				StringList           []interface{}          `genflag:""`
				StringMap            map[string]interface{} `genflag:""`
				BoolMap              map[string]interface{} `genflag:"explicit"`
				StringPtrMap         map[string]interface{} `genflag:""`
				BoolPtrMap           map[string]interface{} `genflag:""`
				CombinedMap          map[string]interface{} `genflag:""`
				EmptyWithMap         interface{}            `genflag:""`
				Ignored              string
				UntaggedNestedStruct interface{} `genflag:""`
				TaggedNestedStruct   interface{} `genflag:""`
			}{
				SwitchPtr: boolToPtr(true),
				Switch:    true,
				StringList: []interface{}{
					"opt1",
					"opt2",
				},
				StringMap: map[string]interface{}{
					"opt1": "opt2",
				},
				BoolMap: map[string]interface{}{
					"opt3": true,
					"opt4": false,
				},
				StringPtrMap: map[string]interface{}{
					"opt5": stringToPtr("opt6"),
				},
				BoolPtrMap: map[string]interface{}{
					"opt7": boolToPtr(true),
				},
				CombinedMap: map[string]interface{}{
					"opt8":  stringToPtr("opt9"),
					"opt10": boolToPtr(true),
				},
				EmptyWithMap: map[string]interface{}{
					"opt11": "opt12",
					"opt13": stringToPtr("opt14"),
					"opt15": true,
					"opt16": boolToPtr(true),
				},
				Ignored: "should be ignored",
				UntaggedNestedStruct: NestedStructWithoutTaggedField{
					UntaggedField: "should be ignored",
				},
				TaggedNestedStruct: NestedStructWithTaggedField{
					TaggedField: "taggedvalue",
				},
			},
			expectedFlags: []string{
				"--switchptr",
				"--switch",
				"--stringlist opt1",
				"--stringlist opt2",
				"--opt1 opt2",
				"--opt3 true",
				"--opt4 false",
				"--opt5 opt6",
				"--opt7",
				"--opt8 opt9",
				"--opt10",
				"--opt11 opt12",
				"--opt13 opt14",
				"--opt15",
				"--opt16",
				"--taggedfield taggedvalue",
			},
		},
		{
			name:        "Errors on listed structs",
			input:       newListStruct(5, "opt"),
			errExpected: true,
		},
		{
			name:        "Errors on listed struct pointers",
			input:       newListStructPtr(5, "opt"),
			errExpected: true,
		},
		{
			name: "Key value map",
			input: map[string]string{
				"opt1": "opt2",
				"opt3": "opt4",
			},
			expectedFlags: []string{
				"--opt1 opt2",
				"--opt3 opt4",
			},
		},
		{
			name: "Switch map",
			input: map[string]bool{
				"opt1": true,
				"opt2": true,
			},
			expectedFlags: []string{
				"--opt1",
				"--opt2",
			},
		},
		{
			name: "Nil values in struct",
			input: struct {
				Field1    *bool             `genfiag:""`
				Items     []string          `genflag:""`
				Switches  map[string]bool   `genflag:""`
				KeyValues map[string]string `genflag:""`
			}{
				Field1:    nil,
				Items:     nil,
				Switches:  nil,
				KeyValues: nil,
			},
			expectedFlags: []string{},
		},
		{
			name: "Errors on tagged list of structs within struct",
			input: struct {
				Nested  []listStruct `genflag:""`
				Another string       `genflag:""`
			}{
				Nested:  newListStruct(3, "opt"),
				Another: "value",
			},
			errExpected: true,
		},
		{
			name: "Errors on tagged list of struct pointers within struct",
			input: struct {
				Nested  []*listStruct `genflag:""`
				Another string        `genflag:""`
			}{
				Nested:  newListStructPtr(3, "opt"),
				Another: "value",
			},
			errExpected: true,
		},
		{
			name: "List of structs without tag are ignored",
			input: struct {
				Nested  []listStruct
				Another string `genflag:""`
			}{
				Nested:  newListStruct(3, "opt"),
				Another: "value",
			},
			expectedFlags: []string{
				"--another value",
			},
		},
		{
			name: "List of struct pointers without tag are ignored",
			input: struct {
				Nested  []*listStruct
				Another string `genflag:""`
			}{
				Nested:  newListStructPtr(3, "opt"),
				Another: "value",
			},
			expectedFlags: []string{
				"--another value",
			},
		},
		{
			name:        "Errors on top-level nil",
			input:       nil,
			errExpected: true,
		},
		{
			name: "Errors on unexported fields",
			input: struct {
				field string `genflag:""`
			}{
				field: "hello",
			},
			errExpected: true,
		},
		{
			name: "Errors on string list",
			input: []string{
				"opt1",
			},
			errExpected: true,
		},
		{
			name: "Errors on nested string slices",
			input: struct {
				Item [][]string `genflag:""`
			}{
				Item: [][]string{
					{
						"opt1",
						"opt2",
					},
				},
			},
			errExpected: true,
		},
		{
			name: "Errors on bool list",
			input: struct {
				Items []bool `genflag:""`
			}{
				Items: []bool{true, false, true},
			},
			errExpected: true,
		},
		{
			name: "Errors on string slice list not having unique values",
			input: struct {
				Items []string `genflag:""`
			}{
				Items: []string{"val1", "val1"},
			},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			results, err := MarshalFlags(testCase.input)

			if testCase.errExpected {
				assert.Error(t, err)
				t.Log(err)
				return
			}

			assert.NoError(t, err)

			actual := []string{}
			for _, result := range results {
				flag, err := result.String()
				assert.NoError(t, err)
				actual = append(actual, flag)
			}

			if testCase.matchOrder {
				assert.Equal(t, testCase.expectedFlags, actual)
			} else {
				assert.Equal(t, stringSliceToMap(testCase.expectedFlags), stringSliceToMap(actual))
			}
		})
	}
}

func TestValidateFlagList(t *testing.T) {
	testCases := []struct {
		name        string
		flags       []Flag
		errExpected bool
	}{
		{
			name: "List of bool flags with different name",
			flags: []Flag{
				boolFlag{
					name:  "opt",
					value: false,
				},
				boolFlag{
					name:  "another",
					value: false,
				},
			},
		},
		{
			name: "List of bool flags with same name collide",
			flags: []Flag{
				boolFlag{
					name:  "opt",
					value: false,
				},
				boolFlag{
					name:  "opt",
					value: false,
				},
			},
			errExpected: true,
		},
		{
			name: "List of string flags with different names and values",
			flags: []Flag{
				stringFlag{
					name:  "opt",
					value: "value",
				},
				stringFlag{
					name:  "anotheropt",
					value: "anothervalue",
				},
			},
		},
		{
			name: "List of string flags with same names and different values",
			flags: []Flag{
				stringFlag{
					name:  "opt",
					value: "value",
				},
				stringFlag{
					name:  "opt",
					value: "anothervalue",
				},
			},
		},
		{
			name: "List of string flags with same names and same values",
			flags: []Flag{
				stringFlag{
					name:  "opt",
					value: "value",
				},
				stringFlag{
					name:  "opt",
					value: "value",
				},
			},
			errExpected: true,
		},
		{
			name: "List of string flags with same names and same values",
			flags: []Flag{
				stringFlag{
					name:  "opt",
					value: "value",
				},
				stringFlag{
					name:  "opt",
					value: "value",
				},
			},
			errExpected: true,
		},
		{
			name:  "List flags with different values",
			flags: NewListFlagsOrDie("opt", []string{"val1", "val2", "val3"}),
		},
		{
			name: "List flags with same values",
			// Must call this twice because NewListFlags() checks for
			// duplicate values.
			flags:       append(NewListFlagsOrDie("opt", []string{"val1"}), NewListFlagsOrDie("opt", []string{"val1"})...),
			errExpected: true,
		},
		{
			name: "List flags with different values collides with string flags with same values",
			flags: append(NewListFlagsOrDie("opt", []string{"val1", "val2", "val3"}), stringFlag{
				name:  "opt",
				value: "val1",
			}),
			errExpected: true,
		},
		{
			name: "List flags with different values collides with bool flag with same name",
			flags: append(NewListFlagsOrDie("opt", []string{"val1", "val2", "val3"}), boolFlag{
				name: "opt",
			}),
			errExpected: true,
		},
		{
			name: "Switch flags have different names",
			flags: NewSwitchFlagsOrDie(map[string]bool{
				"opt1": true,
				"opt2": true,
				"opt3": true,
			}),
		},
		{
			name: "Switch flags collide with bool flags with the same name",
			flags: append(NewSwitchFlagsOrDie(map[string]bool{
				"opt": true,
			}), boolFlag{
				name:  "opt",
				value: false,
			}),
			errExpected: true,
		},
		{
			name: "Switch flags collide with string flags with the same name",
			flags: append(NewSwitchFlagsOrDie(map[string]bool{
				"opt": true,
			}), stringFlag{
				name:  "opt",
				value: "avalue",
			}),
			errExpected: true,
		},
		{
			name: "Key value flags with different keys and values",
			flags: NewKeyValueFlagsOrDie(map[string]string{
				"opt1": "val1",
				"opt2": "val2",
			}),
		},
		{
			name: "Key value flags collide with string flags with the same names and values",
			flags: append(NewKeyValueFlagsOrDie(map[string]string{
				"opt1": "val1",
				"opt2": "val2",
			}), stringFlag{
				name:  "opt1",
				value: "val1",
			}),
			errExpected: true,
		},
		{
			name: "Key value flags collide with list flags with the same names and values",
			flags: append(NewKeyValueFlagsOrDie(map[string]string{
				"opt1": "val1",
				"opt2": "val2",
			}), NewListFlagsOrDie("opt1", []string{"val1", "val2"})...),
			errExpected: true,
		},
		{
			name: "Bool flags collide with string flag with same name",
			flags: []Flag{
				boolFlag{
					name:  "opt1",
					value: true,
				},
				stringFlag{
					name:  "opt1",
					value: "different",
				},
			},
			errExpected: true,
		},
		{
			name: "Custom flag type collides with string flag value",
			flags: []Flag{
				newCustomMarshaler("opt", "value"),
				stringFlag{
					name:  "opt",
					value: "value",
				},
			},
			errExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &marshaler{}
			err := m.validateFlagList(testCase.flags)

			if err != nil {
				t.Log(err)
			}

			if testCase.errExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type NestedStructWithTaggedField struct {
	TaggedField string `genflag:""`
}

type NestedStructWithoutTaggedField struct {
	UntaggedField string
}

type listStruct struct {
	Field string `genflag:""`
}

func newListStructPtr(n int, val string) []*listStruct {
	out := []*listStruct{}

	for i := 1; i <= n; i++ {
		out = append(out, &listStruct{
			Field: fmt.Sprintf("%s-%d", val, i),
		})
	}

	return out
}

func newListStruct(n int, val string) []listStruct {
	out := []listStruct{}

	for i := 1; i <= n; i++ {
		out = append(out, listStruct{
			Field: fmt.Sprintf("%s-%d", val, i),
		})
	}

	return out
}

type customMarshaler struct {
	name  string
	value string
}

func newCustomMarshaler(name, value string) customMarshaler {
	return customMarshaler{
		name:  name,
		value: value,
	}
}

func newCustomMarshalerPtr(name, value string) *customMarshaler {
	cm := newCustomMarshaler(name, value)
	return &cm
}

func (c customMarshaler) Name() string {
	return c.name
}

func (c customMarshaler) Value() string {
	return c.value
}

func (c customMarshaler) String() (string, error) {
	val, err := c.Segmented()
	return strings.Join(val, " "), err
}

func (c customMarshaler) Segmented() ([]string, error) {
	return []string{c.getDashedName(), c.value}, nil
}

func (c customMarshaler) getDashedName() string {
	return fmt.Sprintf("--%s", c.name)
}

func (c customMarshaler) MarshalFlags() ([]Flag, error) {
	return []Flag{c}, nil
}

func boolToPtr(val bool) *bool {
	return &val
}

func stringToPtr(s string) *string {
	return &s
}
