# golang-crud-api

Golang/REST Code Challenge
Using Go as your language, create a CRUD API to manage a list of Books, fulfilling the
following requirements:

1. Books should have the following Attributes:
  - Title
  - Author
  - Publisher
  - Publish Date
  - Rating (1-3)
  - Status (CheckedIn, CheckedOut)
2. Each endpoint should have test coverage of both successful and failed (due to user error)
requests.
3. Use a data store of your choice.
4. Include unit tests that exercise the endpoints such that both 200-level and 400-level
responses are induced.
5. The app should be stood up in Docker, and client code (such as cURL requests and your unit
tests) should be executed on the host machine, against the containerized app.
6. Send the project along as a git repository.
7. Please do not use go-swagger to generate the server side code. Part of the goal of this challenge is to see your coding skills. :)

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
