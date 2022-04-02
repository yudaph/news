package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"news/domain/entities"
	"news/domain/news"
	"news/shared/failure"
)

func AddNews(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.NewsDto
		err := c.BodyParser(&requestBody)
		if err != nil {
			log.Trace().Err(err)
			return ErrorResponse(c, err)
		}
		err = requestBody.Validate()
		if err != nil {
			log.Trace().Err(err)
			return ErrorResponse(c, err)
		}
		result, err := service.Create(c.Context(), &requestBody)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusCreated, result)
	}
}

func GetNewsBySlug(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetBySlug(c.Context(), c.Params("slug"))
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func GetNewsByTopic(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		topic, err := url.QueryUnescape(c.Params("topic"))
		if err != nil {
			return ErrorResponse(c, failure.BadRequestWithString("bad request"))
		}
		result, err := service.GetByTopic(c.Context(), topic)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func GetNewsByStatus(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status, err := entities.StringToNewsStatus(c.Params("status"))
		if err != nil {
			return ErrorResponse(c, failure.BadRequestWithString("bad request"))
		}
		result, err := service.GetByStatus(c.Context(), status)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func GetAllNews(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetAll(c.Context())
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func UpdateNews(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.NewsDto
		err := c.BodyParser(&requestBody)
		if err != nil {
			log.Trace().Err(err)
			return ErrorResponse(c, err)
		}
		err = requestBody.Validate()
		if err != nil {
			log.Trace().Err(err)
			return ErrorResponse(c, err)
		}
		requestBody.ID = c.Params("id")
		err = service.Update(c.Context(), &requestBody)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, &fiber.Map{"message": "updated"})
	}
}

func DeleteNews(service news.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := service.Delete(c.Context(), c.Params("id"))
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, &fiber.Map{"message": "news has been deleted"})
	}
}
