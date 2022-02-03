package main

import (
	"gorm.io/driver/mysql"

	//_ "gorm.io/driver/mysql"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
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

//insert single post structure into db and request comments for this post
func AddPostToDb(wgposts *sync.WaitGroup, userid int, id int, title string, body string) {
	defer wgposts.Done()

	post := Post{Userid: userid, ID: id, Title: title, Body: body}
	result := Db.Create(&post)
	if result.Error != nil {
		fmt.Println("insert into posts error")
	}
	fmt.Println(result.RowsAffected)

	//get comments for postID
	GetComments(strconv.Itoa(id))
}

//insert single comment struct into db
func AddCommentToDb(wgcomments *sync.WaitGroup, postid int, id int, name string, email string, body string) {
	defer wgcomments.Done()

	comment := Comment{PostId: postid, ID: id, Name: name, Email: email, Body: body}
	result := Db.Create(&comment)
	if result.Error != nil {
		fmt.Println("insert into comments error")
	}
	fmt.Println(result.RowsAffected)
	fmt.Println("AddCommentToDb for postid" + strconv.Itoa(postid) + ": " + strconv.Itoa(id))
}

//get comments for PostId and start routines to insert it into db
func GetComments(postid string) {
	req := "https://jsonplaceholder.typicode.com/comments?postId=" + postid

	resp, err := http.Get(req)
	if err != nil {
		fmt.Printf("Request Failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll Failed: %s", err)
	}

	comments := Comments{}
	err = json.Unmarshal(body, &comments)
	if err != nil {
		log.Printf("Comments unmarshaling failed: %s", err)
		return
	}

	//create subroutines to insert comments into db
	var wgcomments sync.WaitGroup
	for _, value := range comments {
		wgcomments.Add(1)
		go AddCommentToDb(&wgcomments, value.PostId, value.ID, value.Name, value.Email, value.Body)
	}
	wgcomments.Wait()
}

//Main
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

	//get posts for userId=7
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err != nil {
		log.Printf("Request Failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll Failed: %s", err)
	}
	posts := Posts{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Printf("Posts unmarshaling failed: %s", err)
		return
	}

	//create routines to insert posts into db
	var wgposts sync.WaitGroup
	for _, value := range posts {
		wgposts.Add(1)
		go AddPostToDb(&wgposts, value.Userid, value.ID, value.Title, value.Body)
	}

	wgposts.Wait()
	fmt.Println("Done")
}
