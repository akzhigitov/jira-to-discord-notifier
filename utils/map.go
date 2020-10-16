package utils

import "strings"

func String2Map(raw string, itemSep string, kvSep string) map[string]string {
	result := map[string]string{}

	for _, item := range strings.Split(raw, itemSep) {
		key, value := parseKeyValue(item, kvSep)
		result[key] = value
	}

	return result
}

func parseKeyValue(item string, sep string) (key string, value string) {
	result := strings.Split(item, sep)
	return result[0], result[1]
}
