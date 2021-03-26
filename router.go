package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//go:embed tmpl/*
var filesEmbed embed.FS

func templ() *Template {
	files := []string{
		"tmpl/home.html", "tmpl/acount.html", "tmpl/login.html", "tmpl/sign.html", "tmpl/stores.html", "tmpl/mystore.html",
        "tmpl/upload.html", "tmpl/product.html", "tmpl/products.html", "tmpl/part/header.html", "tmpl/part/footer.html",
	}

    return &Template{templates: template.Must(template.ParseFS(filesEmbed, files...))}
}
//t, err := template.ParseFS(assetData, "tmpl/")
//if err != nil {
//		fmt.Println(err)
//}


// path file is depends to enveronment.
/*
func templ() *Template {
	var p string
	if os.Getenv("USERNAME") != "fedor" {
		p = "/root/store/"
	}
	files := []string{
		p + "tmpl/home.html", p + "tmpl/acount.html", p + "tmpl/login.html", p + "tmpl/sign.html", p + "tmpl/stores.html", p + "tmpl/mystore.html",
		p + "tmpl/upload.html", p + "tmpl/product.html", p + "tmpl/products.html", p + "tmpl/part/header.html", p + "tmpl/part/footer.html",
	}
	return &Template{templates: template.Must(template.ParseFiles(files...))}
}
*/


