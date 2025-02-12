package videomaker

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simpleAI/service-video-maker/app/structs/model"
)

func (r *videoMakerHandler) MakeVideo(c *fiber.Ctx) error {
	ctx := c.UserContext()

	request := new(model.GenerateVideoRequest)
	errParse := c.BodyParser(request)
	if errParse != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(errParse.Error()))
	}

	videoPath, err := r.makeVideoService.GenerateVideo(ctx, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send([]byte(err.Error()))
	}
	defer r.makeVideoService.CleanUp(ctx, request.Id)

	return c.Status(fiber.StatusOK).SendFile(videoPath)
}
