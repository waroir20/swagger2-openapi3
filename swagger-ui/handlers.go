package swagger_ui

import (
	"encoding/json"
	"html/template"
	"net/http"
)

const (
	assetsBase  = "{{ .BasePath }}"
	faviconBase = "{{ .BasePath }}"
)

// Handler handles swagger UI request.
type Handler struct {
	config
	ConfigJson template.JS
	tpl        *template.Template
}

// NewHandlerWithConfig creates HTTP handler for Swagger UI.
func NewHandlerWithConfig(docsLocation string, opts ...Option) *Handler {
	defaultOpts := []Option{
		WithFaviconBaseURL(faviconBase),
		WithSwagBundleBaseURL(assetsBase),
		WithHTMLTitle("Swagger UI"),
		WithSettingsUI(make(map[string]string)),
	}
	var cfg config
	cfg.SwaggerJSON = docsLocation
	for _, do := range defaultOpts {
		do(&cfg)
	}
	for _, o := range opts {
		o(&cfg)
	}
	return newHandlerWithConfig(cfg)
}

// NewHandlerWithConfig returns a HTTP handler for swagger UI.
func newHandlerWithConfig(cfg config) *Handler {
	h := &Handler{
		config: cfg,
	}

	j, err := json.Marshal(h.config)
	if err != nil {
		panic(err)
	}

	h.ConfigJson = template.JS(j)

	templat := IndexTpl(cfg.swagBundleBaseURL, cfg.faviconBaseURL, cfg)
	h.tpl, err = template.New("index").Parse(templat)
	if err != nil {
		panic(err)
	}
	return h
}

// ServeHTTP implements http.Handler interface to handle swagger UI request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r
	w.Header().Set("Content-Type", "text/html")

	if err := h.tpl.Execute(w, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
