package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/config"
)

type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// VersionInfo returns version information as JSON
func (h *VersionHandler) VersionInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	versionInfo := config.GetVersionInfo()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"data":    versionInfo,
		"message": "Version information retrieved successfully",
	})
}

// VersionText returns version information as plain text
func (h *VersionHandler) VersionText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(config.GetFullVersionInfo()))
}
