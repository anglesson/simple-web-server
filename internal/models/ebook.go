package models

import (
	"regexp"
	"strings"

	"github.com/anglesson/simple-web-server/pkg/utils"
	"gorm.io/gorm"
)

type Ebook struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	SalesPage   string  `json:"sales_page"` // Conteúdo da página de vendas
	Value       float64 `json:"value"`
	Status      bool    `json:"status"`
	Image       string  `json:"image"`
	Slug        string  `json:"slug" gorm:"uniqueIndex"` // URL amigável
	CreatorID   uint    `json:"creator_id"`
	Creator     Creator `gorm:"foreignKey:CreatorID"`
	Files       []*File `gorm:"many2many:ebook_files;"`

	// Campos para SEO e marketing
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	Keywords        string `json:"keywords"`

	// Estatísticas
	Views int `json:"views" gorm:"default:0"`
	Sales int `json:"sales" gorm:"default:0"`
}

func NewEbook(title, description, salesPage string, value float64, creator Creator) *Ebook {
	return &Ebook{
		Title:       title,
		Description: description,
		SalesPage:   salesPage,
		Value:       value,
		Status:      true,
		CreatorID:   creator.ID,
		Slug:        generateSlug(title),
	}
}

func (e *Ebook) GetValue() string {
	return utils.FloatToBRL(e.Value)
}

func (e *Ebook) GetLastUpdate() string {
	return e.UpdatedAt.Format("02-01-2006 15:04")
}

func (e *Ebook) AddFile(file *File) {
	e.Files = append(e.Files, file)
}

func (e *Ebook) RemoveFile(fileID uint) {
	for i, file := range e.Files {
		if file.ID == fileID {
			e.Files = append(e.Files[:i], e.Files[i+1:]...)
			break
		}
	}
}

func (e *Ebook) GetTotalFileSize() int64 {
	var total int64
	for _, file := range e.Files {
		total += file.FileSize
	}
	return total
}

func (e *Ebook) GetFileCount() int {
	return len(e.Files)
}

func (e *Ebook) IncrementViews() {
	e.Views++
}

func (e *Ebook) IncrementSales() {
	e.Sales++
}

// generateSlug cria uma URL amigável baseada no título
func generateSlug(title string) string {
	// Implementação básica - pode ser melhorada com biblioteca de slug
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "ç", "c")
	slug = strings.ReplaceAll(slug, "ã", "a")
	slug = strings.ReplaceAll(slug, "á", "a")
	slug = strings.ReplaceAll(slug, "à", "a")
	slug = strings.ReplaceAll(slug, "â", "a")
	slug = strings.ReplaceAll(slug, "é", "e")
	slug = strings.ReplaceAll(slug, "ê", "e")
	slug = strings.ReplaceAll(slug, "í", "i")
	slug = strings.ReplaceAll(slug, "ó", "o")
	slug = strings.ReplaceAll(slug, "ô", "o")
	slug = strings.ReplaceAll(slug, "ú", "u")
	slug = strings.ReplaceAll(slug, "ü", "u")
	slug = strings.ReplaceAll(slug, "ñ", "n")

	// Remove caracteres especiais
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")

	// Remove hífens duplicados
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove hífens no início e fim
	slug = strings.Trim(slug, "-")

	return slug
}
