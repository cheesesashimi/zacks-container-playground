package genflag

func flagsToStrings(flags []Flag) ([]string, error) {
	out := []string{}
	for _, flag := range flags {
		s, err := flag.String()
		if err != nil {
			return nil, err
		}

		out = append(out, s)
	}

	return out, nil
}

func boolToPtr(val bool) *bool {
	return &val
}

func stringToPtr(s string) *string {
	return &s
}

// Variadic helper function for running a series of functions with return
// flags. This will execute each function, halt on any errors, and combine all
// of the flags into a singular array.
func combineFlags(flagFuncs ...func() ([]Flag, error)) ([]Flag, error) {
	out := []Flag{}

	for _, getFlagFunc := range flagFuncs {
		f, err := getFlagFunc()
		if err != nil {
			return nil, err
		}

		out = append(out, f...)
	}

	return out, nil
}
