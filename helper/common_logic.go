package helper

func CountUniqueElements(slice []string) int {
	uniqueMap := make(map[string]bool)

	// Iterate over the slice and add each element to the map
	for _, element := range slice {
		uniqueMap[element] = true
	}

	// Return the number of unique elements (length of the map)
	return len(uniqueMap)
}
