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
    files := []string{
        "tmpl/home.html", "tmpl/acount.html","tmpl/login.html","tmpl/sign.html","tmpl/stores.html","tmpl/mystore.html",
        "tmpl/upload.html","tmpl/product.html","tmpl/products.html","tmpl/part/header.html","tmpl/part/footer.html"}
    return &Template{templates: template.Must(template.ParseFiles(files...))}
}

