package repo

import (
	"database/sql"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

func toTag(m *models.Tag) *kokizami.Tag {
	return &kokizami.Tag{
		ID:    m.ID,
		Label: m.Label,
	}
}

func (t *TagRepo) FindTagByID(id int) (*kokizami.Tag, error) {
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

func (t *TagRepo) FindAllTags() ([]*kokizami.Tag, error) {
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

func (t *TagRepo) FindTagsByKizamiID(kizamiID int) ([]*kokizami.Tag, error) {
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

func (t *TagRepo) FindTagsByLabels(labels []string) ([]*kokizami.Tag, error) {
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

func (t *TagRepo) InsertTags(labels []string) error {
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

func (t *TagRepo) Delete(id int) error {
	m, err := models.TagByID(t.db, id)
	if err != nil {
		return err
	}

	return m.Delete(t.db)
}
