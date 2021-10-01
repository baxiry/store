package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// updateFotos updates photos of products
func updateProdFotos(c echo.Context) error {
	
    pid := c.Param("id") 
    id, err := strconv.Atoi(pid)
    if err != nil {fmt.Println("id error", err)}

    // from her : 
    form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]
	//fmt.Println("files is :", files[0].Filename)
	picts := ""
	for _, v := range files {
		picts += v.Filename
		picts += "];["
		// TODO Rename pictures.
	}

    // databas function
    err = updateProductFotos(picts, id)

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

    return c.Redirect(http.StatusSeeOther, "/mystore")
}

// TODO redirect to latest page after login.

func updateProd(c echo.Context) error {
    // TODO  separate edit photos


	pid := c.Param("id") 
    id, err := strconv.Atoi(pid)
    if err != nil {fmt.Println("id error", err)}

    title := c.FormValue("title")
    catig := c.FormValue("catigory")
    descr := c.FormValue("description")
    price := c.FormValue("price")
    photos := c.FormValue("files")
    
    err = updateProduct(title, catig, descr, price, photos, id)
    if err != nil {
        // TODO send error to client with ajax
        fmt.Println("error when update product: ", err)
        return err
    }
    return c.Redirect(http.StatusSeeOther, "/mystore")
}


// delete product
func deleteProd(c echo.Context) error {
    // TODO we need checkout sesston ?
    
	sess, _ := session.Get("session", c)
    ownerid := sess.Values["userid"]
    if ownerid == nil {
        return c.Redirect(http.StatusSeeOther, "/mystore")
    }

    
    id := c.Param("id")
    fmt.Println("id is ", id)
    i, _ := strconv.Atoi(id)
    err = deleteProducte(i)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    
    // return string to ajax resever 
    return c.String(http.StatusOK, "success!") 
}

// perhaps is beter ignoring this feater ??!
func myStores(c echo.Context) error { // TODO rename to myproduct ??
    fmt.Println("at myStores function ")
	sess, _ := session.Get("session", c)
	name := sess.Values["name"]

    if name == nil {
		return c.Redirect(http.StatusSeeOther, "/login") // 303 code
    } 
    
    data := make(map[string]interface{}, 3)
    userid := sess.Values["userid"]
	data["name"] = name // from session or from memcach ?
    data["userid"] = userid // from session or from memcach ?

    data["products"] = myProducts(userid.(int))
    if err != nil {
        fmt.Println("err in product", err)
    }

    return c.Render(200, "mystore.html", data)
}


// TODO redirect to latest page after login.
func getOneProd(c echo.Context) error {
	
    data := make(map[string]interface{})
    
    sess, _ := session.Get("session", c)
	name := sess.Values["name"]
	userid := sess.Values["userid"]
    
    // User ID from path `users/:id`
	id := c.Param("id") // TODO home or catigory.html ?
	productId, _ := strconv.Atoi(id)

    data["name"] = name
    data["userid"] = userid
	data["product"], err = getProduct(productId)

    if err != nil {
		fmt.Println("with gitCatigories: ", err)
	}
    return c.Render(http.StatusOK, "product.html", data)
}

// getProduct get all data of one product from db, and reder it
func getProds(c echo.Context) error {
	data := make(map[string]interface{})

	sess, _ := session.Get("session", c)
    uid := sess.Values["userid"]
	
    catigory := c.Param("catigory") // TODO home or catigory.html ?

	data["name"] = sess.Values["name"]
    data["userid"] = uid
    data["subCatigories"] =  catigories[catigory] // from router.go
    data["products"], _ = getProductes(catigory)
	
    // TODO : handle or ignore this error ?
	//if err != nil {
	//	fmt.Println("in gitCatigories: ", err)
    //}

	return c.Render(http.StatusOK, "products.html", data)
}

// upload uploads new product
func upload(c echo.Context) error {
	// TODO: how upload this ?.  definde uploader by session
	sess, _ := session.Get("session", c)
    ownerid := sess.Values["userid"]

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
		picts += v.Filename
		picts += "];["
		fmt.Println(picts)
		// TODO Rename pictures.
	}

    err = insertProduct( title, catigory, details, picts, ownerid.(int), price)

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

// perhaps is beter ignoring this feater ??!
func stores(c echo.Context) error {
	sess, _ := session.Get("session", c)
    uid := sess.Values["userid"]
    data := make(map[string]interface{}, 2)
	name := sess.Values["name"]

	data["name"] = name // from session or from memcach ?
    data["userid"] = uid
    return c.Render(200, "stores.html", data)
}


// TODO url := c.Request().URL  we need change url path ? example /cats/ to /cats
