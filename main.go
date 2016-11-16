package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Joiner = "\\"
)

func init() {
	if os.IsPathSeparator('\\') {
		Joiner = "\\"
	} else {
		Joiner = "/"
	}
}

func main() {
	RunServer()
}

func RunServer() {
	r := gin.Default()

	r.LoadHTMLGlob("res/*.html")
	r.Static("/static", "res/static")
	r.Static("/archive", "archive")

	r.GET("/list", PostList)
	r.GET("/", EditorGet)
	r.POST("/new", EditorPost)
	r.PUT("/update", EditorPut)

	r.Run(":1024")
}

func PostList(c *gin.Context) {
	list := getFileList("./archive")
	fmt.Println(list)
	c.HTML(http.StatusOK, "list.html", gin.H{"list": list})
}

func EditorGet(c *gin.Context) {
	now := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(now, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	c.HTML(http.StatusOK, "editor.html", gin.H{"token": token})

}

func EditorPost(c *gin.Context) {

	path, err := createArchiveDir()
	if err != nil {
		fmt.Println("目录创建失败：", err)
		return
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)
	mdName := now + ".md"
	htmlName := now + ".html"

	NewFile(c.Request.FormValue("md"), path+Joiner+mdName)
	fmt.Println(NewPost(c, path+Joiner+htmlName))

	c.Redirect(http.StatusMovedPermanently, "/list")
}

func EditorPut(c *gin.Context) {
	fmt.Println(c.Request.FormValue("html"))
	fmt.Println(c.Request.FormValue("md"))
}

func createArchiveDir() (path string, err error) {
	now := time.Now().Format("2006-01-02")
	dirs := strings.Split(now, "-")
	path, _ = os.Getwd() //当前目录
	path = path + Joiner + "archive"

	for _, v := range dirs {

		path = path + Joiner + v

	}
	err = os.MkdirAll(path, os.ModePerm)
	return
}

func NewFile(str string, fPath string) (err error) {
	f, err := os.Open(fPath)

	if err != nil && os.IsNotExist(err) {
		f, err = os.Create(fPath)

		if err != nil && os.IsNotExist(err) {
			fmt.Println("文件创建失败：", err)
		}

		f.WriteString(str)

		fmt.Println("文件创建成功")
	} else {
		fmt.Println("文件已存在", fPath)
	}

	f.Close()

	return
}

type Post struct {
	Title string
	Time  string
	Md    string //md file path
	Html  string //html file path
}

func NewPost(c *gin.Context, fPath string) (err error) {
	NewFile("", fPath)

	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}
	defer f.Close()

	t, err := template.ParseFiles("tpl/detail.html")
	m := make(map[string]interface{})
	m["title"] = c.Request.FormValue("title")
	m["content"] = template.HTML(c.Request.FormValue("html"))
	err = t.Execute(f, m)

	//TODO 读写json文件
	var posts []Post
	data, _ := os.OpenFile("data/data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	b, _ := ioutil.ReadAll(data)
	fmt.Println(string(b))
	fmt.Println(json.Unmarshal(b, &posts))
	data.Close()

	return
}

func Build() {

	list := getFileList("./archive")

	os.MkdirAll("./data", os.ModePerm)

	fPath := "./data/data.json"

	f, err := os.Open(fPath)

	if err != nil && os.IsNotExist(err) {
		f, err = os.Create(fPath)

		if err != nil && os.IsNotExist(err) {
			fmt.Println("Data文件创建失败：", err)
		}

		fmt.Println("文件创建成功")

	}

	b, _ := json.Marshal(list)

	_, err = f.WriteString(string(b))

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

}

func getFileList(path string) (list []string) {

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		path = strings.Replace(path, "\\", "/", -1)

		if strings.EqualFold(filepath.Ext(path), ".html") {
			list = append(list, path)
		}
		return nil

	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
