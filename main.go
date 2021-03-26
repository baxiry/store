package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//var asset embed.FS
// hella
func assets() string {
	if os.Getenv("USERNAME") != "fedor" {
		return "/root/store/assets"
	}
	return "assets"
}

//go:embed assets/*
var content embed.FS

var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))
var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})
var e = echo.New()
func SetupRoutes() {
   e.GET("/*", contentHandler, contentRewrite)
}

func main() {

	db := setdb()
	defer db.Close()
    SetupRoutes()
    //e.GET("/a", contentHandler, contentRewrite)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Renderer = templ()

    e.Static("/a", assets())
	e.Static("/fs", photoFold())

	e.GET("/", home)
	e.GET("/sign", signPage)
	e.POST("/sign", signup)
	e.GET("/login", loginPage)
	e.POST("/login", login)
	//e.GET("/stores", stores)
	e.GET("/acount/:name", acount)
	e.GET("/:catigory", getProds)
	//e.GET("/product/:id", getOneProd)
	e.GET("/:catigory/:id", getOneProd) // whech is beter ? :catigory or /product ?
	e.GET("/upload", uploadPage)
	e.POST("/upload", upload)
	// e.GET("/:user", getUser)

	e.Logger.Fatal(e.Start(":8080"))
}
