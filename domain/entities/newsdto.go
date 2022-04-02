package entities

import (
	"news/shared/failure"
	"strings"
)

type NewsDto struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Slug    string   `json:"slug"`
	Content string   `json:"content"`
	Status  string   `json:"status"`
	Tags    []string `json:"tags"`
	Topic   string   `json:"topic"`
}

func (n *NewsDto) Validate() error {
	var errString []string
	if n.Title == "" {
		errString = append(errString, "title can't be null")
	}
	if n.Content == "" {
		errString = append(errString, "content can't be null")
	}
	if n.Status == "" {
		errString = append(errString, "status can't be null")
	}
	_, err := StringToNewsStatus(n.Status)
	if err != nil {
		errString = append(errString, "status not valid")
	}
	if len(n.Tags) < 1 {
		errString = append(errString, "please insert tags")
	}
	if n.Topic == "" {
		errString = append(errString, "topic can't be null")
	}
	if len(errString) > 0 {
		return failure.BadRequestWithString(strings.Join(errString, ", "))
	}
	return nil
}

func (n *NewsDto) ToNews() (news *News, err error) {
	var status NewsStatus
	if n.Status != "" {
		status, err = StringToNewsStatus(n.Status)
	}
	if err != nil {
		return
	}
	news = NewNews(n.ID, n.Title, n.Slug, n.Content, status, n.Tags, n.Topic)
	return
}

type SliceNewsDto []NewsDto
