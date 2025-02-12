package videomaker

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/simpleAI/service-video-maker/app/repository"
	"github.com/simpleAI/service-video-maker/app/structs/model"
)

type IServiceVideoMaker interface {
	GenerateVideo(ctx context.Context, request *model.GenerateVideoRequest) (string, error)
	CleanUp(ctx context.Context, requestId string) error
}

type serviceVideoMaker struct {
	restClient *resty.Client
	repository repository.IRepository
}

func NewServiceVideoMaker(restClient *resty.Client, repository repository.IRepository) IServiceVideoMaker {
	return &serviceVideoMaker{
		restClient: restClient,
		repository: repository,
	}
}
