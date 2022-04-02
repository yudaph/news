package entities

type NewsTag struct {
	NewsID string `db:"news_id"`
	TagID  string `db:"tag_id"`
}

type SliceNewsTag []NewsTag

func (s *SliceNewsTag) ToMapTag() map[string][]string {
	result := map[string][]string{}
	for _, newsTag := range *s {
		result[newsTag.NewsID] = append(result[newsTag.NewsID], newsTag.TagID)
	}
	return result
}
