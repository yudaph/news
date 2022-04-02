package tag

//go:generate go run github.com/golang/mock/mockgen -source repository.go -destination mock/repository_mock.go -package tag_mock
import (
	"context"
	"github.com/jmoiron/sqlx"
	"news/domain/entities"
	"news/shared/failure"
	"news/shared/logger"
)

type Repository interface {
	CreateTag(ctx context.Context, tag *entities.Tag) (result *entities.Tag, err error)
	GetAllTag(ctx context.Context) (result *entities.Tags, err error)
	GetTagLike(ctx context.Context, like string) (result *entities.Tags, err error)
	GetTagByIds(ctx context.Context, id []string) (result *entities.Tags, err error)
	UpdateTag(ctx context.Context, tag *entities.Tag) (result *entities.Tag, err error)
	DeleteTag(ctx context.Context, id string) (err error)
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) *repository {
	return &repository{DB: DB}
}

func (r *repository) CreateTag(ctx context.Context, tag *entities.Tag) (result *entities.Tag, err error) {
	query := "INSERT INTO `tags`(`id`, `name`, `status`) VALUES (:id, :name, :status)"
	stmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	_, err = stmt.Exec(tag)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return tag, nil
}

func (r *repository) GetAllTag(ctx context.Context) (result *entities.Tags, err error) {
	return r.selectTag(ctx, "")
}

func (r *repository) GetTagLike(ctx context.Context, like string) (result *entities.Tags, err error) {
	return r.selectTag(ctx, "WHERE name LIKE '%?%' and status = ?", like, entities.TagActive)
}

func (r *repository) GetTagByIds(ctx context.Context, id []string) (result *entities.Tags, err error) {
	query, args, err := sqlx.In("WHERE id IN (?) and status = ?", id, entities.TagActive)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return r.selectTag(ctx, query, args...)
}

func (r *repository) UpdateTag(ctx context.Context, tag *entities.Tag) (result *entities.Tag, err error) {
	oldTag, err := r.selectTag(ctx, "WHERE id = ?", tag.ID)
	if err != nil {
		return
	}
	result = &(*oldTag)[0]
	result.UpdateTag(tag)
	err = r.updateTag(ctx, result)
	return
}

func (r *repository) DeleteTag(ctx context.Context, id string) (err error) {
	oldTag, err := r.selectTag(ctx, "WHERE id = ?", id)
	if err != nil {
		return
	}
	deletedTag := &(*oldTag)[0]
	deletedTag.Delete()
	err = r.updateTag(ctx, deletedTag)
	return
}

func (r *repository) selectTag(ctx context.Context, where string, args ...interface{}) (tags *entities.Tags, err error) {
	tags = new(entities.Tags)
	query := "SELECT `id`, `name`, `status` FROM `tags` " + where
	err = r.DB.SelectContext(ctx, tags, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	if len(*tags) < 1 {
		err = failure.NotFound("tag not found")
	}
	return
}

func (r *repository) updateTag(ctx context.Context, tag *entities.Tag) (err error) {
	query := "UPDATE `tags` SET name = :name, status = :status WHERE id = :id"
	stmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	_, err = stmt.Exec(tag)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return
}
