package util

func GetMapValueOrDefault(dict map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if x, found := dict[key]; found {
		return x
	}

	return defaultValue
}
