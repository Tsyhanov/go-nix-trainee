package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "test-http/docs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const secret = "secret"

type jwtCustomClaims struct {
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

var googleOauthConfig *oauth2.Config
var oauthStateString = "pseudo-random"

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

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
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

	//e.POST("/", login)
	//e.POST("/", handleGoogleLogin)
	e.GET("/", handleMain)
	e.GET("/login", handleGoogleLogin)
	e.GET("/callback", handleGoogleCallback)
	//e.GET("/swagger/*", echoSwagger.WrapHandler)

	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:      &jwtCustomClaims{},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:Authorization",
	}
	r.Use(middleware.JWTWithConfig(config))

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	r.GET("/posts", getPosts)
	r.GET("/posts/:id", getPostById)
	r.GET("/comments", getComments)
	r.GET("/comments/:id", getCommentById)
	r.GET("/posts/:id/comments", getCommentsByPostId)
	r.POST("/posts/add", addPost)
	r.POST("/posts/:id/comments/add", addComment)
	r.PUT("/posts/:id/edit", editPost)
	r.PUT("/comments/:id/edit", editComment)
	r.DELETE("/posts/:id/delete", deletePost)
	r.DELETE("comments/:id/delete", deleteComment)

	e.Logger.Fatal(e.Start(":8080"))
}

//jwt login
func login(c echo.Context) error {

	fmt.Println("login handle")

	username := c.FormValue("username")
	password := c.FormValue("password")
	username = "alex"
	password = "password"

	//check it from database
	if username != "alex" || password != "password" {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		Name: "Alex Tsyhanov",
		UUID: "999",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "Authorization",
		Value:    t,
		Path:     "/restricted",
		HttpOnly: true,
	})
	if err != nil {
		return err
	}

	return nil
	//	return c.JSON(http.StatusOK, map[string]string{
	//		"token": t,
	//	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

/*
func restricted(c echo.Context) error {
	fmt.Println("call restrictted")
	//	user := c.Get("user").(*jwt.Token)
	//	claims := user.Claims.(*jwtCustomClaims)
	//	name := claims.Name
	//	return c.String(http.StatusOK, "Welcome "+name+"!")
	return c.String(http.StatusOK, "Welcome !")
}
*/
func handleMain(c echo.Context) error {
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`
	return c.HTML(http.StatusOK, htmlIndex)
}

func handleGoogleLogin(c echo.Context) error {
	fmt.Println("handleGoogleLogin")
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	fmt.Println(url)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c echo.Context) error {
	fmt.Println("handleGoogleCallback")
	_, err := getUserInfo(c.FormValue("state"), c.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("handleGoogleCallback getUserInfo error")
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	//	return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")

	login(c)

	return c.Redirect(http.StatusTemporaryRedirect, "/restricted/swagger/index.html")

	/*
		if accessible(c) == nil {
			fmt.Println("handleGoogleCallback accessible!")
			//		return c.Redirect(http.StatusTemporaryRedirect, "/restricted/swagger/index.html")
			return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")

		} else {
			fmt.Println("handleGoogleCallback accessible error")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
	*/
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
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
