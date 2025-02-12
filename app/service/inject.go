package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/simpleAI/service-video-maker/app/repository"
	videomaker "github.com/simpleAI/service-video-maker/app/service/video_maker"
)

type IService interface {
	GetServiceVideoMaker() videomaker.IServiceVideoMaker
}

type service struct {
	videomaker.IServiceVideoMaker
}

func NewService(restClient *resty.Client, repository repository.IRepository) IService {
	videoMakerService := videomaker.NewServiceVideoMaker(restClient, repository)

	return &service{
		videoMakerService,
	}
}

func (s *service) GetServiceVideoMaker() videomaker.IServiceVideoMaker {
	return s.IServiceVideoMaker
}
