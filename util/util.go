package util

func ConvertToStrings[T any](objects []T, toString func(T) string) []string {
	var result []string
	for _, obj := range objects {
		result = append(result, toString(obj))
	}
	return result
}

func ContainsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
