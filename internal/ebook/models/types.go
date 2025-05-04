package models

type EbookRequest struct {
	Title       string  `validate:"required,min=5,max=120" json:"title"`
	Description string  `validate:"required,max=120" json:"description"`
	Value       float64 `validate:"required,gt=0" json:"value"`
	Status      bool    `validate:"required" json:"status"`
}
