package preview

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/logging"

	"io"
	"net/http"
)

var logger = logging.NewLogger("preview")

type Client struct {
	jobs chan string
	db   *sql.DB
}

// New returns a new [Client] and starts waiting for jobs
func New(ctx context.Context, db *sql.DB) *Client {
	queries := queries.New(db)
	client := &Client{
		jobs: make(chan string, 1),
		db:   db,
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

func (c *Client) fetch(ctx context.Context, page *queries.Page) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	q := queries.New(c.db)
	state := "success"
	if err := fetchPreview(ctx, page); err != nil {
		state = "error"
	}
	err := q.UpdatePagePreview(ctx, queries.UpdatePagePreviewParams{
		Title:        page.Title,
		Description:  page.Description,
		ImageUrl:     page.ImageUrl,
		ID:           page.ID,
		PreviewState: state,
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Failed to update page preview")
	}

	isReadable := isPageReadable(page)
	err = q.UpdatePageReadability(ctx, queries.UpdatePageReadabilityParams{
		Readable: isReadable,
		ID:       page.ID,
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Failed to update page readability")
	}
	return
}

func isPageReadable(page *queries.Page) bool {
	// TODO: implement
	return true
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

func fetchPreview(ctx context.Context, page *queries.Page) error {
	now := time.Now()
	page.Updated = now

	// fetch the document
	content, err := fetchUrl(page.Url)
	if err != nil {
		log.Printf("fetch error, %s", err)
		return err
	}

	previewData, err := fetchMeta(content)
	if err != nil {
		return err
	}

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

	return nil
}
