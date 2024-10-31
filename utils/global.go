package utils

func ErrServerInternal(s string) *ApiException { return NewApiException(500, s) }
func ErrNotFound(s string) *ApiException       { return NewApiException(404, s) }
func ErrValidateFailed(s string) *ApiException { return NewApiException(400, s) }
