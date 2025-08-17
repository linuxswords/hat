package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Navigate up from internal/handlers/ to project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func RenderWithLayout(c *gin.Context, contentTemplate string, data interface{}) {
	projectRoot := getProjectRoot()
	
	templateFuncs := map[string]any{
		"contains": strings.Contains,
		"add":      func(a, b int) int { return a + b },
		"sub":      func(a, b int) int { return a - b },
		"mul":      func(a, b int) int { return a * b },
		"div":      func(a, b float64) float64 { return a / b },
	}
	tmpl, err := template.New("base").Funcs(templateFuncs).ParseFiles(
		filepath.Join(projectRoot, "templates/layouts/base.html"),
		filepath.Join(projectRoot, "templates/partials/header.html"),
		filepath.Join(projectRoot, "templates/partials/navigation.html"),
		filepath.Join(projectRoot, fmt.Sprintf("templates/%s.html", contentTemplate)),
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
