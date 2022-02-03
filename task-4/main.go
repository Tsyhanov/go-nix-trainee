package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "./docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

var Db *gorm.DB
var Err error

type Posts []struct {
	Userid int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comments []struct {
	PostId int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
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

func isRegExpMatched(pattern string, source string) bool {

	matched, _ := regexp.MatchString(pattern, source)

	if matched {
		return true
	} else {
		return false
	}
}

// getPosts godoc
// @Summary Retrieves posts
// @Produce json
// @Produce xml
// @Success 200 {object} Post
// @Router /posts [get]
func getPosts() (p []Post) {
	fmt.Println("getPosts")

	result := Db.Find(&p)

	if result.Error != nil {
		fmt.Println("select from posts error")
	}
	fmt.Println(result.RowsAffected)
	return p
}

// getPostById godoc
// @Summary Retrieves posts based on given ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /posts/{id} [get]
func getPostById(id int) (p Post) {
	fmt.Println("getPostById")

	result := Db.First(&p, id)

	if result.Error != nil {
		fmt.Println("select from posts error")
	}
	fmt.Println(result.RowsAffected)
	return p
}

// getComments godoc
// @Summary Retrieves comments
// @Produce json
// @Produce xml
// @Success 200 {object} Comment
// @Router /comments [get]
func getComments() (c []Comment) {
	fmt.Println("getComments")

	result := Db.Find(&c)

	if result.Error != nil {
		fmt.Println("select from comments error")
	}
	fmt.Println(result.RowsAffected)
	return c
}

// getCommentsByPostId godoc
// @Summary Retrieves comments based on post ID
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Comment
// @Router /comments/{id} [get]
func getCommentsByPostId(post_id int) (c []Comment) {
	fmt.Println("getComments by Post id")

	result := Db.Where("post_id = ?", post_id).Find(&c)

	if result.Error != nil {
		fmt.Println("select from comments error")
	}
	fmt.Println(result.RowsAffected)
	return c
}

// addPost godoc
// @Summary Add post
// @Produce json
// @Produce xml
// @Param id formData int true "User Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Comment
// @Router /posts/add [post]
func addPost(userid string, title string, body string) {
	fmt.Println("addPost")
	fmt.Println(userid + ":" + title + ":" + body)

	i, _ := strconv.Atoi(userid)
	post := Post{Userid: i, Title: title, Body: body}
	result := Db.Select("Userid", "Title", "Body").Create(&post)
	if result.Error != nil {
		fmt.Println("insert into posts error")
	}
	fmt.Println(result.RowsAffected)
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
// @Router /comments/add [post]
func addComment(postid string, name string, email string, body string) {
	fmt.Println("addComment")
	fmt.Println(postid + ":" + name + ":" + email + ":" + body)

	i, _ := strconv.Atoi(postid)
	comment := Comment{PostId: i, Name: name, Email: email, Body: body}
	result := Db.Select("PostId", "Name", "Email", "Body").Create(&comment)
	if result.Error != nil {
		fmt.Println("insert into comments error")
	}
	fmt.Println(result.RowsAffected)
}

// editPost godoc
// @Summary Edit post
// @Produce json
// @Produce xml
// @Param userid formData int true "User Id"
// @Param id path int true "Post Id"
// @Param title formData string true "Post Title"
// @Param body formData string true "Post Body"
// @Success 200 {object} Comment
// @Router /posts/{id}/edit [put]
func editPost(userid string, id int, title string, body string) {
	fmt.Println("editPost:" + userid + ":" + strconv.Itoa(id) + ":" + title + ":" + body)

	i, _ := strconv.Atoi(userid)
	p := Post{Userid: i, ID: id, Title: title, Body: body}
	Db.Save(&p)
}

// editComment godoc
// @Summary Edit comment for given Id
// @Produce json
// @Produce xml
// @Param pid path int true "post Id"
// @Param post_id formData int true "post Id"
// @Param id path int true "comment Id"
// @Param name formData string true "Name"
// @Param email formData string true "E-mail"
// @Param body formData string true "Body"
// @Success 200 {object} Comment
// @Router /posts/{pid}/comment/{id}/edit [put]
func editComment(postid string, id int, name string, email string, body string) {
	fmt.Println("editComment:" + postid + ":" + strconv.Itoa(id) + ":" + name + ":" + email + ":" + body)

	i, _ := strconv.Atoi(postid)
	c := Comment{PostId: i, ID: id, Name: name, Email: email, Body: body}
	Db.Save(&c)
}

// deletePost godoc
// @Summary Delete post based on post ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Post Id"
// @Success 200 {object} Post
// @Router /posts/{id}/delete [delete]
func deletePost(post_id int) {
	var p Post
	p.ID = post_id
	fmt.Println("deletePost ", strconv.Itoa(post_id))
	Db.Debug().Delete(p)
}

// deleteComment godoc
// @Summary Delete comments based on comment ID
// @Accept json
// @Produce json
// @Produce xml
// @Param id path int true "Comment Id"
// @Success 200 {object} Post
// @Router /posts/{post_id}/comment/{id}/delete [delete]
func deleteComment(comment_id int) {
	var c Comment
	c.ID = comment_id
	fmt.Println("deleteComment ", strconv.Itoa(comment_id))
	Db.Debug().Delete(c)
}

//mux
func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	case "PUT":
		handleModify(w, r)
	case "DELETE":
		handleDelete(w, r)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleGet")
	var d interface{}
	ep := r.URL.Path
	fmt.Println(ep)
	switch ep {
	case "/swagger/index.html":
		fmt.Println("swagger")
		http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	case "/posts":
		d = getPosts()
	case "/comments":
		d = getComments()
	default:
		switch {
		case isRegExpMatched(`/posts/\d*/comments`, ep):
			ep = path.Dir(ep)
			i, _ := strconv.Atoi(path.Base(ep))
			d = getCommentsByPostId(i)
		case isRegExpMatched(`/posts/\d*`, ep):
			i, _ := strconv.Atoi(path.Base(r.URL.Path))
			d = getPostById(i)
		default:
			fmt.Println("regexp default")
		}
	}

	if r.Header.Get("Content-Type") == "application/xml" {
		x, err := xml.MarshalIndent(d, "", "")
		if err != nil {
			fmt.Println("xml error", err)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(x)
	} else {
		b, err := json.Marshal(d)
		if err != nil {
			fmt.Println("json error", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlePost")
	ep := r.URL.Path
	r.ParseMultipartForm(0)
	fmt.Println(ep)
	switch {
	case isRegExpMatched(`/posts/add`, ep):
		addPost(r.FormValue("id"), r.FormValue("title"), r.FormValue("body"))
	case isRegExpMatched(`/comments/add`, ep):
		addComment(r.FormValue("id"), r.FormValue("name"), r.FormValue("email"), r.FormValue("body"))
	default:
		fmt.Println("regexp default")
	}
}

func handleModify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleModify")
	ep := r.URL.Path
	fmt.Println(ep)
	switch {
	case isRegExpMatched(`/posts/\d*/comment/\d*/edit`, ep):
		ep = path.Dir(ep)
		i, _ := strconv.Atoi(path.Base(ep))
		editComment(r.FormValue("post_id"), i, r.FormValue("name"), r.FormValue("email"), r.FormValue("body"))
	case isRegExpMatched(`/posts/\d*/edit`, ep):
		ep = path.Dir(ep)
		i, _ := strconv.Atoi(path.Base(ep))
		editPost(r.FormValue("userid"), i, r.FormValue("title"), r.FormValue("body"))
	default:
		fmt.Println("regexp default")
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleDelete")
	ep := r.URL.Path
	fmt.Println(ep)
	switch {
	case isRegExpMatched(`/posts/\d*/comment/\d*/delete`, ep):
		ep = path.Dir(ep)
		i, _ := strconv.Atoi(path.Base(ep))
		deleteComment(i) //delete comment
	case isRegExpMatched(`/posts/\d*/delete`, ep):
		ep = path.Dir(ep)
		i, _ := strconv.Atoi(path.Base(ep))
		deletePost(i)
	default:
		fmt.Println("regexp default")
	}
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

// @host 127.0.0.1:8080
func main() {
	fmt.Println("Start")
	//connect to mysql. We can do it once as gorm support pooling
	dsn := "root:test@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	Db, Err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if Err != nil {
		fmt.Println("error connection")
	}
	fmt.Println("connection is OK")
	Db.AutoMigrate(&Post{}, &Comment{})

	//start web server
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", handleRequest)

	server.ListenAndServe()

	fmt.Println("Done")
}
