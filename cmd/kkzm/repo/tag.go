package repo

import (
	"database/sql"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
)

// TagRepo is an implementation of TagRepository
type TagRepo struct {
	db *sql.DB
}

// NewTagRepo returns an implementation of TagRepository with sqlite3
func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

func toTag(m *models.Tag) *kokizami.Tag {
	return &kokizami.Tag{
		ID:    m.ID,
		Label: m.Label,
	}
}

// FindByID finds a tag with specified ID
func (t *TagRepo) FindByID(id int) (*kokizami.Tag, error) {
	tag, err := models.TagByID(t.db, id)
	if err != nil {
		return nil, err
	}

	ret := &kokizami.Tag{
		ID:    tag.ID,
		Label: tag.Label,
	}

	return ret, nil
}

// FindAll returns all tags
func (t *TagRepo) FindAll() ([]*kokizami.Tag, error) {
	ms, err := models.AllTags(t.db)
	if err != nil {
		return nil, err
	}

	tags := make([]kokizami.Tag, len(ms))
	for i, v := range ms {
		tags[i].ID = v.ID
		tags[i].Label = v.Label
	}

	ret := make([]*kokizami.Tag, len(tags))
	for i := range tags {
		ret[i] = &tags[i]
	}

	return ret, nil
}

// FindByKizamiID returns tags they are held by specified kizami
func (t *TagRepo) FindByKizamiID(kizamiID int) ([]*kokizami.Tag, error) {
	ms, err := models.TagsByKizamiID(t.db, kizamiID)
	if err != nil {
		return nil, err
	}

	tags := make([]kokizami.Tag, len(ms))
	for i, v := range ms {
		tags[i].ID = v.ID
		tags[i].Label = v.Label
	}

	ret := make([]*kokizami.Tag, len(tags))
	for i := range tags {
		ret[i] = &tags[i]
	}

	return ret, nil
}

// FindByLabels finds tags by specified labels
func (t *TagRepo) FindByLabels(labels []string) ([]*kokizami.Tag, error) {
	ms, err := models.TagsByLabels(t.db, labels)
	if err != nil {
		return nil, err
	}

	ts := make([]*kokizami.Tag, len(ms))
	for i := range ms {
		ts[i] = toTag(ms[i])
	}

	return ts, nil
}

// Insert inserts tags with specified labels
func (t *TagRepo) Insert(labels []string) error {
	ts := models.Tags(make([]models.Tag, len(labels)))

	for i := range ts {
		// skip empty string
		if len(labels[i]) == 0 {
			ts = ts[:len(ts)-1]
			continue
		}
		ts[i].Label = labels[i]
	}

	return ts.BulkInsert(t.db)
}

// Delete deletes a tag by specified ID
func (t *TagRepo) Delete(id int) error {
	m, err := models.TagByID(t.db, id)
	if err != nil {
		return err
	}

	return m.Delete(t.db)
}
