package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"task-7/database"
	//_ "task-7/database"
	_ "task-7/docs"

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
	database.Db, database.Err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if database.Err != nil {
		fmt.Println("error connection")
	}
	fmt.Println("connection is OK")
	database.Db.AutoMigrate(&database.Post{}, &database.Comment{})

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

	r.GET("/posts", database.GetPosts)
	r.GET("/posts/:id", database.GetPostById)
	r.GET("/comments", database.GetComments)
	r.GET("/comments/:id", database.GetCommentById)
	r.GET("/posts/:id/comments", database.GetCommentsByPostId)
	r.POST("/posts/add", database.AddPost)
	r.POST("/posts/:id/comments/add", database.AddComment)
	r.PUT("/posts/:id/edit", database.EditPost)
	r.PUT("/comments/:id/edit", database.EditComment)
	r.DELETE("/posts/:id/delete", database.DeletePost)
	r.DELETE("comments/:id/delete", database.DeleteComment)

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
