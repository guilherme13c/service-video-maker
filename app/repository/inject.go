package repository

import (
	commandRepository "github.com/simpleAI/service-video-maker/app/repository/command"
)

type IRepository interface {
	GetCommandRepository() commandRepository.IComandRepository
}

type repository struct {
	commandRepository commandRepository.IComandRepository
}

func NewRepository() IRepository {
	commandRepository := commandRepository.NewCommandRepository()

	return &repository{
		commandRepository: commandRepository,
	}
}

func (r *repository) GetCommandRepository() commandRepository.IComandRepository {
	return r.commandRepository
}
