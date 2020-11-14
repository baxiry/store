package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	db := setdb()
	defer db.Close()

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Renderer = templ()

	e.Static("/a", "assets")
	e.Static("acount/a", "assets")
	e.Static("acount/files", "../files")
	e.Static("files", "../files")

	e.GET("/", home)
	e.GET("/sign", signPage)
	e.POST("/sign", signup)
	e.GET("/login", loginPage)
	e.GET("/stores", stores)
	e.GET("/acount/:name", acount)
	e.POST("/login", login)
	e.GET("/:user", getUser)

	e.GET("/upload", uploadPage)
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":8080"))
}
