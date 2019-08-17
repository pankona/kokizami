package main

import (
	"database/sql"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
)

type tagRepo struct {
	db *sql.DB
}

func (t *tagRepo) FindTagByID(id int) (*kokizami.Tag, error) {
	tag, err := models.TagByID(t.db, id)
	if err != nil {
		return nil, err
	}

	ret := &kokizami.Tag{ID: tag.ID, Tag: tag.Tag}

	return ret, nil
}

func (t *tagRepo) FindAllTags() ([]*kokizami.Tag, error) {
	tags, err := models.AllTags(t.db)
	if err != nil {
		return nil, err
	}

	ret := make([]*kokizami.Tag, len(tags))
	for i, v := range tags {
		ret[i].ID = v.ID
		ret[i].Tag = v.Tag
	}

	return ret, nil
}

func (t *tagRepo) FindTagsByKizamiID(kizamiID int) ([]*kokizami.Tag, error) {
	tags, err := models.TagsByKizamiID(t.db, kizamiID)
	if err != nil {
		return nil, err
	}

	ret := make([]*kokizami.Tag, len(tags))
	for i, v := range tags {
		ret[i].ID = v.ID
		ret[i].Tag = v.Tag
	}

	return ret, nil
}

func toTag(m *models.Tag) *kokizami.Tag {
	return &kokizami.Tag{
		ID:  m.ID,
		Tag: m.Tag,
	}
}

func (t *tagRepo) FindTagsByLabels(labels []string) ([]*kokizami.Tag, error) {
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

func (t *tagRepo) InsertTags(labels []string) error {
	ts := models.Tags(make([]models.Tag, len(labels)))

	for i := range ts {
		// skip empty string
		if len(labels[i]) == 0 {
			ts = ts[:len(ts)-1]
			continue
		}
		ts[i].Tag = labels[i]
	}

	return ts.BulkInsert(t.db)
}

func (t *tagRepo) Delete(id int) error {
	m, err := models.TagByID(t.db, id)
	if err != nil {
		return err
	}
	return m.Delete(t.db)
}
