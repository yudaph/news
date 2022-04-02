package main

import (
	"fmt"
	"log"
	"news/app"
	"news/configs"
	"news/domain/news"
	"news/domain/tag"
	"news/infras"
	"news/shared/logger"
)

func main() {
	logger.InitLogger()
	configuration := configs.Get()
	cache := infras.RedisNewClient(configuration)
	mysql, err := infras.MysqlNewClient(configuration)
	if err != nil {
		fmt.Println(err)
	}
	newsRepo := news.NewRepository(mysql)
	tagsRepo := tag.NewRepository(mysql)
	newsCache := news.NewCacheImpl(cache, configuration.Cache.Redis.Expired.News)
	newsService := news.NewService(newsRepo, tagsRepo, newsCache)
	tagService := tag.NewService(tagsRepo)

	fmt.Println(cache, mysql)
	app := app.CreateApp(newsService, tagService)

	log.Fatal(app.Listen(":" + configuration.Server.Port))
}
