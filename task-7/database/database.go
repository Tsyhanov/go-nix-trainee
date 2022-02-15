package database

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	_ "task-7/docs"
)

var Db *gorm.DB
var Err error

type User struct {
	//gorm.Model
	Email string
	ID    int
}

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

// GetPosts godoc
// @Summary Retrieves posts
// @Produce json
// @Produce xml
// @Success 200 {object} Post
// @Router /restricted/posts [get]
func GetPosts(c echo.Context) error {
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

// GetComments godoc
// @Summary Retrieves comments
// @Produce json
// @Produce xml
// @Success 200 {object} Comment
// @Router /restricted/comments [get]
func GetComments(c echo.Context) error {
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

// GetPostById godoc
// @Summary Retrieves posts based on given ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /restricted/posts/{id} [get]
func GetPostById(c echo.Context) error {
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

// GetCommentById godoc
// @Summary Retrieves comment based on given ID
// @Produce json
// @Produce xml
// @Param id path int true "Comment Id"
// @Success 200 {object} Post
// @Router /restricted/comments/{id} [get]
func GetCommentById(c echo.Context) error {
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

// GetCommentsByPostId godoc
// @Summary Retrieves comments based on post ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Comment
// @Router /restricted/posts/{id}/comments [get]
func GetCommentsByPostId(c echo.Context) error {
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

// AddPost godoc
// @Summary Add post
// @Produce json
// @Produce xml
// @Param id formData int true "User Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Post
// @Router /restricted/posts/add [post]
func AddPost(c echo.Context) error {
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

// AddComment godoc
// @Summary Add comment for given post Id
// @Produce json
// @Produce xml
// @Param id formData int true "post Id"
// @Param name formData string true "Name"
// @Param email formData string true "E-mail"
// @Param body formData string true "Body"
// @Success 200 {object} Comment
// @Router /restricted/posts/{id}/comments/add [post]
func AddComment(c echo.Context) error {
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

// EditPost godoc
// @Summary Edit post
// @Produce json
// @Produce xml
// @Param userid formData int true "User Id"
// @Param id path int true "Post Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Post
// @Router /restricted/posts/{id}/edit [put]
func EditPost(c echo.Context) error {
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

// EditComment godoc
// @Summary Edit comment for given Id
// @Produce json
// @Produce xml
// @Param post_id formData int true "post Id"
// @Param id path int true "comment Id"
// @Param name formData string true "Name"
// @Param email formData string true "E-mail"
// @Param body formData string true "Body"
// @Success 200 {object} Comment
// @Router /restricted/comments/{id}/edit [put]
func EditComment(c echo.Context) error {
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

// DeletePost godoc
// @Summary Delete post based on post ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /restricted/posts/{id}/delete [delete]
func DeletePost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //post id
	var p Post
	p.ID = id
	Db.Debug().Delete(p)
	return c.NoContent(http.StatusNoContent)
}

// DeleteComment godoc
// @Summary Delete comments based on comment ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Comment Id"
// @Success 200 {object} Comment
// @Router /restricted/comments/{id}/delete [delete]
func DeleteComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //comment id

	var cmt Comment
	cmt.ID = id
	fmt.Println("deleteComment ", strconv.Itoa(id))
	Db.Debug().Delete(c)
	return c.NoContent(http.StatusNoContent)
}
