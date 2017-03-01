package kokizami

// Tagger represents interface of tag
type Tagger interface {
	ID() int
	Label() string
}

type tag struct {
	id    int
	label string
}

// AddTag adds new tag by specified label
func AddTag(l string) error {
	return nil
}

// RemoveTag removes a tag by specified id
func RemoveTag(id int) error {
	return nil
}

// GetTag returns a tag by specified id
func GetTag(id int) (Tagger, error) {
	return nil, nil
}

// GetTagList returns list of tags
func GetTagList() ([]Tagger, error) {
	return nil, nil
}
