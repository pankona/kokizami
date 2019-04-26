package kokizami

import "github.com/pankona/kokizami/models"

// Tag represents a tag
type Tag struct {
	ID  int
	Tag string
}

func toTag(m *models.Tag) *Tag {
	return &Tag{
		ID:  m.ID,
		Tag: m.Tag,
	}
}
