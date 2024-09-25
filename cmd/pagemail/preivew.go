package main

import (
	"io"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	PREVIEW_UNKNOWN string = "unknown"
	PREVIEW_PENDING string = "pending"
	PREVIEW_SUCCESS string = "success"
	PREVIEW_FAILURE string = "failure"
)

type PageData struct {
	Title       string
	Description string
}

func getText(node *html.Node, buf io.Writer) {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.TextNode {
			buf.Write([]byte(n.Data))
		} else {
			getText(n, buf)
		}
	}
}

func NodeText(node *html.Node) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == "content" && attr.Val != "" {
			return attr.Val
		}
	}
	buf := new(strings.Builder)
	getText(node, buf)
	return buf.String()
}

func GetPreview(pageUrl *url.URL) (*PageData, error) {
	res := new(PageData)

	doc, err := htmlquery.LoadURL(pageUrl.String())
	if err != nil {
		return nil, err
	}

	titleQueries := []string{
		"/html/head/meta[@name=\"og:title\"]",
		"/html/head/meta[@name=\"title\"]",
		"/html/head/meta[@name=\"Title\"]",
		"/html/head/title",
	}
	for _, q := range titleQueries {
		title, _ := htmlquery.Query(doc, q)
		if txt := NodeText(title); txt != "" {
			res.Title = txt
			break
		}
	}

	descQueries := []string{
		"/html/head/meta[@name=\"og:description\"]",
		"/html/head/meta[@name=\"description\"]",
		"/html/head/meta[@name=\"Description\"]",
	}
	for _, q := range descQueries {
		desc, _ := htmlquery.Query(doc, q)
		if txt := NodeText(desc); txt != "" {
			res.Description = txt
			break
		}
	}

	return res, nil
}