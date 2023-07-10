package models

import "time"

type User struct {
	Id    string
	Email string
	Name  string
}

type PageRecord struct {
	Url     string `json:"url"`
	Created time.Time
}

type UrlPreviewData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UrlData struct {
	UrlPreviewData
	PageRecord
}
