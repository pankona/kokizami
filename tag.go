package kokizami

import "github.com/pankona/kokizami/models"

// Tag represents a tag
type Tag struct {
	ID  int
	Tag string // TODO: change var name to label
}

func toTag(m *models.Tag) Tag {
	return Tag{
		ID:  m.ID,
		Tag: m.Tag,
	}
}

type TagRepository interface {
	FindTagByID(id int) (*models.Tag, error)
	FindAllTags() ([]*models.Tag, error)
	FindTagsByKizamiID(kizamiID int) ([]*models.Tag, error)
	FindTagsByLabels(labels []string) ([]*models.Tag, error)
	InsertTags(models.Tags) error
}
