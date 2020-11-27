package main

import (
	"fmt"
	"net/http"

	"goblog/app/http/middlewares"
	"goblog/app/models/article"
	"goblog/bootstrap"
	"goblog/config"
	"goblog/pkg/model"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {

	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	var article article.Article
	model.DB.First(&article)
	fmt.Println(article)

	http.ListenAndServe(":3001", middlewares.RemoveTrailingSlash(router))

}
