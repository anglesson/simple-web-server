package repositories

type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}

func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 1000:
		pageSize = 1000
	case pageSize <= 0:
		pageSize = 10
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) GetLimit() int {
	return p.PageSize
}

func (p Pagination) End() int {
	if int(p.Total) <= 0 {
		return 0
	}
	if p.Page*p.PageSize < int(p.Total) {
		return int(p.Total)
	}
	return p.Page * p.PageSize
}

func (p Pagination) Start() int {
	if int(p.Total) <= 0 {
		return 0
	}
	return (p.Page-1)*p.PageSize + 1
}

func (p Pagination) PrevPage() int {
	if p.Page <= 1 {
		return 1
	}
	return p.Page - 1
}

func (p Pagination) NextPage() int {
	return p.Page + 1
}

func (p Pagination) TotalPages() int {
	totalPages := p.Total / int64(p.PageSize)
	return int(totalPages)
}
