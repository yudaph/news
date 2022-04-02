package tag

import (
	"context"
	"news/domain/entities"
)

type Service interface {
	Create(ctx context.Context, dto *entities.CreateTag) (*entities.TagDto, error)
	Update(ctx context.Context, dto *entities.TagDto) (*entities.TagDto, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) (*[]entities.TagDto, error)
	Search(ctx context.Context, name string) (*[]entities.TagDto, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s service) Create(ctx context.Context, dto *entities.CreateTag) (result *entities.TagDto, err error) {
	tag, err := s.repo.CreateTag(ctx, dto.ToTag())
	if err != nil {
		return
	}
	result = tag.ToDto()
	return
}

func (s service) Update(ctx context.Context, dto *entities.TagDto) (result *entities.TagDto, err error) {
	tag, err := s.repo.UpdateTag(ctx, dto.ToTag())
	result = tag.ToDto()
	return
}

func (s service) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteTag(ctx, id)
}

func (s service) GetAll(ctx context.Context) (result *[]entities.TagDto, err error) {
	tags, err := s.repo.GetAllTag(ctx)
	if err != nil {
		return
	}
	result = tags.ToTagsDto()
	return
}

func (s service) Search(ctx context.Context, name string) (result *[]entities.TagDto, err error) {
	tags, err := s.repo.GetTagLike(ctx, name)
	if err != nil {
		return
	}
	result = tags.ToTagsDto()
	return
}
