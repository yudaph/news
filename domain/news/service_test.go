package news_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"news/domain/entities"
	"news/domain/news"
	news_mock "news/domain/news/mock"
	tag_mock "news/domain/tag/mock"
	"news/shared/Date"
	"news/shared/IDGEN"
	"news/shared/failure"
	"testing"
	"time"
)

func TestNewsService(t *testing.T) {
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

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)

		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, input entities.NewsDto)
			input          entities.NewsDto
			expectedResult *entities.NewsDto
			expectedError  error
		}{
			{
				testTitle: "create success",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, input entities.NewsDto) {
					dto, _ := input.ToNews()
					repo.EXPECT().CreateNews(ctx, dto).Return(nil)
				},
				input: entities.NewsDto{
					Title:   "first title",
					Content: "content first",
					Topic:   "football",
					Status:  "deleted",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedResult: &entities.NewsDto{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "deleted",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedError: nil,
			},
			{
				testTitle: "error wrong status",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, input entities.NewsDto) {

				},
				input: entities.NewsDto{
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "delete",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedResult: nil,
				expectedError:  failure.NotFound("status not found"),
			},
			{
				testTitle: "error repository",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, input entities.NewsDto) {
					dto, _ := input.ToNews()
					repo.EXPECT().CreateNews(ctx, dto).Return(failure.InternalServerError)
				},
				input: entities.NewsDto{
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "deleted",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedResult: nil,
				expectedError:  failure.InternalServerError,
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, test.input)
				actual, err := service.Create(ctx, &test.input)
				assert.Equal(t, err, test.expectedError)
				if test.expectedResult != nil {
					assert.Equal(t, *actual, *test.expectedResult)
				} else {
					assert.Equal(t, actual, test.expectedResult)
				}
			})
		}
	})

	t.Run("testGetNewsBySlug", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)

		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string)
			input          string
			expectedResult *entities.NewsDto
			expectedError  error
		}{
			{
				testTitle: "success from cache",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string) {
					cache.EXPECT().GetNews(ctx, "slug:"+slug).Return(&entities.NewsDto{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "deleted",
						Tags:    []string{"tags1", "tags2"},
					}, nil)
				},
				input: "news-title",
				expectedResult: &entities.NewsDto{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "deleted",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedError: nil,
			},
			{
				testTitle: "success from DB",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string) {
					cache.EXPECT().GetNews(ctx, "slug:"+slug).Return(nil, failure.InternalServerError)
					repo.EXPECT().GetNewsBySlug(ctx, slug).Return(&entities.News{
						ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:     "first title",
						Slug:      "first-title",
						Content:   "content first",
						Topic:     "football",
						Status:    entities.NewsPublish,
						CreatedAt: mockTime,
						DeletedAt: nil,
						Tags:      []string{"id1", "id2"},
					}, nil)
					tagRepo.EXPECT().GetTagByIds(ctx, []string{"id1", "id2"}).Return(&entities.Tags{
						{
							ID:     "id1",
							Name:   "tags1",
							Status: entities.TagActive,
						},
						{
							ID:     "id2",
							Name:   "tags2",
							Status: entities.TagActive,
						},
					}, nil)
					cache.EXPECT().SetNews(ctx, "slug:"+slug, &entities.NewsDto{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}).Return(nil)
				},
				input: "news-title",
				expectedResult: &entities.NewsDto{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedError: nil,
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, mockTagRepo, mockCache, test.input)
				actual, err := service.GetBySlug(ctx, test.input)
				assert.Equal(t, err, test.expectedError)
				if test.expectedResult != nil {
					assert.Equal(t, *actual, *test.expectedResult)
				} else {
					assert.Equal(t, actual, test.expectedResult)
				}
			})
		}
	})

	t.Run("testGetNewsByTopic", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)

		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string)
			input          string
			expectedResult *entities.SliceNewsDto
			expectedError  error
		}{
			{
				testTitle: "success from cache",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string) {
					cache.EXPECT().GetSliceNews(ctx, "topic:"+slug).Return(&entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "deleted",
						Tags:    []string{"tags1", "tags2"},
					}}, nil)
				},
				input: "topic",
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "deleted",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
			{
				testTitle: "success from DB",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug string) {
					cache.EXPECT().GetSliceNews(ctx, "topic:"+slug).Return(nil, failure.InternalServerError)
					repo.EXPECT().GetNewsByTopic(ctx, slug).Return(&entities.SliceNews{{
						ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:     "first title",
						Slug:      "first-title",
						Content:   "content first",
						Topic:     "football",
						Status:    entities.NewsPublish,
						CreatedAt: mockTime,
						DeletedAt: nil,
						Tags:      []string{"id1", "id2"},
					}}, nil)
					tagRepo.EXPECT().GetTagByIds(ctx, []string{"id1", "id2"}).Return(&entities.Tags{
						{
							ID:     "id1",
							Name:   "tags1",
							Status: entities.TagActive,
						},
						{
							ID:     "id2",
							Name:   "tags2",
							Status: entities.TagActive,
						},
					}, nil)
					cache.EXPECT().SetSliceNews(ctx, "topic:"+slug, &entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}}).Return(nil)
				},
				input: "topic",
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, mockTagRepo, mockCache, test.input)
				actual, err := service.GetByTopic(ctx, test.input)
				assert.Equal(t, err, test.expectedError)
				if test.expectedResult != nil {
					assert.Equal(t, *actual, *test.expectedResult)
				} else {
					assert.Equal(t, actual, test.expectedResult)
				}
			})
		}
	})

	t.Run("testGetNewsByStatus", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)

		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug entities.NewsStatus)
			input          entities.NewsStatus
			expectedResult *entities.SliceNewsDto
			expectedError  error
		}{
			{
				testTitle: "success from cache",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug entities.NewsStatus) {
					cache.EXPECT().GetSliceNews(ctx, "status:"+slug.String()).Return(&entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}}, nil)
				},
				input: entities.NewsPublish,
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
			{
				testTitle: "success from DB",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache, slug entities.NewsStatus) {
					cache.EXPECT().GetSliceNews(ctx, "status:"+slug.String()).Return(nil, failure.InternalServerError)
					repo.EXPECT().GetNewsByStatus(ctx, slug).Return(&entities.SliceNews{{
						ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:     "first title",
						Slug:      "first-title",
						Content:   "content first",
						Topic:     "football",
						Status:    entities.NewsPublish,
						CreatedAt: mockTime,
						DeletedAt: nil,
						Tags:      []string{"id1", "id2"},
					}}, nil)
					tagRepo.EXPECT().GetTagByIds(ctx, []string{"id1", "id2"}).Return(&entities.Tags{
						{
							ID:     "id1",
							Name:   "tags1",
							Status: entities.TagActive,
						},
						{
							ID:     "id2",
							Name:   "tags2",
							Status: entities.TagActive,
						},
					}, nil)
					cache.EXPECT().SetSliceNews(ctx, "status:"+slug.String(), &entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}}).Return(nil)
				},
				input: entities.NewsPublish,
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, mockTagRepo, mockCache, test.input)
				actual, err := service.GetByStatus(ctx, test.input)
				assert.Equal(t, err, test.expectedError)
				if test.expectedResult != nil {
					assert.Equal(t, *actual, *test.expectedResult)
				} else {
					assert.Equal(t, actual, test.expectedResult)
				}
			})
		}
	})

	t.Run("testGetNewsAll", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)

		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache)
			expectedResult *entities.SliceNewsDto
			expectedError  error
		}{
			{
				testTitle: "success from cache",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache) {
					cache.EXPECT().GetSliceNews(ctx, "all:").Return(&entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}}, nil)
				},
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
			{
				testTitle: "success from DB",
				mockSetup: func(ctx context.Context, repo *news_mock.MockRepository, tagRepo *tag_mock.MockRepository, cache *news_mock.MockCache) {
					cache.EXPECT().GetSliceNews(ctx, "all:").Return(nil, failure.InternalServerError)
					repo.EXPECT().GetAllNews(ctx).Return(&entities.SliceNews{{
						ID:        "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:     "first title",
						Slug:      "first-title",
						Content:   "content first",
						Topic:     "football",
						Status:    entities.NewsPublish,
						CreatedAt: mockTime,
						DeletedAt: nil,
						Tags:      []string{"id1", "id2"},
					}}, nil)
					tagRepo.EXPECT().GetTagByIds(ctx, []string{"id1", "id2"}).Return(&entities.Tags{
						{
							ID:     "id1",
							Name:   "tags1",
							Status: entities.TagActive,
						},
						{
							ID:     "id2",
							Name:   "tags2",
							Status: entities.TagActive,
						},
					}, nil)
					cache.EXPECT().SetSliceNews(ctx, "all:", &entities.SliceNewsDto{{
						ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
						Title:   "first title",
						Slug:    "first-title",
						Content: "content first",
						Topic:   "football",
						Status:  "publish",
						Tags:    []string{"tags1", "tags2"},
					}}).Return(nil)
				},
				expectedResult: &entities.SliceNewsDto{{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				}},
				expectedError: nil,
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, mockTagRepo, mockCache)
				actual, err := service.GetAll(ctx)
				assert.Equal(t, err, test.expectedError)
				if test.expectedResult != nil {
					assert.Equal(t, *actual, *test.expectedResult)
				} else {
					assert.Equal(t, actual, test.expectedResult)
				}
			})
		}
	})

	t.Run("testUpdate", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)
		setup := func(ctx context.Context, repo *news_mock.MockRepository, input *entities.NewsDto, err error) {
			dto, _ := input.ToNews()
			repo.EXPECT().UpdateNews(ctx, dto).Return(err)
		}
		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, input *entities.NewsDto, err error)
			input          *entities.NewsDto
			expectedResult error
		}{
			{
				testTitle: "update success",
				mockSetup: setup,
				input: &entities.NewsDto{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedResult: nil,
			},
			{
				testTitle: "update fail",
				mockSetup: setup,
				input: &entities.NewsDto{
					ID:      "d2668631-1563-46bd-9498-5bfac7eed17a",
					Title:   "first title",
					Slug:    "first-title",
					Content: "content first",
					Topic:   "football",
					Status:  "publish",
					Tags:    []string{"tags1", "tags2"},
				},
				expectedResult: failure.NotFound("news not found"),
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, test.input, test.expectedResult)
				err := service.Update(ctx, test.input)
				assert.Equal(t, err, test.expectedResult)
			})
		}
	})

	t.Run("testDelete", func(t *testing.T) {
		//mock uuid
		IDGEN.NewUUID = func() string {
			return "d2668631-1563-46bd-9498-5bfac7eed17a"
		}

		//mock time
		mockTime := time.Now()
		Date.Now = func() time.Time {
			return mockTime
		}

		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNewsRepo := news_mock.NewMockRepository(ctrl)
		mockTagRepo := tag_mock.NewMockRepository(ctrl)
		mockCache := news_mock.NewMockCache(ctrl)
		service := news.NewService(mockNewsRepo, mockTagRepo, mockCache)
		setup := func(ctx context.Context, repo *news_mock.MockRepository, input string, err error) {
			repo.EXPECT().DeleteNews(ctx, input).Return(err)
		}
		sliceTest := []struct {
			testTitle      string
			mockSetup      func(ctx context.Context, repo *news_mock.MockRepository, input string, err error)
			input          string
			expectedResult error
		}{
			{
				testTitle:      "update success",
				mockSetup:      setup,
				input:          "d2668631-1563-46bd-9498-5bfac7eed17a",
				expectedResult: nil,
			},
			{
				testTitle:      "update fail",
				mockSetup:      setup,
				input:          "d2668631",
				expectedResult: failure.NotFound("news not found"),
			},
		}

		for _, test := range sliceTest {
			t.Run(test.testTitle, func(t *testing.T) {
				ctx := context.Background()
				test.mockSetup(ctx, mockNewsRepo, test.input, test.expectedResult)
				err := service.Delete(ctx, test.input)
				assert.Equal(t, err, test.expectedResult)
			})
		}
	})
}
