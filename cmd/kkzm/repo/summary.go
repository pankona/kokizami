package repo

import (
	"database/sql"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
)

// SummaryRepo is an implementation of SummaryRepository
type SummaryRepo struct {
	db *sql.DB
}

// NewSummaryRepo returns a struct that implements SummaryRepository with sqlite3
func NewSummaryRepo(db *sql.DB) *SummaryRepo {
	return &SummaryRepo{db: db}
}

// ElapsedOfMonthByDesc returns an array of Elapsed time to summarize them by desc
func (r *SummaryRepo) ElapsedOfMonthByDesc(yyyymm string) ([]*kokizami.Elapsed, error) {
	ms, err := models.ElapsedOfMonthByDesc(r.db, yyyymm)
	if err != nil {
		return nil, err
	}

	es := make([]kokizami.Elapsed, len(ms))
	for i := range ms {
		es[i].Tag = ms[i].Tag
		es[i].Desc = ms[i].Desc
		es[i].Count = ms[i].Count
		es[i].Elapsed = ms[i].Elapsed
	}

	ret := make([]*kokizami.Elapsed, len(es))
	for i := range es {
		ret[i] = &es[i]
	}

	return ret, nil
}

// ElapsedOfMonthByTag returns an array of Elapsed time to summarize them by tag
func (r *SummaryRepo) ElapsedOfMonthByTag(yyyymm string) ([]*kokizami.Elapsed, error) {
	ms, err := models.ElapsedOfMonthByTag(r.db, yyyymm)
	if err != nil {
		return nil, err
	}

	es := make([]kokizami.Elapsed, len(ms))
	for i := range ms {
		es[i].Tag = ms[i].Tag
		es[i].Desc = ms[i].Desc
		es[i].Count = ms[i].Count
		es[i].Elapsed = ms[i].Elapsed
	}

	ret := make([]*kokizami.Elapsed, len(es))
	for i := range es {
		ret[i] = &es[i]
	}

	return ret, nil
}
