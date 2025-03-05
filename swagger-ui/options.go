package swagger_ui

import (
	"fmt"
	"strings"
)

type Option = func(d *config)

func WithFaviconBaseURL(url string) Option {
	return func(d *config) {
		if url != faviconBase {
			d.faviconBaseURL = strings.TrimSuffix(url, "/") + "/"
		} else {
			d.faviconBaseURL = url
		}
	}
}

func WithSwagBundleBaseURL(url string) Option {
	return func(d *config) {
		if url != assetsBase {
			d.swagBundleBaseURL = strings.TrimSuffix(url, "/") + "/"
		} else {
			d.swagBundleBaseURL = url
		}
	}
}

func WithSettingsUI(settings map[string]string) Option {
	if settings == nil {
		panic("map cannot be nil")
	}
	return func(d *config) {
		d.settingsUI = settings
	}
}

func WithCurlDisabled() Option {
	return func(d *config) {
		d.HideCurl = true
	}
}

func WithSmartBearHeader() Option {
	return func(d *config) {
		d.ShowTopBar = true
	}
}

func WithHTMLTitle(title string) Option {
	return func(d *config) {
		if title != "" {
			d.Title = title
		}
	}
}

func WithJSONEditor() Option {
	return func(d *config) {
		d.JsonEditor = true
	}
}

func WithRelativeSpecPath() Option {
	return func(d *config) {
		if !strings.HasPrefix(d.SwaggerJSON, "rel://") {
			d.SwaggerJSON = fmt.Sprintf("rel://%s", d.SwaggerJSON)
		}
	}
}

func WithBasePath(path string) Option {
	return func(d *config) {
		d.BasePath = strings.TrimSuffix(path, "/") + "/"
	}
}
