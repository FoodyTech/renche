package render

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/svg"
)

// Params regarding rendering in Headless Chrome.
type Params struct {
	URL    string
	Idle   bool
	Wait   time.Duration
	Minify bool
}

// Result represents data from Headless Chrome.
type Result struct {
	URL     string `json:"url"`
	Options Params `json:"options"`

	Browser struct {
		UserAgent string `json:"userAgent"`
		Version   string `json:"version"`
	} `json:"browser"`

	Page struct {
		Content  []byte `json:"content"`
		Title    string `json:"title"`
		URL      string `json:"url"`
		Viewport struct {
			X      float64 `json:"x"`
			Y      float64 `json:"y"`
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
			Scale  float64 `json:"scale"`
		} `json:"viewport"`
	} `json:"page"`
}

// Render the page with a given params in Headless Chrome.
func Render(ctxPar context.Context, params Params) (*Result, error) {
	res := &Result{
		URL:     params.URL,
		Options: params,
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoSandbox,
		chromedp.Headless,
		chromedp.UserAgent("SSRBrowser/1.0"),
	}

	// ctxTimeout, timeoutCancel := context.WithTimeout(ctxPar, params.Wait)
	// defer timeoutCancel()

	allocContext, _ := chromedp.NewExecAllocator(ctxPar, opts...)
	ctx, cancel := chromedp.NewContext(allocContext)
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(params.URL),

		// chromedp.WaitReady("html"),
		chromedp.OuterHTML("html", &str),

		chromedp.Title(&res.Page.Title),
		chromedp.Location(&res.Page.URL),

		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, _, ua, jsver, errVer := browser.GetVersion().Do(ctx)
			if errVer != nil {
				return errVer
			}
			res.Browser.Version = jsver
			res.Browser.UserAgent = ua

			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			res.Page.Viewport.X = contentSize.X
			res.Page.Viewport.Y = contentSize.Y
			res.Page.Viewport.Width = contentSize.Width
			res.Page.Viewport.Height = contentSize.Height
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	if params.Minify {
		minified, errMin := minifyHTML([]byte(str))
		if errMin != nil {
			return nil, errMin
		}
		res.Page.Content = minified
	} else {
		res.Page.Content = []byte(str)
	}
	return res, nil
}

func minifyHTML(buf []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	b, err := m.Bytes("text/html", buf)
	if err != nil {
		return nil, err
	}
	return b, nil
}
