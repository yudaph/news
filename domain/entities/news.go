package entities

import (
	"news/shared/Date"
	"news/shared/IDGEN"
	"news/shared/Slug"
	"news/shared/failure"
	"time"
)

type NewsStatus int

const (
	NewsDraft NewsStatus = iota + 1
	NewsPublish
	NewsDeleted
)

func StringToNewsStatus(status string) (NewsStatus, error) {
	mapStatus := map[string]NewsStatus{
		"draft":   NewsDraft,
		"publish": NewsPublish,
		"deleted": NewsDeleted,
	}
	if v, ok := mapStatus[status]; ok {
		return v, nil
	}
	return 0, failure.NotFound("status not found")
}

func (n NewsStatus) String() string {
	stringer := []string{"not found", "draft", "publish", "deleted"}
	return stringer[n]
}

type News struct {
	ID        string     `db:"id"`
	Title     string     `db:"title"`
	Slug      string     `db:"slug"`
	Content   string     `db:"content"`
	Status    NewsStatus `db:"status"`
	Tags      []string
	Topic     string     `db:"topic"`
	CreatedAt time.Time  `db:"createdAt"`
	DeletedAt *time.Time `db:"deletedAt"`
}

func NewNews(id string, title string, slug string, content string, status NewsStatus, tags []string, topic string) *News {
	if slug == "" {
		slug = Slug.Create(title)
	}
	if id == "" {
		id = IDGEN.NewUUID()
	}
	return &News{ID: id, Title: title, Slug: slug, Content: content, Status: status, Tags: tags, Topic: topic, CreatedAt: Date.Now()}
}

func (n *News) ToSliceNewsTag() (sliceNewsTag []NewsTag) {
	for _, v := range n.Tags {
		sliceNewsTag = append(sliceNewsTag, NewsTag{TagID: v, NewsID: n.ID})
	}
	return
}

func (n *News) Update(new News) {
	if new.Title != "" {
		n.Title = new.Title
	}
	if new.Slug != "" {
		n.Slug = new.Slug
	}
	if new.Content != "" {
		n.Content = new.Content
	}
	if new.Status > 0 {
		n.Status = new.Status
	}
	if new.Topic != "" {
		n.Topic = new.Topic
	}
	if len(new.Tags) > 0 {
		n.Tags = new.Tags
	}
	if new.DeletedAt != nil {
		n.DeletedAt = new.DeletedAt
	}
}

func (n *News) Delete() {
	n.Status = NewsDeleted
	now := Date.Now()
	n.DeletedAt = &now
}

func (n *News) ToNewsDto(tagsMap ...map[string]Tag) *NewsDto {
	tags := n.Tags
	if len(tagsMap) > 0 {
		tags = []string{}
		i := 0
		for range tagsMap[0] {
			tags = append(tags, (tagsMap[0])[n.Tags[i]].Name)
			i++
		}
	}
	res := NewsDto{
		ID:      n.ID,
		Title:   n.Title,
		Slug:    n.Slug,
		Content: n.Content,
		Topic:   n.Topic,
		Status:  n.Status.String(),
		Tags:    tags,
	}
	return &res
}

func (n *News) SetTagsFromMapTags(mapTags map[string]Tag) {
	var newTag []string
	for _, tagID := range n.Tags {
		if tag, ok := mapTags[tagID]; ok {
			newTag = append(newTag, tag.Name)
		}
	}
	n.Tags = newTag
}

type SliceNews []News

func (s *SliceNews) ToSliceNewsDto(mapTags ...func() map[string]Tag) *SliceNewsDto {
	var search bool
	var tags map[string]Tag
	if len(mapTags) > 0 {
		tags = mapTags[0]()
		search = true
	}
	var res SliceNewsDto
	for _, news := range *s {
		if search {
			news.SetTagsFromMapTags(tags)
		}
		res = append(res, *news.ToNewsDto())
	}
	return &res
}

func (s *SliceNews) GetSliceTagIds() (tagIds []string) {
	duplication := map[string]bool{}
	for _, news := range *s {
		for _, tag := range news.Tags {
			if _, ok := duplication[tag]; ok {
				continue
			}
			duplication[tag] = true
			tagIds = append(tagIds, tag)
		}
	}
	return
}
