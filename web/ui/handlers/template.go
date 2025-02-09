package handlers

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

// TemplateHandler manages template parsing and execution
type TemplateHandler struct {
	templates map[string]*template.Template
	mutex     sync.RWMutex
	baseDir   string
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(baseDir string) *TemplateHandler {
	return &TemplateHandler{
		templates: make(map[string]*template.Template),
		baseDir:   baseDir,
	}
}

// LoadTemplates preloads all templates
func (th *TemplateHandler) LoadTemplates() error {
	layouts, err := filepath.Glob(filepath.Join(th.baseDir, "layouts/*.html"))
	if err != nil {
		return err
	}

	pages, err := filepath.Glob(filepath.Join(th.baseDir, "pages/*.html"))
	if err != nil {
		return err
	}

	// Load each page template with the base layout
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.ParseFiles(append(layouts, page)...)
		if err != nil {
			return err
		}

		th.mutex.Lock()
		th.templates[name] = tmpl
		th.mutex.Unlock()
	}

	return nil
}

// Render executes the template and writes to response
func (th *TemplateHandler) Render(c *gin.Context, name string, data interface{}) error {
	th.mutex.RLock()
	tmpl, exists := th.templates[name]
	th.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("template %s not found", name)
	}

	return tmpl.Execute(c.Writer, data)
}
