package prepared

import (
	"html/template"
	"path/filepath"

	"github.com/TrueHopolok/braincode-/back-end/config"
)

var Templates *template.Template

func Init() (err error) {
	pattern := filepath.Join(config.Get().TemplatesPath, "*.html")
	Templates, err = template.ParseGlob(pattern)
	return
}
