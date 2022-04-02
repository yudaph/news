package handlers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"news/domain/entities"
	"news/domain/tag"
	"news/shared/failure"
)

func AddTag(service tag.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.CreateTag
		err := c.BodyParser(&requestBody)
		if err != nil {
			return ErrorResponse(c, failure.BadRequestWithString("bad request"))
		}

		err = requestBody.Validate()
		if err != nil {
			return ErrorResponse(c, err)
		}

		result, err := service.Create(c.Context(), &requestBody)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusCreated, result)
	}
}

func GetAllTag(service tag.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetAll(c.Context())
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func UpdateTag(service tag.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.TagDto
		err := c.BodyParser(&requestBody)
		if err != nil {
			return ErrorResponse(c, failure.BadRequestWithString("bad request"))
		}

		err = requestBody.Validate()
		if err != nil {
			return ErrorResponse(c, err)
		}

		requestBody.ID = c.Params("id")
		result, err := service.Update(c.Context(), &requestBody)
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, result)
	}
}

func DeleteTag(service tag.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := service.Delete(c.Context(), c.Params("id"))
		if err != nil {
			return ErrorResponse(c, err)
		}
		return SuccessResponse(c, http.StatusOK, &fiber.Map{
			"message": "success",
		})
	}
}
