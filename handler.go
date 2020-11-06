package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func uploadPage(c echo.Context) error {
	// c.Response().Status
	return c.Render(200, "upload.html", nil)
}

func upload(c echo.Context) error {
	// Read form fields
	//name := c.FormValue("name")
	//email := c.FormValue("email")

	// Read files
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("../files/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	return c.Render(http.StatusOK, "home.html", nil) //fmt.Sprintf("<p>Uploaded successfully %d files with fields name=%s and email=%s.</p>",len(files), name, email))
}

// this map for store and manage user session
var userSession = map[interface{}]string{}

func mysess(c echo.Context, email string) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // * 7,
		HttpOnly: true,
	}
	sess.Values["email"] = email
	sess.Save(c.Request(), c.Response())
}

func login(c echo.Context) error {
	femail := c.FormValue("email")
	fpass := c.FormValue("password")
	name, email, pass := getUsername(femail)

	if pass == fpass && femail == email {
		userSession[email] = name
		mysess(c, email)
		return c.Redirect(http.StatusSeeOther, "/") // 303 code
	}
	return c.Render(200, "login.html", "Username or password is wrong")
}

func signup(c echo.Context) error {
	name := c.FormValue("username")
	pass := c.FormValue("password")
	email := c.FormValue("email")
	phon := c.FormValue("phon")
	err := insertUser(name, pass, email, phon)
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "sign.html", "wrrone")
	}
	return c.Redirect(http.StatusSeeOther, "/login") // 303 code
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	email := sess.Values["email"]
	fmt.Println(email)

	file, err := os.Open("../files")
	if err != nil {
		fmt.Printf("failed opening directory: %s\n", err)
	}
	defer file.Close()

	Photonames := make([]string, 0)
	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		Photonames = append(Photonames, name)
	}

	data := make(map[string]interface{}, 2)
	data["name"] = userSession[email]
	data["photo"] = Photonames

	return c.Render(http.StatusOK, "home.html", data) //, userSession[email]) //getUserSess(email))
}

func signPage(c echo.Context) error {
	return c.Render(200, "sign.html", nil)
}

func loginPage(c echo.Context) error {
	return c.Render(200, "login.html", "hello")
}

func stores(c echo.Context) error {
	sess, _ := session.Get("session", c)
	email := sess.Values["email"]

	return c.Render(200, "stores.html", userSession[email])
}

func acount(c echo.Context) error {
	return c.Render(200, "acount.html", "") //[]string{name, "my acount"})
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.Render(http.StatusOK, "user.html", id)
}

/* Cookies

func writeCookie(c echo.Context, email string) error {
	cookie := new(http.Cookie)
	cookie.Name = "email"
	cookie.Value = email
	cookie.Expires = time.Now().Add(1 * time.Minute)
	c.SetCookie(cookie)
	//return c.String(http.StatusOK, "write a cookie")
	return nil
}

func readCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("username")
	if err != nil {
		return "", err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return cookie.Value, nil //c.String(http.StatusOK, "read a cookie")
}
*/
