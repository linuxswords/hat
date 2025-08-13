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
		"add":      func(a, b int) int { return a + b },
		"div":      func(a, b float64) float64 { return a / b },
	}
	tmpl, err := template.New("base").Funcs(templateFuncs).ParseFiles(
		"templates/layouts/base.html",
		"templates/partials/header.html",
		"templates/partials/navigation.html",
		fmt.Sprintf("templates/%s.html", contentTemplate),
	)
	if err != nil {
		fmt.Printf("Template parsing error: %v\n", err)
		c.String(500, "Template parsing error: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "layouts/base", data)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
	}
}
