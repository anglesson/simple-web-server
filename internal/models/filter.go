package models

type ClientFilter struct {
	Term       string
	EbookID    uint
	Pagination *Pagination
}

type Pagination struct {
	Page       int
	Limit      int
	Total      int64
	Start      int
	End        int
	PrevPage   int
	NextPage   int
	PageSize   int
	TotalPages int
}

// NewPagination creates a new pagination with calculated fields
func NewPagination(page, limit int) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	start := (page-1)*limit + 1
	end := page * limit

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	nextPage := page + 1

	return &Pagination{
		Page:       page,
		Limit:      limit,
		PageSize:   limit,
		Start:      start,
		End:        end,
		PrevPage:   prevPage,
		NextPage:   nextPage,
		Total:      0,
		TotalPages: 0,
	}
}

// SetTotal updates the pagination with total count and recalculates fields
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))

	if p.End > int(total) {
		p.End = int(total)
	}

	if p.NextPage > p.TotalPages {
		p.NextPage = p.TotalPages
	}
}
