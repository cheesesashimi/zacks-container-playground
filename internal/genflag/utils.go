package genflag

func stringSliceToMap(slice []string) map[string]struct{} {
	out := map[string]struct{}{}

	for _, item := range slice {
		out[item] = struct{}{}
	}

	return out
}

func stringMapToSlice(in map[string]struct{}) []string {
	out := []string{}

	for key := range in {
		out = append(out, key)
	}

	return out
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
