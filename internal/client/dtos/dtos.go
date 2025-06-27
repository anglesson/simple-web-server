package dtos

import (
	common_application "github.com/anglesson/simple-web-server/internal/common/application"
)

type CreateClientInput struct {
	Name         string
	CPF          string
	Phone        string
	BirthDate    string
	Email        string
	EmailCreator string
}

type CreateClientOutput struct {
	ID        int
	Name      string
	CPF       string
	Phone     string
	BirthDate string
	Email     string
	CreatedAt string
	UpdatedAt string
}

type UpdateClientInput struct {
	ID           uint
	Email        string
	Phone        string
	EmailCreator string
}

type ClientQuery struct {
	Term       string
	EbookID    uint
	Pagination *common_application.Pagination
}
