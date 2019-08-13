package kokizami

import "github.com/pankona/kokizami/models"

// Tag represents a tag
type Tag struct {
	ID  int
	Tag string // TODO: change var name to label
}

type Tags []Tag

func toTag(m *models.Tag) Tag {
	return Tag{
		ID:  m.ID,
		Tag: m.Tag,
	}
}

type TagRepository interface {
	FindTagByID(id int) (*Tag, error)
	FindAllTags() ([]*Tag, error)
	FindTagsByKizamiID(kizamiID int) ([]*Tag, error)
	FindTagsByLabels(labels []string) ([]*Tag, error)
	InsertTags(labels []string) error
	Delete(id int) error
}
