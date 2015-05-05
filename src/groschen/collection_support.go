package groschen

import ()

type Predicate func(value string) bool

func AddListToMap(inputList []string, destinationSet map[string]bool) {
	for _, v := range inputList {
		destinationSet[v] = true
	}
}

func AddMapToMap(sourceSet map[string]bool, destinationSet map[string]bool) {
	for v := range sourceSet {
		destinationSet[v] = true
	}
}

func NewFromFilter(source map[string]bool, test Predicate) map[string]bool {
	var result = make(map[string]bool, 0)
	for value := range source {
		if test(value) {
			result[value] = true
		}
	}
	return result
}

func SliceFromMapKeys(source map[string]bool) []string {
	result := make([]string, len(source))

	i := 0
	for k := range source {
		result[i] = k
		i += 1
	}
	return result
}
