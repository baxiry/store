package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// TODO redirect to latest page after login.

func getOneProd(c echo.Context) error {
	// TODO whish is beter all data of product or jast photo ?
	data := make(map[string]interface{})
	sess, _ := session.Get("session", c)
	data["name"] = sess.Values["name"]
	// User ID from path `users/:id`
	id := c.Param("id") // TODO home or catigory.html ?
	productId, _ := strconv.Atoi(id)

	data["product"], err = getProduct(productId)
	if err != nil {
		fmt.Println("with gitCatigories: ", err)
	}
	fmt.Println("product form handle : ", data["product"])
	return c.Render(http.StatusOK, "product.html", data)
}

// e.GET("/users/:id", getUser) ?

func getProds(c echo.Context) error {

	data := make(map[string]interface{})

	sess, _ := session.Get("session", c)
	catigory := c.Param("catigory") // TODO home or catigory.html ?

	data["name"] = sess.Values["name"]
	data["subCatigories"] = catigories[catigory]
	data["products"], err = getProductes(catigory)
	if err != nil {
		fmt.Println("in gitCatigories: ", err)
	}

	err = c.Render(http.StatusOK, "products.html", data)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	name := sess.Values["name"]
	//fmt.Println("name is : ", name)

	file, err := os.Open("../files")
	if err != nil {
		fmt.Printf("failed opening directory: %s\n", err)
	}
	defer file.Close()

	Photonames := make([]string, 0)
	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, fileName := range list {
		Photonames = append(Photonames, fileName)
	}

	data := make(map[string]interface{}, 3)
	data["name"] = name // from session or from memcach ?
	data["photos"] = Photonames

	return c.Render(http.StatusOK, "home.html", data)
}

func stores(c echo.Context) error {
	sess, _ := session.Get("session", c)

	data := make(map[string]interface{}, 3)
	name := sess.Values["name"]
	data["name"] = name // from session or from memcach ?

	return c.Render(200, "stores.html", data)
}

var catigories = map[string][]string{
	"cars":      []string{"mersides", "volswagn", "shefrole", "ford", "jarary", "jawad"},
	"animals":   []string{"dogs", "sheeps", "elephens", "checkens", "lions"},
	"motors":    []string{"harly", "senteroi", "basher", "hddaf", "mobilite"},
	"mobiles":   []string{"sumsung", "apple", "oppo", "netro", "nokia"},
	"computers": []string{"dell", "toshipa", "samsung", "hwawi", "hamed"},
	"services":  []string{"penter", "developer", "cleaner", "shooter", "gamer"}, //services
	"others":    []string{"somthing", "another-somth", "else", "anythings"},
}

func mysess(c echo.Context, name, email string) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3000, // = 60s * 60 = 1h,
		HttpOnly: true, // no websocket or any thing else
	}
	sess.Values["name"] = name
	sess.Values["email"] = email
	sess.Save(c.Request(), c.Response())
}

// upload photos
func uploadPage(c echo.Context) error {
	data := make(map[string]interface{}, 3)
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println("erro upload session is : ", err)
	}
	email := sess.Values["email"]
	name := sess.Values["name"]
	data["name"] = name
	if email == nil {
		// TODO flash here

		return c.Redirect(http.StatusSeeOther, "/login") // 303 code
	}
	// c.Response().Status
	return c.Render(200, "upload.html", data)
}

func login(c echo.Context) error {
	femail := c.FormValue("email")
	fpass := c.FormValue("password")
	name, email, pass := getUsername(femail)

	if pass == fpass && femail == email {
		//userSession[email] = name
		mysess(c, name, email)
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

func photoFold() string {
	if os.Getenv("USERNAME") == "fedora" {
		fmt.Println("username is : ", "fedora")
		return "/home/fedora/repo/files/"
	}
	fmt.Println("username is : ", "root")
	return "/root/files/"
}

func upload(c echo.Context) error {
	// TODO: how upload this ?.  definde uploader by session
	sess, _ := session.Get("session", c)
	email := sess.Values["email"]
	fmt.Println("email of owner session", email)

	title := c.FormValue("title")
	catigory := c.FormValue("catigory")
	details := c.FormValue("description")
	price, _ := strconv.Atoi(c.FormValue("price"))

	// Read files, Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	//fmt.Println("files is :", files[0].Filename)
	picts := ""
	for _, v := range files {
		fmt.Println(v.Filename)
		picts += v.Filename
		picts += "];["
		fmt.Println(picts)
		// TODO Rename pictures.
	}

	err = insertProduct(email.(string), title, catigory, details, picts, price)
	if err != nil {
		fmt.Println()
		fmt.Println("error in insert product", err)
	}

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		// Destination
		dst, err := os.Create(photoFold() + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()
		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}
	// TODO redirect to home or to acount ??
	return c.Redirect(http.StatusSeeOther, "/") // 303 code
}

func signPage(c echo.Context) error {
	return c.Render(200, "sign.html", "hello")
}

func loginPage(c echo.Context) error {

	return c.Render(200, "login.html", "hello")
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.Render(http.StatusOK, "user.html", id)
}

func acount(c echo.Context) error {
	//name := c.Param("name")
	//fmt.Println("param is : ", name)
	sess, _ := session.Get("session", c)
	data := make(map[string]interface{}, 3)
	data["name"] = sess.Values["name"]

	if data["name"] == nil {
		return c.Redirect(http.StatusSeeOther, "/login") // 303 code
	}
	return c.Render(200, "acount.html", data)
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
