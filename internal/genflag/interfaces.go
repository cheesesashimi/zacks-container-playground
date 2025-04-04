package genflag

type FlagContainer interface {
	SetName(string)
	SetOptionFuncs([]optionFunc)
	Flags() ([]Flag, error)
}

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
type Marshaler interface {
	MarshalFlags() ([]Flag, error)
}

// A simple stringer interface reference.
type stringer interface {
	String() string
}
