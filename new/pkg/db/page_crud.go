package db

import (
	"github.com/labstack/echo/v4"
)

func (client *Client) CreatePage(c echo.Context, p *Page) error {
	_, err := client.db.Exec(`
		INSERT INTO pages (id, user_id, url, created, updated)
		VALUES (?, ?, ?, ?, ?)
	`, p.Id, p.UserId, p.Url, p.Created, p.Updated)
	return err
}

func (client *Client) ReadPagesByUserId(c echo.Context, id string) ([]Page, error) {
	var pages []Page
	rows, err := client.db.Query(`
		SELECT * FROM pages WHERE user_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		v := Page{}
		rows.Scan(
			&v.Id, &v.UserId, &v.Url,
			&v.Title, &v.Description, &v.ImageUrl,
			&v.ReadabilityStatus, &v.ReadabilityTaskData, &v.IsReadable,
			&v.Created, &v.Updated,
		)
		pages = append(pages, v)
	}

	return pages, err
}

func (client *Client) UpsertPage(c echo.Context, p *Page) error {
	_, err := client.db.Exec(
		`INSERT OR REPLACE INTO pages VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Id, p.UserId, p.Url, p.Title,
		p.Description, p.ImageUrl, p.ReadabilityStatus,
		p.ReadabilityTaskData, p.IsReadable, p.Created, p.Updated,
	)
	if err == nil {
		listener, ok := PageEventMap[p.UserId]
		if ok {
			listener <- Event[Page]{
				Event:  EventType("Update"),
				Record: p,
			}
		} else {
		}
	}
	return err
}

func (client *Client) DeletePagesByUserId(c echo.Context, id string) (int, error) {
	res, err := client.db.Exec(`
		DELETE FROM pages WHERE user_id = ?
	`, id)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), err
}
