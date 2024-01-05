package db

func (c *Client) InsertUser(u *User) error {
	c.log.Debug().Msgf("Inserting new user with id %s", u.Id)
	res, err := c.DB().Exec(`
		INSERT INTO users (
			id, username, email, password,
			name, avatar, subscribed, shortcutoken,
			has_readability, created, updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Id, u.Username, u.Email, u.Password, u.Name, u.Avatar,
		u.Subscribed, u.ShortcutToken, u.HasReadability, u.Created, u.Updated,
	)
	c.log.Debug().Msgf("Created new user with result %v+", res)
	return err
}

func (c *Client) GetUserById(id string) (*User, error) {
	c.log.Debug().Msgf("Checking user by id %s", id)
	row := c.DB().QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	user := User{}
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.Subscribed,
		&user.ShortcutToken,
		&user.HasReadability,
		&user.Created,
		&user.Updated,
	)
	c.log.Debug().Msgf("Loaded user %+v", user)
	return &user, err
}
func (c *Client) GetUserByEmail(email string) (*User, error) {
	c.log.Debug().Msgf("Checking user by email %s", email)
	row := c.DB().QueryRow(`SELECT * FROM users WHERE email = ?`, email)
	user := User{}
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Avatar,
		&user.Subscribed,
		&user.ShortcutToken,
		&user.HasReadability,
		&user.Created,
		&user.Updated,
	)
	c.log.Debug().Msgf("Loaded user %+v", user)
	return &user, err
}
