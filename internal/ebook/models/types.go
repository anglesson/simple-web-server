package models

type EbookRequest struct {
	Title       string `validate:"required,min=5,max=10"`
	Description string `validate:"required"`
	Value       string `validate:"required"`
	Status      bool   `validate:"required"`
}
