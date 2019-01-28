package stringarray

func Filter(array []string, fn func(string) bool) []string {
	filtered := make([]string, 0)
	for _, v := range array {
		if fn(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Map(vs [][]string, f func([]string) []string) [][]string {
	vsm := make([][]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
