package main

import (
	"database/sql"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
)

type summaryRepo struct {
	db *sql.DB
}

func (r *summaryRepo) ElapsedOfMonthByDesc(yyyymm string) ([]*kokizami.Elapsed, error) {
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

func (r *summaryRepo) ElapsedOfMonthByTag(yyyymm string) ([]*kokizami.Elapsed, error) {

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
