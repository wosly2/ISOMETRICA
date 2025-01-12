package main

// self explanatory. can't believe go doesn't have this already
func contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

// get the unique items from a slice
func uniqueItems[T comparable](items []T) []T {
	var uniqueItems []T
	for _, item := range items {
		if !contains(uniqueItems, item) {
			uniqueItems = append(uniqueItems, item)
		}
	}

	return uniqueItems
}

// invert a map[string]int
func invertMap(m map[string]int) map[int]string {
	inverted := make(map[int]string)
	for key, value := range m {
		inverted[value] = key
	}
	return inverted
}
