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
