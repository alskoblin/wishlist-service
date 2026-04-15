package openapi

import (
	"embed"
	"net/http"
)

//go:embed static/openapi.yaml
var assets embed.FS

func SpecHandler(w http.ResponseWriter, _ *http.Request) {
	data, err := assets.ReadFile("static/openapi.yaml")
	if err != nil {
		http.Error(w, "openapi file not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write(data)
}
