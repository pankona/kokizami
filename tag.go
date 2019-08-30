package kokizami

// Tag represents a tag
type Tag struct {
	ID    int
	Label string // TODO: change var name to label
}

type Tags []Tag

type TagRepository interface {
	FindByID(id int) (*Tag, error)
	FindAll() ([]*Tag, error)
	FindByKizamiID(kizamiID int) ([]*Tag, error)
	FindByLabels(labels []string) ([]*Tag, error)
	Insert(labels []string) error
	Delete(id int) error
}
