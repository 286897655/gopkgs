package utils

type Error struct {
	Code        int
	Description string
}

func (err *Error) Error() string {
	return err.Description
}

var (
	ERR_INVALID_PARAMETER = &Error{Code: 6001, Description: "invalid parameter"} // 6001-invalid parameter
	ERR_UNIMPLEMENTATION  = &Error{Code: 7001, Description: "unimplementation"}  // 7001-unimplementation
	ERR_UNKNOWN           = &Error{Code: 9999, Description: "unknown"}           // 9999-unknown
)
