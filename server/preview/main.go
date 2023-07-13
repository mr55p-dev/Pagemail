package preview

import (
	"net/http"
	"pagemail/server/models"

	"github.com/PuerkitoBio/goquery"
)

func FetchPreview(url string) (*models.PreviewData, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find("meta[property='og:title']").AttrOr("content", "")
	description := doc.Find("meta[property='og:description']").AttrOr("content", "")

	//Fallbacks
	if title == "" {
		title = doc.Find("title").Text()
	}
	if description == "" {
		doc.Find("meta[name='description']").Each(func(i int, sel *goquery.Selection) {
			if desc, exists := sel.Attr("content"); exists && desc != "" {
				description = desc
			}
		})
	}

	data := &models.PreviewData{
		Title:       title,
		Description: description,
	}

	return data, nil
}
