package handlers

import (
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/pkg/template"
)

// TODO: Render errors dynamically
func ErrorView(w http.ResponseWriter, r *http.Request, code int) {
	codeStr := strconv.Itoa(code)
	template.View(w, r, codeStr+"-error", nil, "guest")
}
