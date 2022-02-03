package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "test-http/docs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Err error

type Post struct {
	//gorm.Model
	Userid int `gorm:"column:user_id"`
	ID     int
	Title  string
	Body   string
}

type Comment struct {
	//gorm.Model
	PostId int
	ID     int
	Name   string
	Email  string
	Body   string
}

// @title Nix Education Trainee Task API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	fmt.Println("Start")

	dsn := "root:weak_password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	Db, Err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if Err != nil {
		fmt.Println("error connection")
	}
	fmt.Println("connection is OK")
	Db.AutoMigrate(&Post{}, &Comment{})

	//start web server
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/posts", getPosts)
	e.GET("/posts/:id", getPostById)
	e.GET("/comments", getComments)
	e.GET("/comments/:id", getCommentById)
	e.GET("/posts/:id/comments", getCommentsByPostId)
	e.POST("/posts/add", addPost)
	e.POST("/posts/:id/comments/add", addComment)
	e.PUT("/posts/:id/edit", editPost)
	e.PUT("/comments/:id/edit", editComment)
	e.DELETE("/posts/:id/delete", deletePost)
	e.DELETE("comments/:id/delete", deleteComment)

	e.Logger.Fatal(e.Start(":8080"))
}

// getPosts godoc
// @Summary Retrieves posts
// @Produce json
// @Produce xml
// @Success 200 {object} Post
// @Router /posts [get]
func getPosts(c echo.Context) error {
	fmt.Println("getPosts")

	var p []Post
	result := Db.Find(&p)

	if result.Error != nil {
		fmt.Println("select from posts error")
	}
	fmt.Println(result.RowsAffected)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, p)
	case "text/xml":
		return c.XML(http.StatusOK, p)
	default:
		return c.JSON(http.StatusOK, p)
	}
}

// getComments godoc
// @Summary Retrieves comments
// @Produce json
// @Produce xml
// @Success 200 {object} Comment
// @Router /comments [get]
func getComments(c echo.Context) error {
	fmt.Println("getComments")

	var cmt []Comment
	result := Db.Find(&cmt)

	if result.Error != nil {
		fmt.Println("select from comments error")
	}
	fmt.Println(result.RowsAffected)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, cmt)
	case "text/xml":
		return c.XML(http.StatusOK, cmt)
	default:
		return c.JSON(http.StatusOK, cmt)
	}
}

// getPostById godoc
// @Summary Retrieves posts based on given ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /posts/{id} [get]
func getPostById(c echo.Context) error {
	fmt.Println("getPostById")
	id, _ := strconv.Atoi(c.Param("id"))
	var p Post
	result := Db.First(&p, id)

	if result.Error != nil {
		fmt.Println("select from posts error")
	}
	fmt.Println(result.RowsAffected)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, p)
	case "text/xml":
		return c.XML(http.StatusOK, p)
	default:
		return c.JSON(http.StatusOK, p)
	}
}

// getCommentById godoc
// @Summary Retrieves comment based on given ID
// @Produce json
// @Produce xml
// @Param id path int true "Comment Id"
// @Success 200 {object} Post
// @Router /comments/{id} [get]
func getCommentById(c echo.Context) error {
	fmt.Println("getCommentById")
	id, _ := strconv.Atoi(c.Param("id"))
	var cmt Comment
	result := Db.First(&cmt, id)

	if result.Error != nil {
		fmt.Println("select from comments error")
	}
	fmt.Println(result.RowsAffected)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, cmt)
	case "text/xml":
		return c.XML(http.StatusOK, cmt)
	default:
		return c.JSON(http.StatusOK, cmt)
	}
}

// getCommentsByPostId godoc
// @Summary Retrieves comments based on post ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Comment
// @Router /posts/{id}/comments [get]
func getCommentsByPostId(c echo.Context) error {
	fmt.Println("getComments by Post id")
	id, _ := strconv.Atoi(c.Param("id"))
	var cmt []Comment
	result := Db.Where("post_id = ?", id).Find(&cmt)

	if result.Error != nil {
		fmt.Println("select from comments error")
	}
	fmt.Println(result.RowsAffected)
	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, cmt)
	case "text/xml":
		return c.XML(http.StatusOK, cmt)
	default:
		return c.JSON(http.StatusOK, cmt)
	}
}

// addPost godoc
// @Summary Add post
// @Produce json
// @Produce xml
// @Param id formData int true "User Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Post
// @Router /posts/add [post]
func addPost(c echo.Context) error {
	fmt.Println("addPost")
	userid := c.FormValue("id")
	title := c.FormValue("title")
	body := c.FormValue("body")
	fmt.Println(userid + ":" + title + ":" + body)

	i, _ := strconv.Atoi(userid)
	post := Post{Userid: i, Title: title, Body: body}
	result := Db.Select("Userid", "Title", "Body").Create(&post)
	if result.Error != nil {
		fmt.Println("insert into posts error")
	}
	fmt.Println(result.RowsAffected)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusCreated, post)
	case "text/xml":
		return c.XML(http.StatusCreated, post)
	default:
		return c.JSON(http.StatusCreated, post)
	}
}

// addComment godoc
// @Summary Add comment for given post Id
// @Produce json
// @Produce xml
// @Param id formData int true "post Id"
// @Param name formData string true "Name"
// @Param email formData string true "E-mail"
// @Param body formData string true "Body"
// @Success 200 {object} Comment
// @Router /posts/{id}/comments/add [post]
func addComment(c echo.Context) error {
	fmt.Println("addComment")
	id, _ := strconv.Atoi(c.Param("id")) //post id
	name := c.FormValue("name")
	email := c.FormValue("email")
	body := c.FormValue("body")

	fmt.Println(name + ":" + email + ":" + body)

	comment := Comment{PostId: id, Name: name, Email: email, Body: body}
	result := Db.Select("PostId", "Name", "Email", "Body").Create(&comment)
	if result.Error != nil {
		fmt.Println("insert into comments error")
	}

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusCreated, comment)
	case "text/xml":
		return c.XML(http.StatusCreated, comment)
	default:
		return c.JSON(http.StatusCreated, comment)
	}
}

// editPost godoc
// @Summary Edit post
// @Produce json
// @Produce xml
// @Param userid formData int true "User Id"
// @Param id path int true "Post Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Post
// @Router /posts/{id}/edit [put]
func editPost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //post id
	userid, _ := strconv.Atoi(c.FormValue("userid"))
	title := c.FormValue("title")
	body := c.FormValue("body")

	fmt.Println("editPost:" + c.Param("id") + ":" + c.FormValue("userid") + ":" + title + ":" + body)

	p := Post{Userid: userid, ID: id, Title: title, Body: body}
	Db.Save(&p)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, p)
	case "text/xml":
		return c.XML(http.StatusOK, p)
	default:
		return c.JSON(http.StatusOK, p)
	}
}

// editComment godoc
// @Summary Edit comment for given Id
// @Produce json
// @Produce xml
// @Param post_id formData int true "post Id"
// @Param id path int true "comment Id"
// @Param name formData string true "Name"
// @Param email formData string true "E-mail"
// @Param body formData string true "Body"
// @Success 200 {object} Comment
// @Router /comments/{id}/edit [put]
func editComment(c echo.Context) error {
	postid, _ := strconv.Atoi(c.FormValue("post_id"))
	id, _ := strconv.Atoi(c.Param("id")) //comment id
	name := c.FormValue("name")
	email := c.FormValue("email")
	body := c.FormValue("body")

	fmt.Println("editComment:" + strconv.Itoa(postid) + ":" + strconv.Itoa(id) + ":" + name + ":" + email + ":" + body)

	comment := Comment{PostId: postid, ID: id, Name: name, Email: email, Body: body}
	Db.Save(&comment)

	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, comment)
	case "text/xml":
		return c.XML(http.StatusOK, comment)
	default:
		return c.JSON(http.StatusOK, comment)
	}
}

// deletePost godoc
// @Summary Delete post based on post ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /posts/{id}/delete [delete]
func deletePost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //post id
	var p Post
	p.ID = id
	Db.Debug().Delete(p)
	return c.NoContent(http.StatusNoContent)
}

// deleteComment godoc
// @Summary Delete comments based on comment ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Comment Id"
// @Success 200 {object} Comment
// @Router /comments/{id}/delete [delete]
func deleteComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //comment id

	var cmt Comment
	cmt.ID = id
	fmt.Println("deleteComment ", strconv.Itoa(id))
	Db.Debug().Delete(c)
	return c.NoContent(http.StatusNoContent)
}
