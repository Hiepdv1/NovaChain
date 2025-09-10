package err

import "strings"

const (
	CodeInternal = -32603

	CodeNotFound        = -32001
	CodeInvalidArgument = -32002
	CodeDatabase        = -32010
)

type RPCError struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func ErrNotFound(msg ...string) *RPCError {
	if msg != nil {
		return &RPCError{Code: CodeNotFound, Message: strings.Join(msg, " ")}
	}
	return &RPCError{Code: CodeNotFound, Message: "Resource not found"}
}

func ErrInvalidArgument(msg ...string) *RPCError {
	if msg != nil {
		return &RPCError{Code: CodeInvalidArgument, Message: strings.Join(msg, " ")}
	}
	return &RPCError{Code: CodeInvalidArgument, Message: strings.Join(msg, " ")}
}

func ErrDatabase(msg ...string) *RPCError {
	if msg != nil {
		return &RPCError{Code: CodeDatabase, Message: strings.Join(msg, " ")}
	}
	return &RPCError{Code: CodeDatabase, Message: "Database error"}
}

func ErrInternal(msg ...string) *RPCError {
	if msg != nil {
		return &RPCError{Code: CodeInternal, Message: strings.Join(msg, " ")}
	}
	return &RPCError{Code: CodeInternal, Message: "Internal error"}
}
