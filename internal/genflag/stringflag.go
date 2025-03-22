package genflag

var _ Flag = stringFlag{}
var _ Flag = &stringFlag{}

// Defines a string flag field.
type stringFlag struct {
	// The name of the flag.
	name string
	// The value of the flag. However, the value cannot be empty.
	value string
	// The base flag options struct.
	flagOpts
}

// Returns a new validated string flag.
func NewStringFlag(name, value string, optionFuncs ...optionFunc) (Flag, error) {
	f, err := newStringFlag(name, value, optionFuncs...)

	if err != nil {
		return nil, err
	}

	return f, nil
}

// Returns a new validated string flag as a concrete type. Mostly
// for internal use.
func newStringFlag(name, value string, optionFuncs ...optionFunc) (stringFlag, error) {
	sf := stringFlag{
		name:  name,
		value: value,
	}

	if err := applyOptionFuncs(&sf, optionFuncs...); err != nil {
		return sf, err
	}

	if err := sf.flagOpts.validate(value); err != nil {
		return sf, err
	}

	return sf, nil
}

// Returns the name of the flag.
func (s stringFlag) Name() string {
	return s.name
}

// Returns the value of the flag.
func (s stringFlag) Value() string {
	return s.value
}

// Returns a rendered string of the flag as well as any rendering
// or validation errors.
func (s stringFlag) String() (string, error) {
	return s.flagOpts.render(s.name, s.value)
}

// Returns a segmented representation of the flag as well as any
// rendering or validation errors.
func (s stringFlag) Segmented() ([]string, error) {
	return s.flagOpts.renderSegmented(s.name, s.value)
}
