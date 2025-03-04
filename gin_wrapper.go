package swagger2openapi3

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
)

func AddRoute(router *gin.Engine, path, pathToSpec, title string) {
	split := strings.SplitN(pathToSpec, "/", 3)
	router.Static(path, fmt.Sprintf("%s/%s", split[0], split[1]))
	router.GET(path, gin.WrapH(v5emb.NewHandlerWithConfig(swgui.Config{
		Title:            title,
		SwaggerJSON:      fmt.Sprintf("%s/%s", path, split[2]),
		BasePath:         "https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/",
		InternalBasePath: path,
		ShowTopBar:       true,
		HideCurl:         false,
		JsonEditor:       true,
		SettingsUI: map[string]string{
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
		},
	})))
}
