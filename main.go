package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

//hugo  http://gohugo.org/
func main() {
	//router()

//	f, _ := os.Open("1111.txt")

//	buf := bufio.NewReader(f)

//	code, _ := buf.ReadString('\t')

//	code = strings.TrimSpace(code)
//	fmt.Println(code)

//	strconv.ParseInt("01001000", 2, 16)

	ReadLine("1111.txt" func(s string){
		fmt.Println(s)
	})

}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func router() {

	if len(os.Args) == 1 {
		help()
	} else if strings.EqualFold(os.Args[1], "new") {
		New()
	} else if strings.EqualFold(os.Args[1], "build") {
		Build()
	}
}

func help() {
	fmt.Println(`
 archive command [arguments]

 new          :创建目录和index.md
 new fileName :创建目录和指定名字文件
 build        :更新文章列表
	`)
}

func createDir() (path string, err error) {

	now := time.Now().Format("2006-01-02")
	dirs := strings.Split(now, "-")
	path, _ = os.Getwd()
	fmt.Println("当前目录：", path)
	path = path + Joiner + "archive"

	for _, v := range dirs {

		path = path + Joiner + v

	}

	err = os.MkdirAll(path, os.ModePerm)

	return
}

func New() {
	path, err := createDir()
	if err != nil {
		fmt.Println("目录创建失败：", err)
		return
	}
	fmt.Println("目录创建成功")

	fName := "index.md"

	fPath := path + Joiner + fName

	f, err := os.Open(fPath)

	if err != nil && os.IsNotExist(err) {
		f, err := os.Create(fPath)

		if err != nil && os.IsNotExist(err) {
			fmt.Println("文件创建失败：", err)
		}

		f.WriteString("#Title")
		fmt.Println("文件创建成功")

		defer f.Close()
	} else {
		fmt.Println("文件已存在", fPath)
	}

	defer f.Close()

}

func Edit() {

}

func Build() {
	getFileList("./archive")
}

//37439a79596609d0e0ab20ac37671fe7
//bbbcea93dbda07a315c3e99c74af7f67
//a87ff679a2f3e71d9181a67b7542122c
//bdbf46a337ac08e6b4677c2826519542
//adb79a4937b5a13851528d26525a4113
func getFileList(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil

	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
