package db

import (
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type Page struct {
	Id                  string
	UserId              string
	Url                 string
	Title               string
	Description         string
	ImageUrl            string
	ReadabilityStatus   string
	ReadabilityTaskData string
	IsReadable          bool
	Created             *time.Time
	Updated             *time.Time
}

func NewPage(userId, url string) *Page {
	now := time.Now()
	return &Page{
		Id:      tools.GenerateNewId(20),
		UserId:  userId,
		Url:     url,
		Created: &now,
		Updated: &now,
	}
}
