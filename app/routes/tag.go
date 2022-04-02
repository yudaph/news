package routes

import (
	"github.com/gofiber/fiber/v2"
	"news/app/handlers"
	"news/domain/tag"
)

func TagRouter(app fiber.Router, service tag.Service) {
	app.Get("/", handlers.GetAllTag(service))
	//app.Get("/:like", handlers.GetAllTag(service))
	app.Post("/", handlers.AddTag(service))
	app.Put("/:id", handlers.UpdateTag(service))
	app.Delete("/:id", handlers.DeleteTag(service))
}
