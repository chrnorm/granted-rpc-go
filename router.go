package grantedrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// An RPC router designed to receive requests in JSON format for use with MacOS XPC.
//
// Example request format:
//
//	{
//	  "procuredure": "/example.EchoService/SayHello",
//	  "request": {
//	    "greeting": "hello"
//	  }
//	}
//
// Example response format:
//
//	{
//	  "procuredure": "/example.EchoService/SayHello",
//	  "response": {
//	    "reply": "hello"
//	  }
//	}
type Router struct {
	// operationFuncs is a map of each fully-qualified operation name to
	// it's corresponding handler function.
	// for example: example.EchoService.SayHello -> echoServiceHandler.SayHello()
	operationFuncs map[string]any

	// inputTypes is a map of each fully-qualified operation name
	// to it's corresponding request type.
	// for example: example.EchoService.SayHello -> SayHelloRequest
	inputTypes map[string]any

	// outputTypes is a map of each fully-qualified operation name
	// to it's corresponding request type.
	// for example: example.EchoService.SayHello -> SayHelloResponse
	outputTypes map[string]any
}

type messageHandler interface {
	// SendMessage sends the message over XPC and returns the output JSON string, or an error.
	SendMessage(ctx context.Context, input string) (string, error)
}

type routerMessage struct {
	Procedure string          `json:"procedure"`
	Request   json.RawMessage `json:"request,omitempty"`
	Response  json.RawMessage `json:"response,omitempty"`
}

func NewRouter() *Router {
	return &Router{
		operationFuncs: make(map[string]any),
		inputTypes:     make(map[string]any),
		outputTypes:    make(map[string]any),
	}
}

// Register adds operation handlers to the router
func (r *Router) Register(operation string, inputType any, outputType any, handler any) {
	r.operationFuncs[operation] = handler
	r.inputTypes[operation] = inputType
	r.outputTypes[operation] = outputType
}

// HandleMessage processes an incoming message, calling the registered handler
func (r *Router) HandleMessage(ctx context.Context, input string) (string, error) {
	var msg routerMessage
	err := json.Unmarshal([]byte(input), &msg)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Look up handler for the operation
	fn, ok := r.operationFuncs[msg.Procedure]
	if !ok {
		return "", fmt.Errorf("no handler found for operation %s", msg.Procedure)
	}

	// Get the input type
	inputType, ok := r.inputTypes[msg.Procedure]
	if !ok {
		return "", fmt.Errorf("no input type found for operation %s", msg.Procedure)
	}

	// Create a new instance of the input type
	inputValue := reflect.New(reflect.TypeOf(inputType).Elem()).Interface()

	// Unmarshal the request into the input type
	err = protojson.Unmarshal(msg.Request, inputValue.(proto.Message))
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// Call the handler function
	fnValue := reflect.ValueOf(fn)
	results := fnValue.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(inputValue),
	})

	// Check for error
	if !results[1].IsNil() {
		return "", results[1].Interface().(error)
	}

	// Get response
	response := results[0].Interface().(proto.Message)

	// Marshal response to JSON
	responseBytes, err := protojson.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	// Create response message
	outMsg := routerMessage{
		Procedure: msg.Procedure,
		Response:  responseBytes,
	}

	// Marshal full response
	out, err := json.Marshal(outMsg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response message: %w", err)
	}

	return string(out), nil
}
