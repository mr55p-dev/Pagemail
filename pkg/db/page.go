package db

import (
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type Page struct {
	Id                  string     `db:"id"`
	UserId              string     `db:"user_id"`
	Url                 string     `db:"url"`
	Title               *string    `db:"title"`
	Description         *string    `db:"description"`
	ImageUrl            *string    `db:"image_url"`
	ReadabilityStatus   *string    `db:"readability_status"`
	ReadabilityTaskData *string    `db:"readability_task_data"`
	IsReadable          *bool      `db:"is_readable"`
	Created             *time.Time `db:"created"`
	Updated             *time.Time `db:"updated"`
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
