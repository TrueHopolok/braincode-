package prepared

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/TrueHopolok/braincode-/server/config"
)

var Templates *template.Template

func Init() (err error) {
	pattern := filepath.Join(config.Get().TemplatesPath, "*.html")
	Templates, err = template.ParseGlob(pattern)
	return
}

type T struct {
	Lang string
	Auth bool
}

func (l T) Tr(en, ru string) string {
	if l.IsRU() {
		return ru
	}
	return en
}

func (l T) IsRU() bool {
	return strings.ToLower(string(l.Lang)) == "ru"
}

func (l T) IsEN() bool {
	return !l.IsRU()
}

func (l T) TrURL(url string) string {
	if l.IsEN() {
		return url
	}
	return url + "?lang=RU"
}

func (l T) LangNormalized() string {
	if l.IsEN() {
		return "en"
	}
	return "ru"
}

func TFromBools(isengligh, isauth bool) T {
	var l string
	if isengligh {
		l = "en"
	} else {
		l = "ru"
	}

	return T{
		Lang: l,
		Auth: isauth,
	}
}
