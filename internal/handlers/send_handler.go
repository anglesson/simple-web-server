package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func SendViewHandler(w http.ResponseWriter, r *http.Request) {
	allClients := []map[string]any{
		{
			"ID":   1,
			"Name": "João Silva",
			"CPF":  "123.456.789-00",
			"Contact": map[string]string{
				"Email": "joao.silva@email.com",
				"Phone": "(11) 98765-4321",
			},
		},
		{
			"ID":   2,
			"Name": "Maria Oliveira",
			"CPF":  "987.654.321-00",
			"Contact": map[string]string{
				"Email": "maria.oliveira@email.com",
				"Phone": "(11) 91234-5678",
			},
		},
		{
			"ID":   3,
			"Name": "Pedro Santos",
			"CPF":  "456.789.123-00",
			"Contact": map[string]string{
				"Email": "pedro.santos@email.com",
				"Phone": "(11) 99876-5432",
			},
		},
		{
			"ID":   4,
			"Name": "Ana Costa",
			"CPF":  "789.123.456-00",
			"Contact": map[string]string{
				"Email": "ana.costa@email.com",
				"Phone": "(11) 97777-8888",
			},
		},
		{
			"ID":   5,
			"Name": "Carlos Ferreira",
			"CPF":  "321.654.987-00",
			"Contact": map[string]string{
				"Email": "carlos.ferreira@email.com",
				"Phone": "(11) 96666-7777",
			},
		},
		{
			"ID":   6,
			"Name": "Juliana Almeida",
			"CPF":  "147.258.369-00",
			"Contact": map[string]string{
				"Email": "juliana.almeida@email.com",
				"Phone": "(11) 95555-6666",
			},
		},
		{
			"ID":   7,
			"Name": "Roberto Souza",
			"CPF":  "258.369.147-00",
			"Contact": map[string]string{
				"Email": "roberto.souza@email.com",
				"Phone": "(11) 94444-5555",
			},
		},
		{
			"ID":   8,
			"Name": "Fernanda Lima",
			"CPF":  "369.147.258-00",
			"Contact": map[string]string{
				"Email": "fernanda.lima@email.com",
				"Phone": "(11) 93333-4444",
			},
		},
		{
			"ID":   9,
			"Name": "Ricardo Mendes",
			"CPF":  "741.852.963-00",
			"Contact": map[string]string{
				"Email": "ricardo.mendes@email.com",
				"Phone": "(11) 92222-3333",
			},
		},
		{
			"ID":   10,
			"Name": "Patrícia Costa",
			"CPF":  "852.963.741-00",
			"Contact": map[string]string{
				"Email": "patricia.costa@email.com",
				"Phone": "(11) 91111-2222",
			},
		},
		{
			"ID":   11,
			"Name": "Marcelo Santos",
			"CPF":  "963.741.852-00",
			"Contact": map[string]string{
				"Email": "marcelo.santos@email.com",
				"Phone": "(11) 90000-1111",
			},
		},
		{
			"ID":   12,
			"Name": "Carla Oliveira",
			"CPF":  "159.357.486-00",
			"Contact": map[string]string{
				"Email": "carla.oliveira@email.com",
				"Phone": "(11) 98888-9999",
			},
		},
		{
			"ID":   13,
			"Name": "Bruno Pereira",
			"CPF":  "357.486.159-00",
			"Contact": map[string]string{
				"Email": "bruno.pereira@email.com",
				"Phone": "(11) 97777-8888",
			},
		},
		{
			"ID":   14,
			"Name": "Amanda Rodrigues",
			"CPF":  "486.159.357-00",
			"Contact": map[string]string{
				"Email": "amanda.rodrigues@email.com",
				"Phone": "(11) 96666-7777",
			},
		},
		{
			"ID":   15,
			"Name": "Diego Martins",
			"CPF":  "753.159.486-00",
			"Contact": map[string]string{
				"Email": "diego.martins@email.com",
				"Phone": "(11) 95555-6666",
			},
		},
		{
			"ID":   16,
			"Name": "Luciana Ferreira",
			"CPF":  "159.486.753-00",
			"Contact": map[string]string{
				"Email": "luciana.ferreira@email.com",
				"Phone": "(11) 94444-5555",
			},
		},
		{
			"ID":   17,
			"Name": "Gabriel Alves",
			"CPF":  "486.753.159-00",
			"Contact": map[string]string{
				"Email": "gabriel.alves@email.com",
				"Phone": "(11) 93333-4444",
			},
		},
		{
			"ID":   18,
			"Name": "Mariana Costa",
			"CPF":  "753.486.159-00",
			"Contact": map[string]string{
				"Email": "mariana.costa@email.com",
				"Phone": "(11) 92222-3333",
			},
		},
		{
			"ID":   19,
			"Name": "Rafael Silva",
			"CPF":  "159.753.486-00",
			"Contact": map[string]string{
				"Email": "rafael.silva@email.com",
				"Phone": "(11) 91111-2222",
			},
		},
		{
			"ID":   20,
			"Name": "Beatriz Santos",
			"CPF":  "486.159.753-00",
			"Contact": map[string]string{
				"Email": "beatriz.santos@email.com",
				"Phone": "(11) 90000-1111",
			},
		},
	}

	mockData := map[string]any{
		"Ebooks": []map[string]any{
			{
				"ID":            1,
				"Title":         "Marketing Digital para Iniciantes",
				"Description":   "Aprenda os fundamentos do marketing digital e como aplicá-los em seu negócio.",
				"GetValue":      "R$ 49,90",
				"GetLastUpdate": "15/03/2024",
			},
			{
				"ID":            2,
				"Title":         "Guia Completo de SEO",
				"Description":   "Técnicas avançadas de otimização para mecanismos de busca.",
				"GetValue":      "R$ 79,90",
				"GetLastUpdate": "10/03/2024",
			},
			{
				"ID":            3,
				"Title":         "Copywriting que Vende",
				"Description":   "Aprenda a escrever textos persuasivos que convertem visitantes em clientes.",
				"GetValue":      "R$ 59,90",
				"GetLastUpdate": "20/03/2024",
			},
		},
		"Clients": allClients,
	}

	template.View(w, r, "send_ebook", mockData, "admin")
}
