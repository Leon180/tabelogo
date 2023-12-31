package helpers

func ParamsCombine(url string, params map[string]string) string {
	result := url + "?"
	for key, value := range params {
		result += key + "=" + value + "&"
	}
	return result
}
