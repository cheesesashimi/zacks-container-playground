package genflag

import "fmt"

// A listflag is just a type alias for string flags. The difference is that the
// factory which produces these assumes the same name for the flag, but
// provides different values. This is useful for CLI options which expect many
// flags such as, "--option val1 --option val2 --option val3".
type listFlag struct {
	stringFlag
}

// Constructs a new validated listflag from a stringflag.
func newListFlag(name, value string, optionFuncs ...optionFunc) (listFlag, error) {
	sf, err := newStringFlag(name, value, optionFuncs...)
	return listFlag{sf}, err
}

// Iterates through all of the provided values, instantiating one listflag for
// each value provided. Each listflag has the same name, but a different value.
// Values must be unique.
func NewListFlags(name string, values []string, optionFuncs ...optionFunc) ([]Flag, error) {
	out := []Flag{}

	seen := map[string]struct{}{}

	for _, value := range values {
		if _, ok := seen[value]; ok {
			return nil, fmt.Errorf("values must be unique, found %q more than once", value)
		}

		seen[value] = struct{}{}

		f, err := newListFlag(name, value, optionFuncs...)
		if err != nil {
			return nil, err
		}

		out = append(out, f)
	}

	return out, nil
}

func NewListFlagsOrDie(name string, values []string, optionFuncs ...optionFunc) []Flag {
	f, err := NewListFlags(name, values, optionFuncs...)
	if err != nil {
		panic(err)
	}
	return f
}

// A keyValueFlag is a type alias for string flags. The difference is that the
// factory which produces these assumes the map key name for the flag name, and
// the map value for the flag value. This is useful for situations where one
// wants to be more dynamic about what flags are being instantiated.
type keyValueFlag struct {
	stringFlag
}

// Constructs a new validated keyValueFlag from a stringflag.
func newKeyValueFlag(name string, value string, optionFuncs ...optionFunc) (keyValueFlag, error) {
	sf, err := newStringFlag(name, value, optionFuncs...)
	return keyValueFlag{sf}, err
}

// Iterates through all of the provided keys and values, instantiating one
// keyValueFlag for each value provided. Each keyValueFlag has the name of the
// map key with the value provided by the map value.
func NewKeyValueFlags(items map[string]string, optionFuncs ...optionFunc) ([]Flag, error) {
	out := []Flag{}

	for key, value := range items {
		f, err := newKeyValueFlag(key, value, optionFuncs...)
		if err != nil {
			return nil, err
		}

		out = append(out, f)
	}

	return out, nil
}

func NewKeyValueFlagsOrDie(items map[string]string, optionFuncs ...optionFunc) []Flag {
	f, err := NewKeyValueFlags(items, optionFuncs...)
	if err != nil {
		panic(err)
	}
	return f
}

// A switchFlag is a type alias for boolean flags. The difference is that the
// factory which produces these assumes the map key for the name, and the value
// for the flags' value.
type switchFlag struct {
	boolFlag
}

// Instantiates new switchFlag from a boolean flag.
func newSwitchFlag(name string, value bool, optionFuncs ...optionFunc) (switchFlag, error) {
	bf, err := newBoolFlag(name, value, optionFuncs...)
	return switchFlag{bf}, err
}

// Iterates through all of the provided keys and values, instantiating one
// switchFlag for each key provided.
func NewSwitchFlags(items map[string]bool, optionFuncs ...optionFunc) ([]Flag, error) {
	out := []Flag{}

	for key, val := range items {
		f, err := newSwitchFlag(key, val, optionFuncs...)
		if err != nil {
			return nil, err
		}

		out = append(out, f)
	}

	return out, nil
}

func NewSwitchFlagsOrDie(items map[string]bool, optionFuncs ...optionFunc) []Flag {
	f, err := NewSwitchFlags(items, optionFuncs...)
	if err != nil {
		panic(err)
	}

	return f
}
