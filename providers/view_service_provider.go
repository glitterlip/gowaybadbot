package providers

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"goawaybot/services"
	"html/template"
	"io"
	"io/fs"
)

var Renderer *TemplateRenderer
var TemplateFs embed.FS
var t *template.Template

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
func RegisterViewServiceProvider(e *echo.Echo) {
	RegisterRenderer(e)
	RegisterErrorHandler(e)

}
func RegisterRenderer(e *echo.Echo) {
	t = template.New("verify").Funcs(template.FuncMap{
		"imageToBase64": services.ImgToBase64,
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
		"url": func(s string) template.URL {
			return template.URL(s)
		},
	})

	t.ParseFS(TemplateFs, "resources/templates/*.html")
	entries, _ := TemplateFs.ReadDir("resources/templates")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("resources/templates/" + entry.Name())
			walkDir(TemplateFs, "resources/templates/"+entry.Name())
		}
	}

	Renderer = &TemplateRenderer{
		templates: t,
	}
	e.Renderer = Renderer
}
func RegisterErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, context echo.Context) {
		fmt.Println(err.Error())
	}
}
func walkDir(fileSystem embed.FS, path string) {
	t.ParseFS(fileSystem, path+"/*.html")
	entries, _ := fs.ReadDir(fileSystem, path)
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println(path + "/" + entry.Name())
			walkDir(fileSystem, path+"/"+entry.Name())
		}
	}
}
