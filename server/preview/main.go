package preview

import (
	"bytes"
	"fmt"
	"log"
	"pagemail/server/models"
	"pagemail/server/net"
	"time"

	"pagemail/server/readability"

	"github.com/PuerkitoBio/goquery"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/hook"

	"github.com/pocketbase/pocketbase/core"
)

type DocumentMeta struct {
	Title       string
	Description string
	ImageUrl    string
}

func FetchDocumentMeta(contents []byte) (*DocumentMeta, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}

	page_data := new(DocumentMeta)
	page_data.Title = doc.Find("meta[property='og:title']").AttrOr("content", "")
	page_data.Description = doc.Find("meta[property='og:description']").AttrOr("content", "")

	//Fallbacks
	if page_data.Title == "" {
		page_data.Title = doc.Find("title").Text()
	}
	if page_data.Description == "" {
		doc.Find("meta[name='description']").Each(func(i int, sel *goquery.Selection) {
			if desc, exists := sel.Attr("content"); exists && desc != "" {
				page_data.Description = desc
			}
		})
	}

	return page_data, nil
}

func FetchPreview(url string, cfg readability.ReaderConfig) (*models.Page, error) {
	out := &models.Page{
		LastCrawled:       time.Now(),
		ReadabilityStatus: models.ReadabilityUnknown,
	}
	content, err := net.FetchUrlContents(url)
	if err != nil {
		log.Printf("fetch error, %s", err)
		return out, err
	}

	preview_ch := make(chan DocumentMeta)
	error_ch := make(chan error)
	is_readable_ch := make(chan bool)
	go func() {
		out, err := FetchDocumentMeta(content)
		preview_ch <- *out
		error_ch <- err
	}()

	go func() {
		is_readable_ch <- readability.CheckIsReadable(cfg, url, content)
	}()

	out.IsReadable = <-is_readable_ch

	select {
	case err := <-error_ch:
		log.Printf("Error in goroutine, %s", err)
		return out, err
	case preview_data := <-preview_ch:
		out.Title = preview_data.Title
		out.Description = preview_data.Description
		out.ImageUrl = preview_data.ImageUrl
	}

	return out, nil
}

func PagePreviewHook(app *pocketbase.PocketBase, cfg readability.ReaderConfig) hook.Handler[*core.RecordCreateEvent] {
	return func(e *core.RecordCreateEvent) error {
		url := e.Record.GetString("url")
		if url == "" {
			return fmt.Errorf("url parameter not present in request header")
		}

		go func() {
			res, err := FetchPreview(url, cfg)
			if err != nil {
				log.Printf("Failed to fetch preview, %s", err)
				return
			}

			form := forms.NewRecordUpsert(app, e.Record)
			err = form.LoadData(res.ToMap())
			if err != nil {
				log.Printf("Failed to prepare page preview record update for url %s: %s", url, err)
				return
			}
			err = form.Submit()
			if err != nil {
				log.Printf("Failed to commit page preview record update for url %s: %s", url, err)
				return
			}
		}()

		return nil
	}
}
