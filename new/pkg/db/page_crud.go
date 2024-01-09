package db

import "context"

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
	_, err := client.db.Exec(`
		DELETE FROM pages WHERE id = ?
	`, id)
	return err
}

func (client *Client) ReadPagesByUserId(c context.Context, id string) ([]Page, error) {
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

func (client *Client) UpsertPage(c context.Context, p *Page) error {
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
