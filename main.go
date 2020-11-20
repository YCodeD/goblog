package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是goblog!!</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈建议，请联系" + "<a href=\"mailto:xxx@xxx.com\">xxx@xxx.com</a>")
	} else {
		fmt.Fprint(w, "<h1>请求页面未找到：</h1>" + "<p>如有疑惑，请联系我们。</p>")
	}
	
	// fmt.Fprint(w, "请求路径为：" + r.URL.Path)
	// fmt.Println(r.Header.Get("User-Agent")) // 获取客户端信息
	// w.WriteHeader(http.StatusInternalServerError) // 返回500状态码
}

func main()  {
	http.HandleFunc("/", handlerFunc)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println(err)
	}

}