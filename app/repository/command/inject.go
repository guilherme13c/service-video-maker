package command

import (
	"context"
)

type IComandRepository interface {
	Run(ctx context.Context, cmd string, args ...string) error
}

type commandRepository struct{}

func NewCommandRepository() IComandRepository {
	return &commandRepository{}
}
