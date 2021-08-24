# golang-crud-api

## Quick Startup

```shell
docker-compose up --build

```

### Quick Info

Made with Docker, Golang, and Postgresql

## Routes

### GET

There are 2 get routes:

```text

/books 
/book/[title]
```

/books returns all books in the DB

/book/[title] returns a book with a matching title

### POST

```text

/books
```

example body params:

```json
  {
    "title": "title",
    "author": "author",
    "publishDate": "2009-01-20",
    "publisher": "publisher",
    "rating": 3,
    "status": "CheckedIn"
  }
```

### PUT

```text

/book/[title]
```

/book/[title] updates a book with matching title using supplied body params. Any number of params may be entered matching the POST example.

### DELETE

```text

/book/[title]
```

/book/[title] deletes a book matching the supplied title.

## Write-Up

Within the first four hours, I was able to get the routes created, however that was without database interaction. Database interaction was added an hour or so after that initial 4 hours. I had a lot of issues getting Docker to work on windows, so I saved dockerizing the app for last.

The app makes an assumption that duplicate books are allowed.

I was also unable to do any package splitting(?). For some reason, when I created subfolders with different packages, I could not get the importing to work, so for the sake of time all the files are under the main package. My ideal structure would be something like this:

```text
+-- _server
|   +-- server.go
|   +-- server_test.go
|   +-- server_integration_test.go
+-- _utils
|   +-- utils.go
|   +-- asserts.go
+-- _db
|   +-- database.go
|   +-- controller.go
+-- _models
|   +-- models.go
+-- main.go
```

Tests on the route functionality run when the container is built. This is due to being unable to cli into the container for some reason, so its set where the build fails if the unit tests fail. Integration tests are included and can be run locally, but because I could not cli into the container, I'm unable to run the integration tests manually. I'm sure the is a way to have the tests run after the database has been finished initializing, but I'm unsure how to write that functionality with docker.

### Improvements

A couple things I would implement would be more descriptive error responses - currently I'm only sending back the status code.

Possibly adding a UUID lookup to allow for more accurate searches on books with duplicate titles but different meta data.

Getting the integration tests to run automatically after the postgres database is initialized and populated.
