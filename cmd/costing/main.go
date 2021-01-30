package main

import (
	"fmt"

	"github.com/KeisukeIwabuchi/Costing/internal/apps/totalcosting"
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

	var box totalcosting.Box
	box.Run()
}
