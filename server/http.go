package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MenciusCheng/auto-seat/templates"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	t, err := LoadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/", home)
	r.GET("/ping", ping)
	r.Any("/delay", delay)

	r.POST("/upload", upload)

	return r
}

func LoadTemplate() (*template.Template, error) {
	t := template.New("")
	var err error

	t, err = t.New("index.tmpl").Parse(templates.IndexTmpl)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "自动座位",
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func delay(c *gin.Context) {
	ms := c.Query("ms")
	if len(ms) > 0 {
		fmt.Println(ms, "ms")
		num, err := strconv.Atoi(ms)
		if err == nil && num > 0 {
			time.Sleep(time.Duration(int64(num)) * time.Millisecond)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func upload(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file not found",
		})
		return
	}

	log.Println(file.Filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer src.Close()

	bytes, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	//fmt.Println("bytes", bytes)

	name := c.PostForm("name")
	fmt.Println("name", name)

	//c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

	// 导出文件
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	//c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(bytes)))
	c.Writer.Write(bytes)
}
