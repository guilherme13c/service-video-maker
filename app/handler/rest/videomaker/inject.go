package videomaker

import (
	"github.com/gofiber/fiber/v2"
	serviceVideoMaker "github.com/simpleAI/service-video-maker/app/service/video_maker"
)

type IVideoMakerHandler interface {
	MakeVideo(c *fiber.Ctx) error
}

type videoMakerHandler struct {
	makeVideoService serviceVideoMaker.IServiceVideoMaker
}

func NewVideoMakerHandler(makeVideoService serviceVideoMaker.IServiceVideoMaker) IVideoMakerHandler {
	return &videoMakerHandler{
		makeVideoService: makeVideoService,
	}
}
