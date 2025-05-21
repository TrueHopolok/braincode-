package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/views"
)

const MAX_PROBLEMS_ON_PAGE = 20

func Problemset(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Method=%s is not allowed\nAllowed=GET", r.Method), 405)
		w.Header().Set("Allow", "GET")
		logger.Log.Debug("req=%p method=%s is not allowed", r, r.Method)
		return
	}

	urldata := r.URL.Query()

	ses, isauth := sessionHandler(w, r)

	templ := "index.html"
	lang := urldata.Get("lang")
	if lang == "ru" {
		// templ = "index_ru.html" // TODO(vadim): switch to this when it will be done
		logger.Log.Error("req=%p failed; error=translation is not available yet", r)
		http.Error(w, "Not implemented yet", 503)
		return
	} else if lang != "" && lang != "en" {
		logger.Log.Debug("req=%p failed; error=user trying to select invalid language", r)
		http.Error(w, "Such language selection is not allowed", 406)
		return
	}

	var limit, page int
	limit, err := strconv.Atoi(urldata.Get("limit"))
	if err != nil || limit > MAX_PROBLEMS_ON_PAGE {
		limit = MAX_PROBLEMS_ON_PAGE
	}
	page, err = strconv.Atoi(urldata.Get("page"))
	if err != nil {
		page = 0
	}
	rows, err := models.Problemset(limit, page)

	if err := views.Problemset(w, templ, ses.Name, isauth, rows); err != nil {
		http.Error(w, "Failed to write into the response body", 500)
		logger.Log.Error("req=%p failed; error=%s", r, err)
		return
	}
}
