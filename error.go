package grantedrpc

// RPCError holds details about an error occurrence
type RPCError struct {
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return e.Message
}

// ErrorPayload describes a JSON error response.
type ErrorPayload struct {
	Error RPCError `json:"error"`
}
