package models

import "time"

type User struct {
	Id    string
	Email string
	Name  string
}

type PageRecord struct {
	Id		string
	Url     string `json:"url"`
	Created time.Time
}

type PreviewData struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type PageData struct {
	PreviewData
	PageRecord
}
