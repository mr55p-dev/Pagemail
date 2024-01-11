package db

import (
	"log/slog"
	"time"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

type Page struct {
	Id                  string     `db:"id" log:"id"`
	UserId              string     `db:"user_id" log:"user_id"`
	Url                 string     `db:"url" log:"url"`
	Title               *string    `db:"title" log:"title"`
	Description         *string    `db:"description" log:"description"`
	ImageUrl            *string    `db:"image_url" log:"image_url"`
	ReadabilityStatus   *string    `db:"readability_status" log:"readability_status"`
	ReadabilityTaskData *string    `db:"readability_task_data"`
	IsReadable          *bool      `db:"is_readable" log:"is_readable"`
	Created             *time.Time `db:"created" log:"created"`
	Updated             *time.Time `db:"updated" log:"updated"`
}

func (page *Page) LogValue() slog.Value {
	vals := tools.LogValue(page)
	return slog.GroupValue(vals...)
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
