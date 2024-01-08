package db

func (c *DBDriver) CreatePage(p *Page) error {
	c.log.Debug().Msgf("Creating new Page %+v", p)
	res, err := c.DB().Exec(`
		INSERT INTO pages (id, user_id, url, created, updated)
		VALUES (?, ?, ?, ?, ?)
	`, p.Id, p.UserId, p.Url, p.Created, p.Updated)
	c.log.Debug().Msgf("Created new page with result %s", res)
	return err
}

func (c *DBDriver) ReadPagesByUserId(id string) ([]Page, error) {
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

func (c *DBDriver) UpsertPage(p *Page) error {
	res, err := c.DB().Exec(
		`INSERT OR REPLACE INTO pages VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Id, p.UserId, p.Url, p.Title,
		p.Description, p.ImageUrl, p.ReadabilityStatus,
		p.ReadabilityTaskData, p.IsReadable, p.Created, p.Updated,
	)
	c.log.Debug().Msgf("Insert ran with result %+v", res)
	if err == nil {
		c.log.Info().Msg("Started to fire listener")
		listener, ok := PageEventMap[p.UserId]
		if ok {
			c.log.Info().Msg("Fired listener")
			listener <- Event[Page]{
				Event:  EventType("Update"),
				Record: p,
			}
		} else {
			c.log.Info().Msg("Failed to fire, no open channel")
		}
	} 
	return err
}

func (c *DBDriver) DeletePagesByUserId(id string) (int, error) {
	c.log.Debug().Msgf("Deleting pages with user id %s", id)
	res, err := c.DB().Exec(`
		DELETE FROM pages WHERE user_id = ?
	`, id)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), err
}
