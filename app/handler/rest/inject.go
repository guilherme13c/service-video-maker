package resthandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simpleAI/service-video-maker/app/handler/rest/videomaker"
	"github.com/simpleAI/service-video-maker/app/service"
)

type IRestHandler interface {
	MakeVideo(c *fiber.Ctx) error
}

type restHandler struct {
	videomaker.IVideoMakerHandler
}

func NewRestHandler(makeVideoService service.IService) IRestHandler {
	videoMakerHandler := videomaker.NewVideoMakerHandler(makeVideoService.GetServiceVideoMaker())

	return &restHandler{
		videoMakerHandler,
	}
}
