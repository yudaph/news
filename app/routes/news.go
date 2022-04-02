package routes

import (
	"github.com/gofiber/fiber/v2"
	"news/app/handlers"
	"news/domain/news"
)

func NewsRouter(app fiber.Router, service news.Service) {
	app.Get("/", handlers.GetAllNews(service))
	app.Get("/status/:status", handlers.GetNewsByStatus(service))
	app.Get("/topic/:topic", handlers.GetNewsByTopic(service))
	app.Get("/:slug", handlers.GetNewsBySlug(service))
	app.Post("/", handlers.AddNews(service))
	app.Patch("/:id", handlers.UpdateNews(service))
	app.Delete("/:id", handlers.DeleteNews(service))
}
