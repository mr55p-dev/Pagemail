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

type Client struct {
	jobs    chan string
	queries *dbqueries.Queries
}

// New returns a new [Client] and starts waiting for jobs
func New(ctx context.Context, queries *dbqueries.Queries) *Client {
	client := &Client{
		jobs:    make(chan string, 1),
		queries: queries,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case pageID := <-client.jobs:
				page, err := queries.ReadPageById(ctx, pageID)
				if err != nil {
					// TODO: handle error
					continue
				}
				go client.fetch(ctx, &page)
			}
		}
	}()
	return client
}

// Queue adds a pageID to the queue of previews
func (c *Client) Queue(pageID string) {
	c.jobs <- pageID
}

// Sweep checks for any pages marked unknown and attempts to generate a preview
func (c *Client) Sweep(ctx context.Context) {
	pageIDs, err := c.queries.ReadPageIdsByPreviewState(ctx, "unknown")
	if err != nil {
		panic(err)
	}
	for _, ID := range pageIDs {
		c.Queue(ID)
	}
}

func (c *Client) fetch(ctx context.Context, page *dbqueries.Page) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := fetchPreview(ctx, page)
	pageUpdate := dbqueries.UpdatePagePreviewParams{
		Title:       page.Title,
		Description: page.Description,
		ImageUrl:    page.ImageUrl,
		Updated:     time.Now(),
		ID:          page.ID,
	}
	if err == nil {
		pageUpdate.PreviewState = "success"
	} else {
		pageUpdate.PreviewState = "error"
	}

	err = c.queries.UpdatePagePreview(ctx, pageUpdate)
	if err != nil {
		return
	}
}

func fetchUrl(url string) ([]byte, error) {
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

func fetchMeta(contents []byte) (*DocumentMeta, error) {
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

func fetchPreview(ctx context.Context, page *dbqueries.Page) error {
	now := time.Now()
	page.Updated = now
	page.ReadabilityStatus.String = "unknown"
	page.ReadabilityStatus.Valid = true

	// fetch the document
	content, err := fetchUrl(page.Url)
	if err != nil {
		log.Printf("fetch error, %s", err)
		return err
	}

	// async check readability and fetch the doc info
	previewChan := make(chan DocumentMeta)
	errorChan := make(chan error)
	isReadableChan := make(chan bool)
	go func() {
		out, err := fetchMeta(content)
		previewChan <- *out
		errorChan <- err
	}()

	go func() {
		isReadableChan <- checkIsReadable(ctx, page.Url, content)
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

func checkIsReadable(ctx context.Context, url string, content []byte) bool {
	return true
}
