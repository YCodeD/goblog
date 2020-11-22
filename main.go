package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

/* v1
func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是goblog!!</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈建议，请联系" + "<a href=\"mailto:xxx@xxx.com\">xxx@xxx.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到：</h1>" + "<p>如有疑惑，请联系我们。</p>")
	}

	// fmt.Fprint(w, "请求路径为：" + r.URL.Path)
	// fmt.Println(r.Header.Get("User-Agent")) // 获取客户端信息
	// w.WriteHeader(http.StatusInternalServerError) // 返回500状态码
}
*/

/* v2
func defaultHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是goblog!!</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到：</h1>" +
		 "<p>如有疑惑，请联系我们。</p>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈建议，请联系" +
	 "<a href=\"mailto:xxx@xxx.com\">xxx@xxx.com</a>")

}
*/

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	
	fmt.Fprint(w, "<h1>Hello, 欢迎来到goblog!!</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request)  {
	
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈建议，请联系" +
	 "<a href=\"mailto:xxx@xxx.com\">xxx@xxx.com</a>")
	
}

func notFoundHandler(w http.ResponseWriter, r *http.Request)  {

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID：" + id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "访问文章列表")
}

// 添加新文章
func articlesStoreHandler(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		// 解析错误
		fmt.Fprint(w, "请提供正确的数据！")
		return
	}

	title := r.PostForm.Get("title")

	fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	fmt.Fprintf(w, "title 的值为：%v", title)


	// 可直接使用r.FormValue r.PostFormValue 获取请求内容, 无需使用r.ParseForm
	// fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
    // fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
    // fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
    // fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))


}

// 标头中间件
func forceHTMLMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		h.ServeHTTP(w, r)
	})
}

// 
func articlesCreateHandler(w http.ResponseWriter, r *http.Request)  {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<title>创建文章 -- 我的技术博客</title>
</head>
<body>
	<form action="%s?test=data" method="post">
		<p><input type="text" name="title"></p>
		<p><textarea name="body" cols="30" rows="10"></textarea></p>
		<p><button type="submit">提交</button></p>
	</form>
</body>
</html>
`

	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 除首页以外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		
		// 2. 将请求传递下去
		next.ServeHTTP(w, r)
	})
}


func main()  {

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取URL示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL("id", "23")
	// fmt.Println("articlesURL:", articlesURL)


	http.ListenAndServe(":3001", removeTrailingSlash(router))

}