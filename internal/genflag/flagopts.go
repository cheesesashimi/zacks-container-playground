package genflag

import (
	"fmt"
	"strings"
)

// The basic flag options struct. This holds and implements the
// basic functionality as described within.
type flagOpts struct {
	// This value indicates whether the name and value of the flag
	// should be separated by a space or an equal sign. If not set,
	// will default to an empty space.
	separator string
	// If true, the flag name will be preceded by a single dash
	// (-name). By default, the flag name will be preceded by a
	// double dash (--name).
	single bool
	// If true, the value of the flag will be wrapped in double
	// quotes, e.g., "value".
	quoted bool
}

// Returns a pointer to this instance.
func (f *flagOpts) getFlagOpts() *flagOpts {
	return f
}

// Validates the options given as well as the value provided.
func (f flagOpts) validate(value string) error {
	if value == "" && f.quoted {
		return fmt.Errorf("cannot quote empty value")
	}

	if value == "" && f.separator != "" {
		return fmt.Errorf("cannot separate empty value")
	}

	if err := f.validateSeparator(); err != nil {
		return err
	}

	return nil
}

// Renders the flag with the name and value in seaprate string
// slice elements, e.g., []string{"--name", "value"}.
//
// This is subject to the various options such as the separator
// being set, quoting, etc. If no separator is set or the
// separator is an empty space, the elements will be in separate
// parts of the string slice. If the separator is set to an equal
// sign (=), then the values will be placed into a single element
// within the string slice, e.g.: []string{"--name=value"}.
//
// This will perform validation of the value based upon the
// current settings of the flagOpts struct.
func (f flagOpts) renderSegmented(name, value string) ([]string, error) {
	if err := f.validate(value); err != nil {
		return nil, err
	}

	out := []string{f.getDashes() + name}

	if value == "" {
		return out, nil
	}

	val := f.getValue(value)

	if f.separator == "" || f.separator == " " {
		return append(out, val), nil
	}

	out[0] = out[0] + f.getSeparator() + val
	return out, nil
}

// Renders the flag with the name and value in a single string,
// along with the various modification options.
//
// This will perform validation of the value based upon the
// current settings of the flagOpts struct.
func (f flagOpts) render(name, value string) (string, error) {
	if err := f.validate(value); err != nil {
		return "", err
	}

	val := f.getValue(value)

	out := []string{f.getDashes(), name}

	if val != "" {
		out = append(out, []string{f.getSeparator(), val}...)
	}

	return strings.Join(out, ""), nil
}

// Determines if a single dash or a double-dash should be
// returned.
func (f flagOpts) getDashes() string {
	if f.single {
		return "-"
	}

	return "--"
}

// Returns the value either double-quoted or not.
func (f flagOpts) getValue(value string) string {
	if f.quoted {
		return fmt.Sprintf("%q", value)
	}

	return value
}

// Gets the separator, defaulting to an empty space.
func (f flagOpts) getSeparator() string {
	if f.separator == "" {
		return " "
	}

	return f.separator
}

// Determines whether the provided separator is valid. Valid
// separators include an empty string, an empty space, and an
// equal sign.
func (f flagOpts) validateSeparator() error {
	valid := map[string]struct{}{
		"=": {},
		" ": {},
		"":  {},
	}

	if _, ok := valid[f.separator]; !ok {
		return fmt.Errorf("invalid separator %q", f.separator)
	}

	return nil
}
