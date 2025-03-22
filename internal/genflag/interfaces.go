package genflag

// Each flag implements this interface.
type Flag interface {
	// The name of the flag.
	Name() string
	// The value of the flag.
	Value() string
	// The rendered string of the flag as well as an error if it
	// cannot be rendered.
	String() (string, error)
	// The segmented string of the flag as well as an error if it
	// cannot be rendered.
	Segmented() ([]string, error)
}

// Any type which satisfies this interface will be marshaled.
type FlagMarshaler interface {
	MarshalFlags() ([]Flag, error)
}
