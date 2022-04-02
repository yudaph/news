package entities_test

import (
	"github.com/magiconair/properties/assert"
	"news/domain/entities"
	"news/shared/Date"
	"news/shared/IDGEN"
	"testing"
	"time"
)

func TestNews(t *testing.T) {

	t.Run("testCreateNews", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		sliceTest := []struct {
			testTitle, id, title, slug, content, topic string
			status                                     entities.NewsStatus
			tags                                       []string
			expected                                   entities.News
		}{
			{
				testTitle: "without slug and id",
				title:     "first title",
				content:   "content first",
				topic:     "football",
				status:    entities.NewsDraft,
				tags:      []string{"tags1", "tags2"},
				expected: entities.News{
					ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:     "first title",
					Slug:      "first-title",
					Content:   "content first",
					Topic:     "football",
					Status:    entities.NewsDraft,
					Tags:      []string{"tags1", "tags2"},
					CreatedAt: mockTime,
				},
			},
			{
				testTitle: "without slug",
				id:        "d2668631",
				title:     "first title",
				content:   "content first",
				topic:     "football",
				status:    entities.NewsPublish,
				tags:      []string{"tags1", "tags2"},
				expected: entities.News{
					ID:        "d2668631",
					Title:     "first title",
					Slug:      "first-title",
					Content:   "content first",
					Topic:     "football",
					Status:    entities.NewsPublish,
					Tags:      []string{"tags1", "tags2"},
					CreatedAt: mockTime,
				},
			},
			{
				testTitle: "without id",
				title:     "first title",
				slug:      "first-title",
				content:   "content first",
				topic:     "football",
				status:    entities.NewsDeleted,
				tags:      []string{"tags1", "tags2"},
				expected: entities.News{
					ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:     "first title",
					Slug:      "first-title",
					Content:   "content first",
					Topic:     "football",
					Status:    entities.NewsDeleted,
					Tags:      []string{"tags1", "tags2"},
					CreatedAt: mockTime,
				},
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := entities.NewNews(test.id, test.title, test.slug, test.content, test.status, test.tags, test.topic)
				assert.Equal(t, *actual, test.expected)
			})
		}
	})

	t.Run("testToSliceNewsTag", func(t *testing.T) {
		sliceTest := []struct {
			testTitle string
			input     entities.News
			expected  []entities.NewsTag
		}{
			{
				testTitle: "success",
				input: entities.News{
					ID:   "idnews",
					Tags: []string{"idtags1", "idtags2"},
				},
				expected: []entities.NewsTag{
					{
						NewsID: "idnews",
						TagID:  "idtags1",
					},
					{
						NewsID: "idnews",
						TagID:  "idtags2",
					},
				},
			},
		}
		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := test.input.ToSliceNewsTag()
				assert.Equal(t, actual, test.expected)
			})
		}
	})

	t.Run("testToSliceNewsTag", func(t *testing.T) {
		sliceTest := []struct {
			testTitle string
			input     entities.News
			expected  []entities.NewsTag
		}{
			{
				testTitle: "success",
				input: entities.News{
					ID:   "idnews",
					Tags: []string{"idtags1", "idtags2"},
				},
				expected: []entities.NewsTag{
					{
						NewsID: "idnews",
						TagID:  "idtags1",
					},
					{
						NewsID: "idnews",
						TagID:  "idtags2",
					},
				},
			},
		}
		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := test.input.ToSliceNewsTag()
				assert.Equal(t, actual, test.expected)
			})
		}
	})

	t.Run("testNewsUpdate", func(t *testing.T) {
		sliceTest := []struct {
			testTitle string
			input     entities.News
			input2    entities.News
			expected  entities.News
		}{
			{
				testTitle: "success",
				input: entities.News{
					ID:      "id",
					Title:   "title",
					Content: "content",
					Topic:   "topic",
					Status:  entities.NewsDraft,
					Tags:    []string{"idtags1", "idtags2"},
				},
				input2: entities.News{
					Title:   "title update",
					Content: "content update",
					Topic:   "topic update",
					Status:  entities.NewsPublish,
					Tags:    []string{"tags update", "tags2 update"},
				},
				expected: entities.News{
					ID:      "id",
					Title:   "title update",
					Content: "content update",
					Topic:   "topic update",
					Status:  entities.NewsPublish,
					Tags:    []string{"tags update", "tags2 update"},
				},
			}, {
				testTitle: "delete",
				input: entities.News{
					ID:      "id",
					Title:   "title",
					Content: "content",
					Topic:   "topic",
					Status:  entities.NewsDraft,
					Tags:    []string{"idtags1", "idtags2"},
				},
				input2: entities.News{
					Status: entities.NewsDeleted,
				},
				expected: entities.News{
					ID:      "id",
					Title:   "title",
					Content: "content",
					Topic:   "topic",
					Status:  entities.NewsDeleted,
					Tags:    []string{"idtags1", "idtags2"},
				},
			},
		}
		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := test.input
				actual.Update(test.input2)
				assert.Equal(t, actual, test.expected)
			})
		}
	})

	t.Run("testNewsDeleted", func(t *testing.T) {
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}
		sliceTest := []struct {
			testTitle string
			input     entities.News
			expected  entities.News
		}{
			{
				testTitle: "success",
				input: entities.News{
					ID:      "id",
					Title:   "title",
					Content: "content",
					Topic:   "topic",
					Status:  entities.NewsDraft,
					Tags:    []string{"idtags1", "idtags2"},
				},
				expected: entities.News{
					ID:        "id",
					Title:     "title",
					Content:   "content",
					Topic:     "topic",
					Status:    entities.NewsDeleted,
					Tags:      []string{"idtags1", "idtags2"},
					DeletedAt: &mockTime,
				},
			},
		}
		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := test.input
				actual.Delete()
				assert.Equal(t, actual, test.expected)
			})
		}
	})

	t.Run("testNewsToNewsDto", func(t *testing.T) {
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}
		sliceTest := []struct {
			testTitle string
			input     entities.News
			input2    map[string]entities.Tag
			expected  *entities.NewsDto
		}{
			{
				testTitle: "success",
				input: entities.News{
					ID:      "id",
					Title:   "title",
					Slug:    "title",
					Content: "content",
					Topic:   "topic",
					Status:  entities.NewsDraft,
					Tags:    []string{"idtags1", "idtags2"},
				},
				input2: map[string]entities.Tag{
					"idtags1": entities.Tag{
						Name: "tag 1",
					},
					"idtags2": entities.Tag{
						Name: "tag 2",
					},
				},
				expected: &entities.NewsDto{
					ID:      "id",
					Title:   "title",
					Slug:    "title",
					Content: "content",
					Topic:   "topic",
					Status:  "draft",
					Tags:    []string{"tag 1", "tag 2"},
				},
			},
		}
		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				actual := test.input.ToNewsDto(test.input2)
				assert.Equal(t, actual, test.expected)
			})
		}
	})
}
