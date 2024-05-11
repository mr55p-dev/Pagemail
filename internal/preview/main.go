package preview

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"

	"io"
	"net/http"
)

var READSTAT_UNKNOWN = "unknown"

func FetchUrlContents(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

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

func FetchPreview(ctx context.Context, page *dbqueries.Page) error {
	now := time.Now()
	page.Updated = now
	page.ReadabilityStatus.String = READSTAT_UNKNOWN
	page.ReadabilityStatus.Valid = true

	// fetch the document
	content, err := FetchUrlContents(page.Url)
	if err != nil {
		log.Printf("fetch error, %s", err)
		return err
	}

	// async check readability and fetch the doc info
	previewChan := make(chan DocumentMeta)
	errorChan := make(chan error)
	isReadableChan := make(chan bool)
	go func() {
		out, err := FetchDocumentMeta(content)
		previewChan <- *out
		errorChan <- err
	}()

	go func() {
		isReadableChan <- CheckIsReadable(ctx, page.Url, content)
	}()

	isReadable, ok := <-isReadableChan
	if ok {
		page.IsReadable.Valid = true
		page.IsReadable.Bool = isReadable
	}

	select {
	case err := <-errorChan:
		return err
	case previewData := <-previewChan:
		if previewData.Title != "" {
			page.Title.String = previewData.Title
			page.Title.Valid = true
		}
		if previewData.Description != "" {
			page.Description.String = previewData.Description
			page.Description.Valid = true
		}
		if previewData.ImageUrl != "" {
			page.ImageUrl.String = previewData.ImageUrl
			page.ImageUrl.Valid = true
		}
	}

	return nil
}

func CheckIsReadable(ctx context.Context, url string, content []byte) bool {
	return true
}
