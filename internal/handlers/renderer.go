package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

func RenderWithLayout(c *gin.Context, contentTemplate string, data interface{}) {
	templateFuncs := map[string]any{
		"contains": strings.Contains,
	}
	tmpl, _ := template.New("base").Funcs(templateFuncs).ParseFiles(
		"templates/layouts/base.html",
		"templates/partials/header.html",
		"templates/partials/navigation.html",
		fmt.Sprintf("templates/%s.html", contentTemplate),
	)

	c.Header("Content-Type", "text/html")
	fmt.Println("Rendering template:", contentTemplate)
	fmt.Println("defined templates", tmpl.DefinedTemplates())
	tmpl.ExecuteTemplate(c.Writer, "layouts/base", data)
}
