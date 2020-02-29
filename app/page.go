package app

import (
	"context"
	"log"
	"net/http"

	"github.com/FoodyTech/renche/apputil"

	"github.com/go-chi/chi"
)

type RenderResult struct {
	Info    string
	URLs    []string
	Success int
	Failure int
}

func (app *App) pageStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, errGet := app.cache.Get(context.TODO(), id)
	if errGet != nil {
		return
	}
	// TODO
	_ = res

	resp := struct{}{}
	apputil.WriteJSONResponse(w, 200, resp)
}

func (app *App) pageResultHandler(w http.ResponseWriter, r *http.Request) {
	format := chi.URLParam(r, "format")

	switch format {
	case "json":
		app.resultJSON(w, r)
	case "zip":
		app.resultZip(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (app *App) resultJSON(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_ = id

	res, errRender := app.render(r.Context(), r)
	if errRender != nil {
		log.Print(errRender)
		return
	}
	apputil.WriteZipResponse(w, 200, "page.zip", res.Page.Content)
}

func (app *App) resultZip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_ = id

	res, errRender := app.render(r.Context(), r)
	if errRender != nil {
		log.Print(errRender)
		return
	}
	apputil.WriteZipResponse(w, 200, "page.zip", res.Page.Content)
}
