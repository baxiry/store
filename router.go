package main

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// path file is depends to enveronment.
func templ() *Template {
	var p string
	if os.Getenv("USERNAME") != "fedor" {
		p = "/root/store/"
	}
	files := []string{
        p + "tmpl/home.html", p + "tmpl/acount.html", p + "tmpl/login.html", p + "tmpl/sign.html", p + "tmpl/stores.html", p + "tmpl/mystore.html", p + "tmpl/notfound.html",
		p + "tmpl/upload.html", p + "tmpl/product.html", p + "tmpl/products.html", p + "tmpl/part/header.html", p + "tmpl/part/footer.html",
	}
	return &Template{templates: template.Must(template.ParseFiles(files...))}
}


/*
//go:embed tmpl/*
var filesEmbed embed.FS

func templ() *Template {
	files := []string{
		"tmpl/home.html", "tmpl/acount.html", "tmpl/login.html", "tmpl/sign.html", "tmpl/stores.html", "tmpl/mystore.html",
        "tmpl/upload.html", "tmpl/product.html", "tmpl/products.html", "tmpl/part/header.html", "tmpl/part/footer.html",
	}

    return &Template{templates: template.Must(template.ParseFS(filesEmbed, files...))}
}

//go:embed assets/*
var content embed.FS

var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))
var contentRewrite = middleware.Rewrite(map[string]string{"/": "/static/$1"})
var e = echo.New()
func SetupRoutes() {
    e.GET("/*", contentHandler, contentRewrite)
}

*/

//t, err := template.ParseFS(assetData, "tmpl/")
//if err != nil {
//		fmt.Println(err)
//}



