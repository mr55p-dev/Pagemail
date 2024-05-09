package db

import (
	"context"
	"time"
)

const PAGE_SIZE int = 10

func (client *Client) CreatePage(ctx context.Context, p *Page) error {
	_, err := client.db.ExecContext(ctx, `
		INSERT INTO pages (id, user_id, url, created, updated)
		VALUES (?, ?, ?, ?, ?)
	`, p.Id, p.UserId, p.Url, p.Created, p.Updated)
	return err
}

func (client *Client) ReadPage(ctx context.Context, id string) (*Page, error) {
	out := new(Page)
	err := client.db.GetContext(ctx, out, "SELECT * FROM pages WHERE id = ?", id)
	return out, err
}

func (client *Client) DeletePage(c context.Context, id string) error {
	_, err := client.db.Exec(`DELETE FROM pages WHERE id = ?`, id)
	return err
}

func (client *Client) ReadPagesByUserId(c context.Context, id string, page int) ([]Page, error) {
	var skip, limit int
	if page == 0 {
		page = 1
	}
	if page > 0 {
		skip = PAGE_SIZE * (page - 1)
		limit = PAGE_SIZE
	} else if page <= 0 {
		skip = 0
		limit = 9999
	}
	var pages []Page
	err := client.db.SelectContext(
		c, &pages,
		`SELECT * FROM pages WHERE user_id = ? ORDER BY created DESC LIMIT ? OFFSET ?`,
		id, limit, skip,
	)
	if err != nil {
		return nil, err
	}
	return pages, err
}

func (client *Client) ReadPagesByUserBetween(c context.Context, id string, start, end time.Time) ([]Page, error) {
	pages := make([]Page, 0)
	err := client.db.SelectContext(c, &pages, `SELECT * FROM pages WHERE user_id = ? AND created BETWEEN ? AND ? ORDER BY created DESC`, id, start, end)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (client *Client) UpsertPage(c context.Context, p *Page) error {
	_, err := client.db.Exec(
		`INSERT OR REPLACE INTO pages VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Id, p.UserId, p.Url, p.Title,
		p.Description, p.ImageUrl, p.ReadabilityStatus,
		p.ReadabilityTaskData, p.IsReadable, p.Created, p.Updated,
	)
	return err
}

func (client *Client) DeletePagesByUserId(c context.Context, id string) (int, error) {
	res, err := client.db.Exec(`
		DELETE FROM pages WHERE user_id = ?
	`, id)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), err
}
