package handlers

import (
	logs "gobp/pkg/logger"
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func IndexHandler(c *gin.Context) {
	ctx := c.Request.Context()
	l := logs.GetLoggerctx(ctx)
	l.Info("this is a test info")
	tmpl, err := template.ParseFiles(filepath.Join(viper.GetString("app.uiTemplates"), "index.html"))
	if err != nil {
		l.Sugar().Errorf("parse template failed", err)
		RenderErrorTemplate(c, "Failed to parse form", err)
		return
	}

	// Execute the template and write the output to the response
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		l.Sugar().Errorf("execute template failed", err)
		return
	}
}
