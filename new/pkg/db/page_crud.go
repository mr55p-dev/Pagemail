package db

func (c *Client) CreatePage(p *Page) error {
	c.log.Debug().Msgf("Creating new Page %+v", p)
	res, err := c.DB().Exec(`
		INSERT INTO pages (id, user_id, url, created, updated)
		VALUES (?, ?, ?, ?, ?)
	`, p.Id, p.UserId, p.Url, p.Created, p.Updated)
	c.log.Debug().Msgf("Created new page with result %s", res)
	return err
}

func (c *Client) ReadPagesByUserId(id string) ([]Page, error) {
	var pages []Page
	c.log.Debug().Msgf("Reading pages with user id %s", id)
	rows, err := c.DB().Query(`
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
