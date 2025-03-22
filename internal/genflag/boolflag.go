package genflag

import (
	"fmt"
)

var _ Flag = boolFlag{}
var _ Flag = &boolFlag{}

// Defines a boolean flag field. This flag has a few operational modes:
//
// Implicit mode:
// In this mode, the flag will only be rendered if the value is true. If the
// value is false, an error will be returned instead of the flag being rendered.
//
// Explicit mode:
// In this mode, the flag will be followed by a string indicating whether it is
// true or false. No error will be returned when in this mode.
type boolFlag struct {
	// The name of the flag.
	name string
	// The value of the flag.
	value bool
	// Whether the flag should be considered "explicit" (see above).
	explicit bool
	// Any transformation functions to perform upon converting the value to
	// a string such as making it all uppercoase, titlecase, etc.
	transform func(string) string
	// The base flag options struct.
	flagOpts
}

// Returns a new validated boolean flag.
func NewBoolFlag(name string, value bool, optionFuncs ...optionFunc) (Flag, error) {
	f, err := newBoolFlag(name, value, optionFuncs...)

	if err != nil {
		return nil, err
	}

	return f, nil
}

// Returns a new validated boolean flag as a concrete type.
// Mostly for internal use.
func newBoolFlag(name string, value bool, optionFuncs ...optionFunc) (boolFlag, error) {
	bf := boolFlag{
		name:     name,
		value:    value,
		flagOpts: flagOpts{},
	}

	if err := applyOptionFuncs(&bf, optionFuncs...); err != nil {
		return bf, err
	}

	if err := bf.validate(); err != nil {
		return bf, err
	}

	return bf, nil
}

// Returns the name of the flag.
func (b boolFlag) Name() string {
	return b.name
}

// Returns a string representatino of the flag value.
// Note: This representation is without any of the transformers
// being applied.
func (b boolFlag) Value() string {
	return fmt.Sprintf("%v", b.value)
}

// Returns a rendered string of the flag name and value.
func (b boolFlag) String() (string, error) {
	if err := b.validate(); err != nil {
		return "", err
	}

	return b.render(b.name, b.getStringValue())
}

// Returns a segemented string slice of the flag name and value.
func (b boolFlag) Segmented() ([]string, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	return b.renderSegmented(b.name, b.getStringValue())
}

// Validates whether the flag can be rendered.
func (b boolFlag) validate() error {
	if !b.explicit && !b.value {
		return fmt.Errorf("cannot set implicit false bool flag")
	}

	return b.flagOpts.validate(b.getStringValue())
}

// Gets the explicit string value, applying any transformers such
// as titlecase or upppercase before returning it.
func (b boolFlag) getExplicitStringValue() string {
	val := b.Value()

	if b.transform == nil {
		return val
	}

	return b.transform(val)
}

// Gets the string value if explicit or empty if implicit.
func (b boolFlag) getStringValue() string {
	if b.explicit {
		return b.getExplicitStringValue()
	}

	return ""
}
