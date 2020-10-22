package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

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
	return c.Render(200, "login.html", "wrone password")
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
	//return c.Render(200, "login.html", "welcome")
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	email := sess.Values["email"]

	return c.Render(http.StatusOK, "home.html", getUserSess(email)) //userSession[email.(string)])
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

	return c.Render(200, "stores.html", userSession[email.(string)])
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

// this map for store and manage user session
var userSession = map[string]string{}

// getUserSess takes user email then returns username to display it
func getUserSess(us interface{}) string {
	if len(userSession) == 0 {
		return ""
	}
	return userSession[us.(string)]
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
