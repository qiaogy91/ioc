package utils

func NewApiException(e int, m string) *ApiException {
	return &ApiException{Code: e, Message: m}
}

type ApiException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiException) Error() string {
	return e.Message
}
