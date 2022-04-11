# restapi_example

APIs

```sh

[GIN-debug] GET    /ping                     --> newsplatform/internal/server.SetupServer.func1 (3 handlers)
[GIN-debug] GET    /news                     --> newsplatform/internal/controllers.FindNews (3 handlers)
[GIN-debug] POST   /news                     --> newsplatform/internal/controllers.CreateNews (3 handlers)
[GIN-debug] GET    /news/:keyword            --> newsplatform/internal/controllers.FindNew (3 handlers)
[GIN-debug] PATCH  /news/:id                 --> newsplatform/internal/controllers.UpdateNew (3 handlers)
[GIN-debug] DELETE /news/:id                 --> newsplatform/internal/controllers.DeleteNews (3 handlers)
[GIN-debug] DELETE /newsbyday/:day           --> newsplatform/internal/controllers.DeleteNewsByDay (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :3000

```
