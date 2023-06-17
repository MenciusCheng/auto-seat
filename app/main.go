package main

import (
	"github.com/MenciusCheng/auto-seat/server"
)

func main() {
	r := server.InitRouter()
	err := r.Run(":8080") // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(err)
	}
}
