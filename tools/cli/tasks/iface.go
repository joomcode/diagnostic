package tasks

import (
	"context"

	"github.com/joomcode/diagnostic/tools/cli/logger"
)

type Task interface {
	ID() string
}

type TaskFn func(ctx context.Context, log logger.Logger) error

type genericTask struct {
	id string
	fn TaskFn
}

func NewGenericTask(id string, fn TaskFn) Task {
	return &genericTask{
		id: id,
		fn: fn,
	}
}

func (g *genericTask) ID() string {
	return g.id
}
