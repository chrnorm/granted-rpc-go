package grantedrpc

// RPCError holds details about an error occurrence
type RPCError struct {
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return e.Message
}
