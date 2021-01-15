package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"os"
)

func assets() string {
	if os.Getenv("USERNAME") != "fedora" {
		return "/root/store/assets"
	}
	return "assets"
}

func main() {

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Renderer = templ()

	e.Static("/a", assets())
	e.Static("/fs", photoFold())

	e.GET("/", home)
	e.GET("/sign", signPage)
	e.POST("/sign", signup)
	e.GET("/login", loginPage)
	e.POST("/login", login)
	e.GET("/stores", stores)
	e.GET("/acount/:name", acount)
	e.GET("/:catigory", getProds)
	e.GET("/:catigory/:id", getOneProd)
	e.GET("/upload", uploadPage)
	e.POST("/upload", upload)
	// e.GET("/:user", getUser)

	db := setdb()
	defer db.Close()

	e.Logger.Fatal(e.Start(":8080"))
}
