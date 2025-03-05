package swagger2openapi3

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	swui "github.com/waroir20/swagger2-openapi3/swagger-ui"
)

type Config struct {
	DocsBasePath string

	PathToSpec   string
	SpecBasePath string

	StaticFilesBasePath string
	PathToStaticFiles   string

	Title string
}

func AddRoute(router *gin.Engine, cfg Config) {
	if cfg.PathToSpec == "" {
		panic("PathToSpec must be provided")
	}
	var optsPrelim []swui.Option
	specURL := ""
	if strings.HasPrefix(cfg.PathToSpec, "http://") ||
		strings.HasPrefix(cfg.PathToSpec, "https://") {
		specURL = cfg.PathToSpec
	} else {
		if cfg.SpecBasePath == "" {
			panic("SpecBasePath must be provided")
		} else {
			split := strings.SplitN(cfg.PathToSpec, "/", 3)
			folder := fmt.Sprintf("%s/%s", split[0], split[1])
			router.Static(cfg.SpecBasePath, folder)
			specURL = fmt.Sprintf("%s/%s", cfg.SpecBasePath, split[2])
			optsPrelim = append(optsPrelim, swui.WithRelativeSpecPath())
		}
	}

	opts := []swui.Option{
		swui.WithBasePath(cfg.DocsBasePath),
		swui.WithHTMLTitle(cfg.Title),
		swui.WithSwagBundleBaseURL("https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/"),
		swui.WithSmartBearHeader(),
		swui.WithJSONEditor(),
		swui.WithSettingsUI(map[string]string{
			"filter":                   "true",
			"operationsSorter":         "'alpha'",
			"showExtensions":           "true",
			"showCommonExtensions":     "true",
			"syntaxHighlight":          "true",
			"tryItOutEnabled":          "true",
			"requestSnippetsEnabled":   "true",
			"deepLinking":              "true",
			"displayRequestDuration":   "true",
			"validatorUrl":             "'https://validator.swagger.io/validator'",
			"defaultModelsExpandDepth": "1",
		}),
	}
	opts = append(optsPrelim, opts...)

	if cfg.PathToStaticFiles != "" {
		router.Static(cfg.StaticFilesBasePath, cfg.PathToStaticFiles)
		opts = append(opts, swui.WithFaviconBaseURL(cfg.StaticFilesBasePath))
	}

	handlerFunc := swui.NewHandlerWithConfig(specURL, opts...)
	router.GET(cfg.DocsBasePath, gin.WrapH(handlerFunc))
}
