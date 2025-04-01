package httptool

type HTTP_API_ERROR struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	HTTP_API_ERROR_EXIST             = HTTP_API_ERROR{6001, "already exist"}
	HTTP_API_ERROR_UNEXIST           = HTTP_API_ERROR{6002, "unexist"}
	HTTP_API_ERROR_VALUE_EMPTY       = HTTP_API_ERROR{6005, "empty value"}
	HTTP_API_ERROR_INVALID_PARAMETER = HTTP_API_ERROR{6006, "invalid parameter"}
	HTTP_API_ERROR_FAIL              = HTTP_API_ERROR{6901, "fail"}
)
