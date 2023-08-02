package preview

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pagemail/server/models"
	"time"

	// "pagemail/server/readability"

	"github.com/PuerkitoBio/goquery"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/hook"

	// pb_models "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/core"
)

func fetchUrlContents(url string) (*[]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

type DocumentMeta struct {
	Title       string
	Description string
	ImageUrl    string
}

func FetchDocumentMeta(contents *[]byte) (*DocumentMeta, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*contents))
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

func FetchPreview(url string) (*models.Page, error) {
	out := &models.Page{
		LastCrawled: time.Now(),
		ReadabilityStatus: models.ReadabilityUnknown,
	}
	content, err := fetchUrlContents(url)
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
		// is_readable_ch <- readability.CheckIsReadable(url, content)
		is_readable_ch <- true
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


func PagePreviewHook(app *pocketbase.PocketBase) hook.Handler[*core.RecordCreateEvent] {
	return func(e *core.RecordCreateEvent) error {
		// Fetches and inserts page metadata
		url := e.Record.GetString("url")
		if url == "" {
			return fmt.Errorf("Failed to fetch URL")
		}
		res, err := FetchPreview(url)
		if err != nil {
			log.Printf("Failed to fetch preview, %s", err)
			return err
		}

		form := forms.NewRecordUpsert(app, e.Record)
		form.LoadData(res.ToMap())
		return form.Submit()
	}
}
