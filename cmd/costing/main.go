package main

import (
	"fmt"

	gin "github.com/gin-gonic/gin"
)

func main() {
	message := "原価計算しよう"

	fmt.Println(message)

	router := gin.Default()
	router.LoadHTMLGlob("web/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"data": message,
		})
	})

	router.Run()
}
