package news

import (
	"context"
	"fmt"
	"news/domain/entities"
	"news/domain/tag"
	"news/shared/logger"
)

type Service interface {
	Create(ctx context.Context, dto *entities.NewsDto) (result *entities.NewsDto, err error)
	GetAll(ctx context.Context) (result *entities.SliceNewsDto, err error)
	GetBySlug(ctx context.Context, slug string) (result *entities.NewsDto, err error)
	GetByTopic(ctx context.Context, topic string) (result *entities.SliceNewsDto, err error)
	GetByStatus(ctx context.Context, status entities.NewsStatus) (result *entities.SliceNewsDto, err error)
	Update(ctx context.Context, dto *entities.NewsDto) (err error)
	Delete(ctx context.Context, id string) (err error)
}

type serviceImpl struct {
	repo    Repository
	tagRepo tag.Repository
	cache   Cache
}

func NewService(repo Repository, tagRepo tag.Repository, cache Cache) *serviceImpl {
	return &serviceImpl{repo: repo, tagRepo: tagRepo, cache: cache}
}

func (s *serviceImpl) Create(ctx context.Context, dto *entities.NewsDto) (result *entities.NewsDto, err error) {
	news, err := dto.ToNews()
	if err != nil {
		return
	}

	err = s.repo.CreateNews(ctx, news)
	if err != nil {
		return
	}

	result = news.ToNewsDto()
	return
}

func (s *serviceImpl) GetAll(ctx context.Context) (result *entities.SliceNewsDto, err error) {
	result, err = s.cache.GetSliceNews(ctx, "all:")
	if err == nil {
		fmt.Println("Get from Cache")
		return
	}
	fmt.Println("Search from DB")
	sliceNews, err := s.repo.GetAllNews(ctx)
	if err != nil {
		return
	}

	tags, err := s.tagRepo.GetTagByIds(ctx, sliceNews.GetSliceTagIds())
	if err != nil {
		return
	}

	result = sliceNews.ToSliceNewsDto(tags.ToMapTags)

	errs := s.cache.SetSliceNews(ctx, "all:", result)
	if errs != nil {
		logger.ErrorWithStack(errs)
	}
	return
}

func (s *serviceImpl) GetByTopic(ctx context.Context, topic string) (result *entities.SliceNewsDto, err error) {
	result, err = s.cache.GetSliceNews(ctx, "topic:"+topic)
	if err == nil {
		fmt.Println("Get from Cache")
		return
	}
	fmt.Println("Search from DB")

	sliceNews, err := s.repo.GetNewsByTopic(ctx, topic)
	if err != nil {
		return
	}

	tags, err := s.tagRepo.GetTagByIds(ctx, sliceNews.GetSliceTagIds())
	if err != nil {
		return
	}

	result = sliceNews.ToSliceNewsDto(tags.ToMapTags)

	errs := s.cache.SetSliceNews(ctx, "topic:"+topic, result)
	if errs != nil {
		logger.ErrorWithStack(errs)
	}
	return
}

func (s *serviceImpl) GetByStatus(ctx context.Context, status entities.NewsStatus) (result *entities.SliceNewsDto, err error) {
	result, err = s.cache.GetSliceNews(ctx, "status:"+status.String())
	if err == nil {
		fmt.Println("Get from Cache")
		return
	}
	fmt.Println("Search from DB")

	sliceNews, err := s.repo.GetNewsByStatus(ctx, status)
	if err != nil {
		return
	}

	tags, err := s.tagRepo.GetTagByIds(ctx, sliceNews.GetSliceTagIds())
	if err != nil {
		return
	}

	result = sliceNews.ToSliceNewsDto(tags.ToMapTags)

	errs := s.cache.SetSliceNews(ctx, "status:"+status.String(), result)
	if errs != nil {
		logger.ErrorWithStack(errs)
	}
	return
}

func (s *serviceImpl) GetBySlug(ctx context.Context, slug string) (result *entities.NewsDto, err error) {
	result, err = s.cache.GetNews(ctx, "slug:"+slug)
	if err == nil {
		fmt.Println("Get from Cache")
		return
	}
	fmt.Println("Search from DB")

	news, err := s.repo.GetNewsBySlug(ctx, slug)
	if err != nil {
		return
	}

	tags, err := s.tagRepo.GetTagByIds(ctx, news.Tags)
	if err != nil {
		return
	}

	result = news.ToNewsDto(tags.ToMapTags())

	errs := s.cache.SetNews(ctx, "slug:"+slug, result)
	if errs != nil {
		logger.ErrorWithStack(errs)
	}
	return
}

func (s *serviceImpl) Update(ctx context.Context, dto *entities.NewsDto) (err error) {
	news, err := dto.ToNews()
	if err != nil {
		return
	}

	return s.repo.UpdateNews(ctx, news)
}

func (s *serviceImpl) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteNews(ctx, id)
}
