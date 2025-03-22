package genflag

import (
	"fmt"
	"strings"
)

type genflagTagOpt string

const (
	genFlagKeyName               string        = "genflag"
	genFlagEqualSeparated        genflagTagOpt = "equaled"
	genFlagExplicitBoolTitleCase genflagTagOpt = "titlecase"
	genFlagExplicitBoolUppercase genflagTagOpt = "uppercase"
	genFlagExplicitOpt           genflagTagOpt = "explicit"
	genFlagQuoted                genflagTagOpt = "quoted"
	genFlagSingleOpt             genflagTagOpt = "single"
)

// Implements a parser for the genflag struct tags.
type genFlagOpts struct {
	// The name of the struct field. For example:
	// type astruct struct {
	//     FieldName string `genflag:""`
	// }
	//
	// The name would be "FieldName".
	name string
	// Holds the options that are set.
	setOpts map[genflagTagOpt]bool
}

// Constructs a new genFlagOpts instance, validating what options are set in
// the process.
func newGenFlagOpts(fieldName, in string) (*genFlagOpts, error) {
	g := &genFlagOpts{
		name:    strings.ToLower(fieldName),
		setOpts: getSetOpts([]genflagTagOpt{}),
	}

	if in == "" {
		return g, nil
	}

	split := strings.Split(in, ",")

	invalid := map[string]struct{}{}

	for _, item := range split {
		if item == "" || item == " " || strings.Contains(item, " ") {
			return nil, fmt.Errorf("empty space not allowed")
		}

		isSet, ok := g.setOpts[genflagTagOpt(item)]
		if ok {
			if isSet {
				invalid[item] = struct{}{}
			} else {
				g.setOpts[genflagTagOpt(item)] = true
			}
		} else {
			invalid[item] = struct{}{}
		}
	}

	if g.setOpts[genFlagExplicitBoolUppercase] && g.setOpts[genFlagExplicitBoolTitleCase] {
		return nil, fmt.Errorf("only one of %v may be used, not both", []genflagTagOpt{genFlagExplicitBoolUppercase, genFlagExplicitBoolTitleCase})
	}

	if len(invalid) == 0 {
		return g, nil
	}

	invalidKeywords := stringMapToSlice(invalid)

	// If there is only one "invalid" keyword and it is found in the first
	// position, this should be used as the name override for the field name.
	if len(invalid) == 1 && invalidKeywords[0] == split[0] {
		g.name = invalidKeywords[0]
		return g, nil
	}

	return nil, fmt.Errorf("found multiple invalid keywords: %v", invalidKeywords)
}

// Gets all of the optionFuncs that correspond to the matching keywords.
func (g *genFlagOpts) getOptionFuncs() []optionFunc {
	setOpt := []genflagTagOpt{}

	for opt, isSet := range g.setOpts {
		if isSet {
			setOpt = append(setOpt, opt)
		}
	}

	return getMatchingOptionFuncs(setOpt)
}

// Constructs a new string flag with the given value, optionfuncs, and either
// the field name or the overridden name.
func (g *genFlagOpts) newStringFlagWithName(val string) (Flag, error) {
	return NewStringFlag(g.name, val, g.getOptionFuncs()...)
}

// Constructs a new boolean flag with the given value, optionfuncs, and either
// the field name or the overridden name.
func (g *genFlagOpts) newBoolFlagWithName(val bool) (Flag, error) {
	return NewBoolFlag(g.name, val, g.getOptionFuncs()...)
}

// Constructs list flags with the given values, optionfuncs, and either the
// field name or the overridden name.
func (g *genFlagOpts) newListFlag(values []string) ([]Flag, error) {
	return NewListFlags(g.name, values, g.getOptionFuncs()...)
}

// Constructs key/value flags with the given keys / values and optionfuncs.
func (g *genFlagOpts) newKeyValueFlag(items map[string]string) ([]Flag, error) {
	return NewKeyValueFlags(items, g.getOptionFuncs()...)
}

// Constructs switch flags with the given keys / values and optionfuncs.
func (g *genFlagOpts) newSwitchFlag(items map[string]bool) ([]Flag, error) {
	return NewSwitchFlags(items, g.getOptionFuncs()...)
}

// Returns which options in the input are set.
func getSetOpts(input []genflagTagOpt) map[genflagTagOpt]bool {
	out := map[genflagTagOpt]bool{
		genFlagEqualSeparated:        false,
		genFlagExplicitBoolTitleCase: false,
		genFlagExplicitBoolUppercase: false,
		genFlagExplicitOpt:           false,
		genFlagQuoted:                false,
		genFlagSingleOpt:             false,
	}

	for _, item := range input {
		_, ok := out[genflagTagOpt(item)]
		out[item] = ok
	}

	return out
}

// Maps the keywords to the optionfuncs which implement them.
func getOptionFuncOptMapping() map[genflagTagOpt]optionFunc {
	return map[genflagTagOpt]optionFunc{
		genFlagExplicitOpt:           Explicit,
		genFlagSingleOpt:             Single,
		genFlagExplicitBoolUppercase: Uppercase,
		genFlagExplicitBoolTitleCase: Title,
		genFlagQuoted:                Quoted,
		genFlagEqualSeparated:        EqualSeparator,
	}
}

// Given a list of input keywords, returns a list of matching optionfuncs.
func getMatchingOptionFuncs(input []genflagTagOpt) []optionFunc {
	mapping := getOptionFuncOptMapping()
	out := []optionFunc{}

	for _, item := range input {
		tf, ok := mapping[item]
		if ok {
			out = append(out, tf)
		}
	}

	return out
}
