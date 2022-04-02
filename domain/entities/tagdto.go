package entities

import "news/shared/failure"

type TagDto struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (t *TagDto) ToTag() *Tag {
	return newTag(t.ID, t.Name)
}

func (t *TagDto) Validate() error {
	if t.Name == "" {
		return failure.BadRequestWithString("name can't be null")
	}
	return nil
}

type CreateTag struct {
	Name string `json:"name"`
}

func (c *CreateTag) Validate() error {
	if c.Name == "" {
		return failure.BadRequestWithString("name can't be empty")
	}
	return nil
}

func (c *CreateTag) ToTag() *Tag {
	return newTag("", c.Name)
}
