package utils

func IsStatusCodeOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
