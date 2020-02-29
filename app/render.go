package app

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/FoodyTech/renche/apputil"
	"github.com/FoodyTech/renche/render"

	"github.com/go-chi/chi"
)

func (app *App) renderHandler(w http.ResponseWriter, r *http.Request) {
	format := chi.URLParam(r, "format")

	switch format {
	case "json":
		app.renderJSON(w, r)
	case "html":
		app.renderHTML(w, r)
	case "zip":
		app.renderZip(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (app *App) renderJSON(w http.ResponseWriter, r *http.Request) {
	res, errRender := app.render(r.Context(), r)
	if errRender != nil {
		log.Print(errRender)
		return
	}
	apputil.WriteJSONResponse(w, 200, res.Page.Content)
}

func (app *App) renderHTML(w http.ResponseWriter, r *http.Request) {
	res, errRender := app.render(r.Context(), r)
	if errRender != nil {
		log.Print(errRender)
		return
	}
	apputil.WriteHTMLResponse(w, 200, res.Page.Content)
}

func (app *App) renderZip(w http.ResponseWriter, r *http.Request) {
	res, errRender := app.render(r.Context(), r)
	if errRender != nil {
		log.Print(errRender)
		return
	}
	apputil.WriteZipResponse(w, 200, "page.zip", res.Page.Content)
}

func (app *App) render(ctx context.Context, r *http.Request) (*render.Result, error) {
	params := r.URL.Query()

	res, errRender := render.Render(r.Context(), render.Params{
		URL:    params.Get("url"),
		Idle:   params.Get("idle") == "true",
		Wait:   time.Duration(asInt(params, "wait", 0)),
		Minify: params.Get("minify") == "true",
	})
	return res, errRender
}

func asInt(vals url.Values, param string, def int) int {
	v := vals.Get(param)
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
