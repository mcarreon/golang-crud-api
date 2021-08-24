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
