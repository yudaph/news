package news

//go:generate go run github.com/golang/mock/mockgen -source repository.go -destination mock/repository_mock.go -package news_mock

import (
	"context"
	"github.com/jmoiron/sqlx"
	"news/domain/entities"
	"news/shared/failure"
	"news/shared/logger"
)

type Repository interface {
	CreateNews(ctx context.Context, news *entities.News) error
	GetNewsBySlug(ctx context.Context, slug string) (*entities.News, error)
	GetNewsByTopic(ctx context.Context, topic string) (*entities.SliceNews, error)
	GetNewsByStatus(ctx context.Context, status entities.NewsStatus) (*entities.SliceNews, error)
	GetAllNews(ctx context.Context) (*entities.SliceNews, error)
	UpdateNews(ctx context.Context, news *entities.News) error
	DeleteNews(ctx context.Context, id string) error
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) *repository {
	return &repository{DB: DB}
}

func (r *repository) CreateNews(ctx context.Context, news *entities.News) (err error) {
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	err = r.insertNews(tx, news)
	if err != nil {
		tx.Rollback()
		return
	}
	err = r.insertNewsTags(tx, news)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (r *repository) insertNews(tx *sqlx.Tx, news *entities.News) (err error) {
	query := "INSERT INTO `news`(`id`, `title`, `slug`, `content`, `topic`, `status`, `createdAt`) " +
		"VALUES (:id, :title, :slug, :content, :topic, :status, :createdAt)"
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	_, err = stmt.Exec(news)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return
}

func (r *repository) insertNewsTags(tx *sqlx.Tx, news *entities.News) (err error) {
	query := "INSERT INTO `news_tags`(`news_id`, `tag_id`) VALUES(:news_id, :tag_id)"
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	sliceNewsTag := news.ToSliceNewsTag()
	for _, newsTag := range sliceNewsTag {
		_, err = stmt.Exec(newsTag)
		if err != nil {
			logger.ErrorWithStack(err)
			err = failure.InternalServerError
			break
		}
	}
	return
}

func (r *repository) GetNewsBySlug(ctx context.Context, slug string) (news *entities.News, err error) {
	sliceNews, err := r.selectNews(ctx, "WHERE slug = ? AND status = ?", slug, entities.NewsPublish)
	if err != nil {
		return
	}
	news = &(*sliceNews)[0]
	tags, err := r.selectNewsTag(ctx, "WHERE news_id = ?", news.ID)
	if err != nil {
		return
	}
	news.Tags = tags.ToMapTag()[news.ID]
	return
}

func (r *repository) GetNewsByTopic(ctx context.Context, topic string) (sliceNews *entities.SliceNews, err error) {
	sliceNews, err = r.selectNews(ctx, "WHERE topic = ? and status = ?", topic, entities.NewsPublish)
	if err != nil {
		return
	}
	tags, err := r.selectNewsTagByNewsIds(ctx, extractNewsId(sliceNews))
	if err != nil {
		return
	}
	compositeNewsTags(sliceNews, tags)
	return
}

func (r *repository) GetNewsByStatus(ctx context.Context, status entities.NewsStatus) (sliceNews *entities.SliceNews, err error) {
	sliceNews, err = r.selectNews(ctx, "WHERE status = ?", status)
	if err != nil {
		return
	}
	tags, err := r.selectNewsTagByNewsIds(ctx, extractNewsId(sliceNews))
	if err != nil {
		return
	}
	compositeNewsTags(sliceNews, tags)
	return
}

func (r *repository) GetAllNews(ctx context.Context) (sliceNews *entities.SliceNews, err error) {
	sliceNews, err = r.selectNews(ctx, "WHERE status <> ?", entities.NewsDeleted)
	if err != nil {
		return
	}
	tags, err := r.selectNewsTagByNewsIds(ctx, extractNewsId(sliceNews))
	if err != nil {
		return
	}
	compositeNewsTags(sliceNews, tags)
	return
}

func (r *repository) selectNews(ctx context.Context, where string, args ...interface{}) (news *entities.SliceNews, err error) {
	news = new(entities.SliceNews)
	query := "SELECT id, title, slug, content, topic, status, createdAt FROM `news` " + where + " ORDER BY createdAt desc"
	err = r.DB.SelectContext(ctx, news, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	if len(*news) < 1 {
		err = failure.NotFound("news not found")
		return
	}
	return
}

func (r *repository) selectNewsTagByNewsIds(ctx context.Context, ids []string) (tags *entities.SliceNewsTag, err error) {
	where, args, err := sqlx.In("WHERE `news_id` IN (?)", ids)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return r.selectNewsTag(ctx, where, args...)
}

func (r *repository) selectNewsTag(ctx context.Context, where string, args ...interface{}) (tags *entities.SliceNewsTag, err error) {
	tags = new(entities.SliceNewsTag)
	query := "SELECT `news_id`, `tag_id` FROM `news_tags` " + where
	err = r.DB.SelectContext(ctx, tags, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	if len(*tags) < 1 {
		err = failure.NotFound("tags not found")
	}
	return
}

func (r *repository) UpdateNews(ctx context.Context, news *entities.News) (err error) {
	sliceNews, err := r.selectNews(ctx, "WHERE id = ?", news.ID)
	if err != nil {
		return
	}
	newNews := (*sliceNews)[0]
	newNews.Update(*news)

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		logger.ErrorWithStack(err)
		return failure.InternalServerError
	}

	err = r.updateNews(tx, &newNews)
	if err != nil {
		tx.Rollback()
		return
	}

	err = r.deleteNewsTag(tx, newNews.ID)
	if err != nil {
		tx.Rollback()
		return
	}

	err = r.insertNewsTags(tx, &newNews)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (r *repository) DeleteNews(ctx context.Context, id string) (err error) {
	sliceNews, err := r.selectNews(ctx, "WHERE id = ?", id)
	if err != nil {
		return
	}
	newNews := (*sliceNews)[0]
	newNews.Delete()
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	err = r.updateNews(tx, &newNews)
	if err != nil {
		return
	}
	tx.Commit()
	return
}

func (r *repository) updateNews(tx *sqlx.Tx, news *entities.News) (err error) {
	query := "UPDATE `news` SET title = :title, content = :content, topic = :topic, status = :status, " +
		"deletedAt = :deletedAt WHERE id = :id"
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	_, err = stmt.Exec(news)
	if err != nil {
		logger.ErrorWithStack(err)
		err = failure.InternalServerError
		return
	}
	return
}

func (r *repository) deleteNewsTag(tx *sqlx.Tx, id string) error {
	query := "DELETE FROM `news_tags` WHERE `news_id` = ?"
	_, err := tx.Exec(query, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return failure.InternalServerError
	}
	return nil
}

func extractNewsId(sliceNews *entities.SliceNews) (res []string) {
	for _, news := range *sliceNews {
		res = append(res, news.ID)
	}
	return
}

func compositeNewsTags(sliceNews *entities.SliceNews, tags *entities.SliceNewsTag) {
	mapTags := tags.ToMapTag()
	for i, news := range *sliceNews {
		(*sliceNews)[i].Tags = mapTags[news.ID]
	}
}
