package grantedrpc

import "context"

type Transport interface {
	// SendMessage sends a message over the underlying XPC connection.
	SendMessage(ctx context.Context, input string) (string, error)
}
