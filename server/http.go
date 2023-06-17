package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", home)
	r.GET("/ping", ping)
	r.POST("/upload", upload)

	return r
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
