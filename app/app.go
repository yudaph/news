package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"news/app/routes"
	"news/domain/news"
	"news/domain/tag"
)

const (
	v1 = "/api/v1"
)

func CreateApp(newsService news.Service, tagService tag.Service) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to app!"))
	})
	routes.NewsRouter(app.Group(v1+"/news"), newsService)
	routes.TagRouter(app.Group(v1+"/tag"), tagService)
	return app
}
