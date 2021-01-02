package main

import (
	"io"
    "html/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func templ() *Template {
    //return &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}
    files := []string{
        "tmpl/home.html", "tmpl/acount.html","tmpl/login.html","tmpl/sign.html","tmpl/stores.html","tmpl/mystore.html",
        "tmpl/upload.html","tmpl/product.html","tmpl/products.html","tmpl/partial/header.html","tmpl/partial/footer.html"}
        return &Template{templates: template.Must(template.ParseFiles(files...))}
}
