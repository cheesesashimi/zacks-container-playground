package genflag

import (
	"fmt"
	"strings"
)

// Represents the type of flgas which may utilize optionFuncs for
// configuration. This is internal-only.
type internalFlag interface {
	Flag
	getFlagOpts() *flagOpts
}

// Options funcs which may be passed into the various flag constructors to
// control options such as the separator, whether the flag value is quoted, has
// a single dash, etc.
type optionFunc func(internalFlag) error

// When active, ensures that the flag has a single leading dash.
func Single(f internalFlag) error {
	fo := f.getFlagOpts()
	fo.single = true
	return nil
}

// When active, ensures that the flag value is surrounded in double quotes.
func Quoted(f internalFlag) error {
	fo := f.getFlagOpts()
	fo.quoted = true

	_, ok := f.(*boolFlag)
	if !ok {
		return nil
	}

	return Explicit(f)
}

// When active, ensures that the flag name and value are separated by an equal
// sign (=).
func EqualSeparator(f internalFlag) error {
	tf := Separator("=")
	return tf(f)
}

// When active, allows one to set a specific seprator, subject to validation.
func Separator(sep string) optionFunc {
	return func(f internalFlag) error {
		fo := f.getFlagOpts()
		fo.separator = sep

		_, ok := f.(*boolFlag)
		if !ok {
			return nil
		}

		return Explicit(f)
	}
}

// When active, forces a boolean flag's value to explicitly say true or false.
func Explicit(f internalFlag) error {
	bf, ok := f.(*boolFlag)
	if !ok {
		return fmt.Errorf("explicit optionFunc only available on boolFlag")
	}

	bf.explicit = true
	return nil
}

// When active, forces a boolean flag's value to explicitly say TRUE or FALSE.
func Uppercase(f internalFlag) error {
	return boolOptionFunc(strings.ToUpper)(f)
}

// When active, forces a boolean flag's value to exlicitly say True or False.
func Title(f internalFlag) error {
	return boolOptionFunc(strings.Title)(f)
}

// Ensures that a given transformer function can only be applied to a boolean
// flag.
func boolOptionFunc(tf func(string) string) optionFunc {
	return func(f internalFlag) error {
		bf, ok := f.(*boolFlag)
		if !ok {
			return fmt.Errorf("boolean optionFunc only available on boolFlag")
		}

		bf.transform = tf
		return Explicit(bf)
	}
}

// Applies the given optionfuncs to the flagOpts struct retrieved from the
// passed flag..
func applyOptionFuncs(f internalFlag, optionFuncs ...optionFunc) error {
	for _, optionFunc := range optionFuncs {
		if err := optionFunc(f); err != nil {
			return err
		}
	}

	return nil
}

// Constructs a flagOpts struct from the given optionFuncs.
func flagOptsFromOptionFuncs(optionFuncs ...optionFunc) (flagOpts, error) {
	sf := stringFlag{}

	if err := applyOptionFuncs(&sf, optionFuncs...); err != nil {
		return sf.flagOpts, err
	}

	return sf.flagOpts, nil
}
