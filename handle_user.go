package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// acount render profile of user.
func acount(c echo.Context) error {
	sess, _ := session.Get("session", c)
    data := make(map[string]interface{}, 2)
	data["name"] = sess.Values["name"]
    data["id"] = sess.Values["userid"]

    if data["id"] == nil {
		return c.Redirect(http.StatusSeeOther, "/login") // 303 code
	}
	return c.Render(200, "acount.html", data)
}

// remove this function
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.Render(http.StatusOK, "user.html", id)
}

// updateAcount updates Acount information
func updateAcount(c echo.Context) error {
    data := make(map[string]interface{},1)
    sess, _ := session.Get("session", c)
    
    uid := sess.Values["userid"]
    name := sess.Values["name"]

    data["name"] = name
    
    fmt.Println("name is : ", name, "userid is :", uid)
    //userid , _:= strconv.Atoi(uid.(string))

    if uid == nil {
        // login first
        return c.Redirect(http.StatusSeeOther, "/login") // 303 code 
    }
    data["id"] = uid
 
    return c.Render(200, "upacount.html", data)
}

func mysess(c echo.Context, name string, userid int) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
        MaxAge:   60 * 60, // = 1h,
		HttpOnly: true, // no websocket or any thing else
	}
	sess.Values["name"] = name
	sess.Values["userid"] = userid
	sess.Save(c.Request(), c.Response())
}


func login(c echo.Context) error {
	femail := c.FormValue("email")
	fpass := c.FormValue("password")
    userid,  name, email, pass := getUsername(femail)

	if pass == fpass && femail == email {
		//userSession[email] = name
        mysess(c, name, userid)
		return c.Redirect(http.StatusSeeOther, "/") // 303 code
		// TODO redirect to latest page
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
		//fmt.Println(err)
		return c.Render(200, "sign.html", "wrrone")
	}
	return c.Redirect(http.StatusSeeOther, "/login") // 303 code
}

