package app

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/FoodyTech/renche/apputil"

	"github.com/FoodyTech/sitemap"
)

// sitemapHandler POST /sitemap/xml/url
//
func (app *App) sitemapHandler(w http.ResponseWriter, r *http.Request) {
	response, errFetch := (&http.Client{Timeout: 5 * time.Second}).Get("url")
	if errFetch != nil {
		return
	}
	defer response.Body.Close()

	m, errParse := sitemap.Parse(response.Body)
	if errParse != nil {
		return
	}

	limit := 1000
	urls := getURLs(m, limit)
	_ = urls

	var resp = struct {
		Success bool
		Errors  []error
	}{
		Success: true,
		Errors:  []error{},
	}
	apputil.WriteJSONResponse(w, 200, resp)
}

// sitemapFileHandler POST /sitemap/xml/file
//
func (app *App) sitemapFileHandler(w http.ResponseWriter, r *http.Request) {
	var limit = 1000
	var workers = 1
	var delay = 0
	var idle = false
	var wait = 0
	var minify = false

	m, errParse := sitemap.Parse(r.Body)
	if errParse != nil {
		return
	}

	urls := getURLs(m, limit)

	result := struct {
		Info struct {
			ID             string
			SitemapXMLFile string //: file.filename || file.fieldname,
			Limit          int
			Workers        int
			Delay          time.Duration
		}
		URLs    []string
		Results []string
	}{
		URLs: urls,
	}

	id := generateUID()
	errSet := app.cache.Set(context.TODO(), id, result)
	if errSet != nil {
		return
	}
	app.renderCache(id, urls, workers, delay, idle, wait, minify)
}

func (app *App) renderCache(id string, urls []string, workers, delay, idle, wait, minify interface{}) {
	// TODO
}

func getURLs(s *sitemap.Sitemap, limit int) []string {
	var urls []string
	for _, u := range s.URL {
		if len(urls) >= limit {
			break
		}
		if strings.HasPrefix(u.Location, "http://") ||
			strings.HasPrefix(u.Location, "https://") {
			urls = append(urls, u.Location)
		}
	}
	return urls
}

func generateUID() string {
	return fmt.Sprintf("page-%d-%d", time.Now().Unix(), rand.Intn(1<<20))
}
