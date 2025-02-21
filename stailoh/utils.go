package stailoh

func contains[T comparable](a []T, x T) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
