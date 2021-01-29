# PageSaver
A basic read it later framework - get the pages delivered straight to your inbox.
Hosted on [Heroku](http://page-saver.herokuapp.com), and built on [FastAPI](https://fastapi.tiangolo.com/).
## Routes to add:

- DELETE: Delete a user
- UPDATE: Change user info
- UPDATE: User preferences
- DELETE: Delete a post for a user

## Roadmap:

- [ ] Write unit tests
- [ ] Make the shortcuts display the website title.
- [ ] Add a worker which sends emails, and marks pages as sent and unsent.
- [ ] Set up a domain, and a light frontend in pure HTML/CSS/JS.
- [ ] Check the ability to set a custom token duration.
- [ ] Add a worker which gets metadata for each page saved and stores it in a table.
- [ ] Add user permissions.