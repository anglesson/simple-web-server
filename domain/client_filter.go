package domain

type ClientFilter struct {
	Term       string
	EbookID    uint
	Pagination *Pagination
}
