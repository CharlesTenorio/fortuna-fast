package model

import (
	"math"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Paginate struct {
	Total    int64       `json:"total"`
	Currente int64       `json:"current_page"`
	Last     int64       `json:"last_page"`
	Data     interface{} `json:"data"`
	Limit    int64       `json:"-"`
	Page     int64       `json:"-"`
}

func (p *Paginate) GetPaginatedOpts() *options.FindOptions {
	l := p.Limit
	skip := p.Page*p.Limit - p.Limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func (p *Paginate) Paginate(data interface{}) {
	p.Data = data
	p.Currente = p.Page
	d := float64(p.Total) / float64(p.Limit)
	p.Last = int64(math.Ceil(d))
}

func NewPaginate(limit, page, total int64) *Paginate {
	var limitL, pageL int64 = 10, 1

	if limit > 0 {
		limitL = limit
	}
	if page > 0 {
		pageL = page
	}

	return &Paginate{
		Limit: limitL,
		Page:  pageL,
		Total: total,
	}
}
