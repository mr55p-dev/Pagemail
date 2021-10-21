# PageSaver
A basic read it later framework - get the pages delivered straight to your inbox.
Hosted on [Heroku](http://page-saver.herokuapp.com), and built on [FastAPI](https://fastapi.tiangolo.com/).
## Routes to add:

- UPDATE: Change user info
- UPDATE: User preferences
- DELETE: Delete a post for a user

## Roadmap:

- [x] Write unit tests
- [ ] Make the shortcuts display the website title.
- [x] Add a worker which sends emails,
- [ ] Marks pages as seen and unseen.
- [x] Set up a domain, and a light frontend in pure HTML/CSS/JS.
- [ ] Check the ability to set a custom token duration.
- [ ] Add a worker which gets metadata for each page saved and stores it in a table.
- [ ] Add user permissions.

---
# Project restart
To get this show back on the road we are gonna do a couple of things:
- [ ] Finally make the emailing work.
- [ ] Bring the async scheduler back inside this project.
- [ ] Start writing some tests for the API
- [ ] Make the API design more consistent (push to a v1.0) (big redesign)
- [ ] Create some tests for the website; switch to proper hosting to get rid of github and look into react-native.
- [ ] add "sign up with service" options
- [ ] Stop messing about with page summaries, instead just number, date and page title is all we need for now.
