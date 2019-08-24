package kokizami

// Tag represents a tag
type Tag struct {
	ID    int
	Label string // TODO: change var name to label
}

type Tags []Tag

type TagRepository interface {
	FindTagByID(id int) (*Tag, error)
	FindAllTags() ([]*Tag, error)
	FindTagsByKizamiID(kizamiID int) ([]*Tag, error)
	FindTagsByLabels(labels []string) ([]*Tag, error)
	InsertTags(labels []string) error
	Delete(id int) error
}
