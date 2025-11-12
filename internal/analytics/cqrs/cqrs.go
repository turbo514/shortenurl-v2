package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

var queryBus QueryBus

//var commandBus CommandBus

type Query any
type QueryHandler[Q Query, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}
type QueryBus struct {
	handlers map[reflect.Type]any
}

func RegisterQuery[Q Query, R any](query Q, handler QueryHandler[Q, R]) {
	t := reflect.TypeOf(query)
	queryBus.handlers[t] = handler
}

func AskQuery[Q Query, R any](ctx context.Context, query Q) (R, error) {
	t := reflect.TypeOf(query)
	h, ok := queryBus.handlers[t]
	if !ok {
		var zero R
		return zero, fmt.Errorf("no handler registered for %v", t)
	}
	return h.(QueryHandler[Q, R]).Handle(ctx, query)
}

//type Command any
//type CommandHandler interface{}
