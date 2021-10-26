package main

func Merge(receiver *[]string, giver []string) {
	for _, str := range giver {
		*receiver = append(*receiver, str)
	}
}

func containsInt(rack []int, item int) bool {
	for _, elem := range rack {
		if elem == item {
			return true
		}
	}
	return false
}

func containsStr(rack []string, item string) bool {
	for _, elem := range rack {
		if elem == item {
			return true
		}
	}
	return false
}
