package genflag

func flagsToString(flags []Flag) ([]string, error) {
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

func flagsToMap(flags []Flag) (map[string]struct{}, error) {
	s, err := flagsToString(flags)
	if err != nil {
		return nil, err
	}

	return stringSliceToMap(s), nil
}

func isStringSliceValuesUnique(slice []string) (bool, []string) {
	counts := map[string]int{}

	for _, item := range slice {
		_, ok := counts[item]
		if ok {
			counts[item]++
		} else {
			counts[item] = 1
		}
	}

	multiples := []string{}

	multiplesFound := false

	for key, count := range counts {
		if count > 1 {
			multiples = append(multiples, key)
		}
	}

	return multiplesFound, multiples
}

func stringSliceToMap(slice []string) map[string]struct{} {
	out := map[string]struct{}{}

	for _, item := range slice {
		out[item] = struct{}{}
	}

	return out
}

func toPlural(f Flag, err error) ([]Flag, error) {
	return []Flag{f}, err
}
