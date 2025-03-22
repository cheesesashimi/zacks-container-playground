package genflag

import (
	"fmt"
	"strings"
)

type internalFlag interface {
	Flag
	getFlagOpts() *flagOpts
}

type optionFunc func(internalFlag) error

func Single(f internalFlag) error {
	fo := f.getFlagOpts()
	fo.single = true
	return nil
}

func Quoted(f internalFlag) error {
	fo := f.getFlagOpts()
	fo.quoted = true

	_, ok := f.(*boolFlag)
	if !ok {
		return nil
	}

	return Explicit(f)
}

func EqualSeparator(f internalFlag) error {
	tf := Separator("=")
	return tf(f)
}

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

func Explicit(f internalFlag) error {
	bf, ok := f.(*boolFlag)
	if !ok {
		return fmt.Errorf("explicit optionFunc only available on boolFlag")
	}

	bf.explicit = true
	return nil
}

func Uppercase(f internalFlag) error {
	return boolOptionFunc(strings.ToUpper)(f)
}

func Title(f internalFlag) error {
	return boolOptionFunc(strings.Title)(f)
}

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

func applyOptionFuncs(f internalFlag, optionFuncs ...optionFunc) error {
	for _, optionFunc := range optionFuncs {
		if err := optionFunc(f); err != nil {
			return err
		}
	}

	return nil
}

func flagOptsFromOptionFuncs(optionFuncs ...optionFunc) (flagOpts, error) {
	sf := stringFlag{}

	if err := applyOptionFuncs(&sf, optionFuncs...); err != nil {
		return sf.flagOpts, err
	}

	return sf.flagOpts, nil
}
