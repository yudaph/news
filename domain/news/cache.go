package news

//go:generate go run github.com/golang/mock/mockgen -source cache.go -destination mock/cache_mock.go -package news_mock

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"news/domain/entities"
	"time"
)

type Cache interface {
	SetSliceNews(ctx context.Context, key string, sliceNewsDto *entities.SliceNewsDto) error
	SetNews(ctx context.Context, key string, newsDto *entities.NewsDto) error
	GetSliceNews(ctx context.Context, key string) (sliceNewsDto *entities.SliceNewsDto, err error)
	GetNews(ctx context.Context, key string) (newsDto *entities.NewsDto, err error)
}

type cacheImpl struct {
	redis          *redis.Client
	expiredMinutes time.Duration
}

//type newsGenerics interface {
//	entities.SliceNewsDto | entities.NewsDto
//}

//func set[newsType newsGenerics](redis *redis.Client, expired time.Duration, key string, news newsType) (err error) {
//	value, err := json.Marshal(news)
//	if err != nil {
//		return err
//	}
//
//	return redis.Set(key, value, expired).Err()
//}

//func get[newsType newsGenerics](redis *redis.Client, expired time.Duration, key string) (news newsType, err error) {
//	val, err := redis.Get(key).Result()
//	if err != nil {
//		return
//	}
//	err = json.Unmarshal([]byte(val), &news)
//	return news, err
//}

func NewCacheImpl(redis *redis.Client, expiredMinutes int) *cacheImpl {
	return &cacheImpl{redis: redis, expiredMinutes: time.Duration(expiredMinutes) * time.Second}
}

func (c *cacheImpl) SetSliceNews(ctx context.Context, key string, sliceNewsDto *entities.SliceNewsDto) error {
	value, err := json.Marshal(sliceNewsDto)
	if err != nil {
		return err
	}

	return c.redis.Set(key, value, c.expiredMinutes).Err()
}

func (c *cacheImpl) SetNews(ctx context.Context, key string, newsDto *entities.NewsDto) error {
	value, err := json.Marshal(newsDto)
	if err != nil {
		return err
	}

	return c.redis.Set(key, value, c.expiredMinutes).Err()
}

func (c *cacheImpl) GetSliceNews(ctx context.Context, key string) (sliceNewsDto *entities.SliceNewsDto, err error) {
	val, err := c.redis.Get(key).Result()
	if err != nil {
		return
	}
	sliceNewsDto = &entities.SliceNewsDto{}
	err = json.Unmarshal([]byte(val), sliceNewsDto)
	return
}

func (c *cacheImpl) GetNews(ctx context.Context, key string) (newsDto *entities.NewsDto, err error) {
	val, err := c.redis.Get(key).Result()
	if err != nil {
		return
	}
	newsDto = &entities.NewsDto{}
	err = json.Unmarshal([]byte(val), newsDto)
	return
}
