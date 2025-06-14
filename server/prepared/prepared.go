package prepared

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TrueHopolok/braincode-/judge/ml"
	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/session"
)

var Templates *template.Template

func Init() (err error) {
	pattern := filepath.Join(config.Get().TemplatesPath, "*.html")
	Templates, err = template.ParseGlob(pattern)
	ml.AddHTMLTemplate(Templates, "markleftDoc")
	return
}

type T struct {
	Lang     string
	Auth     bool
	Username string
	IsAdmin  bool
	ErrCode  int // 0 = no error
}

func (t T) Tr(en, ru string) string {
	if t.IsRU() {
		return ru
	}
	return en
}

func (t T) IsRU() bool {
	return strings.ToLower(string(t.Lang)) == "ru"
}

func (t T) IsEN() bool {
	return !t.IsRU()
}

func (t T) TrURL(url string) string {
	if t.IsEN() {
		return url
	}
	return url + "?lang=RU"
}

func (t T) LangNormalized() string {
	if t.IsEN() {
		return "en"
	}
	return "ru"
}

func (t T) Session(s session.Session) T {
	t.Auth = !s.IsZero()
	if t.Auth {
		t.Username = s.Name
	} else {
		t.Username = ""
	}
	return t
}

func (t T) LangBool(isengligh bool) T {
	if isengligh {
		t.Lang = "en"
	} else {
		t.Lang = "ru"
	}
	return t
}

func (t T) AuthBool(auth bool, username string) T {
	if !auth {
		t.Auth = false
		t.Username = ""
	} else {
		t.Auth = true
		t.Username = username
	}
	return t
}

func (t T) SetAdmin(isadmin bool) T {
	t.IsAdmin = isadmin
	return t
}

func (t T) Request(r *http.Request) T {
	s := r.URL.Query().Get("error")
	i, err := strconv.Atoi(s)
	if err != nil {
		return t
	}
	t.ErrCode = i
	return t
}

func (t T) SetErr(errcode int) T {
	t.ErrCode = errcode
	return t
}

func (t T) Context(ctx context.Context) T {
	return t.Session(session.Get(ctx))
}

func (t T) E(code int, en, ru string) string {
	if t.ErrCode != code {
		return ""
	}
	if t.IsEN() {
		return en
	} else {
		return ru
	}
}
