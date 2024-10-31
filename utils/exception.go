package utils

func NewApiException(e int, m string) *ApiException {
	return &ApiException{Code: e, Message: m}
}

type ApiException struct {
	Code     int    `json:"code"`
	HttpCode int    `json:"httpCode"`
	Message  string `json:"message"`
}

func (e *ApiException) Error() string {
	return e.Message
}

func (e *ApiException) WithHttpCode(c int) *ApiException {
	e.HttpCode = c
	return e
}
