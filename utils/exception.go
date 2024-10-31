package utils

func NewApiException(c int, m string, err error) *ApiException {
	return &ApiException{
		Code:    c,
		Message: m,
		Cause:   err.Error(),
	}
}

type ApiException struct {
	Code     int    `json:"code"`
	HttpCode int    `json:"httpCode"`
	Message  string `json:"message"`
	Cause    string `json:"cause"`
}

func (e *ApiException) Error() string {
	return e.Message
}

func (e *ApiException) WithHttpCode(c int) *ApiException {
	e.HttpCode = c
	return e
}
