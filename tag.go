package kokizami

// Tag represents a tag
type Tag struct {
	ID    int
	Label string // TODO: change var name to label
}

// Tags represents array of Tag
type Tags []Tag

// TagRepository is an interface to fetch tags from repository
type TagRepository interface {
	FindByID(id int) (*Tag, error)
	FindAll() ([]*Tag, error)
	FindByKizamiID(kizamiID int) ([]*Tag, error)
	FindByLabels(labels []string) ([]*Tag, error)
	Insert(labels []string) error
	Delete(id int) error
}
