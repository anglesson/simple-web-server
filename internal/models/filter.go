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
	HasPrev    bool
	HasNext    bool
	Pages      []int
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

	// Calcular HasPrev e HasNext
	p.HasPrev = p.Page > 1
	p.HasNext = p.Page < p.TotalPages

	// Gerar array de páginas para navegação
	p.Pages = p.generatePageNumbers()
}

// generatePageNumbers generates an array of page numbers to display
func (p *Pagination) generatePageNumbers() []int {
	var pages []int

	if p.TotalPages <= 7 {
		// Se tem 7 páginas ou menos, mostrar todas
		for i := 1; i <= p.TotalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		// Se tem mais de 7 páginas, mostrar páginas estratégicas
		if p.Page <= 4 {
			// Páginas iniciais: 1, 2, 3, 4, 5, ..., TotalPages
			for i := 1; i <= 5; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, -1) // Separador
			pages = append(pages, p.TotalPages)
		} else if p.Page >= p.TotalPages-3 {
			// Páginas finais: 1, ..., TotalPages-4, TotalPages-3, TotalPages-2, TotalPages-1, TotalPages
			pages = append(pages, 1)
			pages = append(pages, -1) // Separador
			for i := p.TotalPages - 4; i <= p.TotalPages; i++ {
				pages = append(pages, i)
			}
		} else {
			// Páginas do meio: 1, ..., Page-1, Page, Page+1, ..., TotalPages
			pages = append(pages, 1)
			pages = append(pages, -1) // Separador
			for i := p.Page - 1; i <= p.Page+1; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, -1) // Separador
			pages = append(pages, p.TotalPages)
		}
	}

	return pages
}
