package processor

import (
	"context"
)

type Input interface {
	GetData() (data []byte, err error)
	GetHeader(key string) string
	GetQueryParam(key string) string
	GetParam(key string) string
}

type Output struct {
	Error      error
	Type       string
	StatusCode string
	Data       interface{}
}

// Service is a common handler so it can be change into Server or gRPC server
type Service func(ctx context.Context, input Input) (out Output)
