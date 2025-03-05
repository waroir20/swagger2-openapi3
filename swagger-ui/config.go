package swagger_ui

// config is used for Swagger UI handler configuration.
type config struct {
	faviconBaseURL     string            `json:"-"`
	swagBundleBaseURL  string            `json:"-"`
	Title              string            `json:"title"`              // Title of index file.
	SwaggerJSON        string            `json:"swaggerJsonUrl"`     // URL to openapi.json/swagger.json document specification.
	BasePath           string            `json:"basePath"`           // Base URL to docs.
	ShowTopBar         bool              `json:"showTopBar"`         // Show navigation top bar, hidden by default.
	HideCurl           bool              `json:"hideCurl"`           // Hide curl code snippet.
	JsonEditor         bool              `json:"jsonEditor"`         // Enable visual json editor support (experimental, can fail with complex schemas).
	PreAuthorizeApiKey map[string]string `json:"preAuthorizeApiKey"` // Map of security name to key value.

	// SettingsUI contains keys and plain javascript values of SwaggerUIBundle configuration.
	// Overrides default values.
	// See https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/ for available options.
	settingsUI map[string]string `json:"-"`
}
