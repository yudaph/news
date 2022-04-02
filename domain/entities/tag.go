package entities

import (
	"news/shared/Date"
	"news/shared/IDGEN"
	"time"
)

type TagStatus int

const (
	TagActive TagStatus = iota + 1
	TagDelete
)

func (t TagStatus) String() string {
	sliceTagStatus := []string{"not found", "active", "delete"}
	return sliceTagStatus[t]
}

type Tag struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Status    TagStatus `db:"status"`
	CreatedAt time.Time `db:"createdAt"`
}

func newTag(id, name string) *Tag {
	if id == "" {
		id = IDGEN.NewUUID()
	}
	return &Tag{ID: id, Name: name, Status: TagActive, CreatedAt: Date.Now()}
}

func (t *Tag) UpdateTag(newTag *Tag) {
	if newTag.Name != "" {
		t.Name = newTag.Name
	}
	if newTag.Status > 0 {
		t.Status = newTag.Status
	}
}

func (t *Tag) Delete() {
	t.Status = TagDelete
}

func (t *Tag) ToDto() *TagDto {
	return &TagDto{
		ID:     t.ID,
		Name:   t.Name,
		Status: t.Status.String(),
	}
}

type Tags []Tag

func (t Tags) ToMapTags() map[string]Tag {
	result := map[string]Tag{}
	for _, tag := range t {
		result[tag.ID] = tag
	}
	return result
}

func (t Tags) ToTagsDto() *[]TagDto {
	var result []TagDto
	for _, tag := range t {
		result = append(result, *tag.ToDto())
	}
	return &result
}
