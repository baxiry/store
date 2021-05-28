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

var currentPage string

// TODO url := c.Request().URL  we need change url path ? example /cats/ to /cats

// costorm error handler
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
    fmt.Println(err)
    //c.Redirect(303, "notfound.html")
    c.Redirect(http.StatusSeeOther, "/notfound") // 303 code
    return
}




// perhaps is beter ignoring this feater ??!
func myStores(c echo.Context) error { // TODO rename to myproduct ??
	sess, _ := session.Get("session", c)
    data := make(map[string]interface{}, 2)
	name := sess.Values["name"]
    email := sess.Values["email"]
	data["name"] = name // from session or from memcach ?
    data["email"] = email // from session or from memcach ?


    data["products"] = myProducts(email.(string))
    fmt.Println("data products is : ", data["products"])

    err =  c.Render(200, "mystore.html", data)
    if err != nil {
        fmt.Println("at myStores function ", err)
        return nil
    } 

    return nil// c.Render(200, "mystore.html", data)
}


func home(c echo.Context) error {

	sess, _ := session.Get("session", c)
	name := sess.Values["name"]
	//fmt.Println("name is : ", name)

	data := make(map[string]interface{}, 3)
	data["name"] = name
	data["catigories"] = catigories
    return c.Render(http.StatusOK, "home.html", data)
}

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
	//fmt.Println("product form handle : ", data["product"])
    return c.Render(http.StatusOK, "product.html", data)
}

// getProduct get all data of one product from db, and reder it
func getProds(c echo.Context) error {
	data := make(map[string]interface{})

	sess, _ := session.Get("session", c)
	catigory := c.Param("catigory") // TODO home or catigory.html ?

	data["name"] = sess.Values["name"]
	data["subCatigories"] = catigories[catigory]
	data["products"], err = getProductes(catigory)
    fmt.Println(data)
	
    // TODO : handle or ignore this error ?
	//if err != nil {
	//	fmt.Println("in gitCatigories: ", err)
    //}

	return c.Render(http.StatusOK, "products.html", data)
}

var catigories = map[string][]string{
	"cars":      {"mersides", "volswagn", "shefrole", "ford", "jarary", "jawad"},
	"animals":   {"dogs", "sheeps", "elephens", "checkens", "lions"},
	"motors":    {"harly", "senteroi", "basher", "hddaf", "mobilite"},
	"mobiles":   {"sumsung", "apple", "oppo", "netro", "nokia"},
	"computers": {"dell", "toshipa", "samsung", "hwawi", "hamed"},
	"services":  {"penter", "developer", "cleaner", "shooter", "gamer"}, //services
	"others":    {"somthing", "another-somth", "else", "anythings"},
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
		fmt.Println(v.Filename)
		picts += v.Filename
		picts += "];["
		fmt.Println(picts)
		// TODO Rename pictures.
	}

	err = insertProduct(email.(string), title, catigory, details, picts, price)
	if err != nil {
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
	//if err != nil {
	//	fmt.Println("redirect err", err)
	//	return nil
	//}
    //return nil
}

func signPage(c echo.Context) error {
	return c.Render(200, "sign.html", "hello")
}

func loginPage(c echo.Context) error {
	return c.Render(200, "login.html", "hello")
}

// notFoundPage
func notFoundPage(c echo.Context) error {
    return c.Render(200, "notfound.html", nil)
}

// acount render profile of user.
func acount(c echo.Context) error {
	sess, _ := session.Get("session", c)
	data := make(map[string]interface{}, 3)
	data["name"] = sess.Values["name"]

	if data["name"] == nil {
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

// perhaps is beter ignoring this feater ??!
func stores(c echo.Context) error {
	sess, _ := session.Get("session", c)
	data := make(map[string]interface{}, 3)
	name := sess.Values["name"]
	data["name"] = name // from session or from memcach ?
	return c.Render(200, "stores.html", data)
}

// folder when photos is stored.
func photoFold() string {
	if os.Getenv("USERNAME") == "fedor" {
		return "/home/fedor/repo/files/"
	}
	return "/root/files/"
}
