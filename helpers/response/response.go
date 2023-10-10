package helpers

func FailedResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "failed",
		"message": message,
	}
}

func SuccessResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
	}
}

func SuccessWithDataResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}
