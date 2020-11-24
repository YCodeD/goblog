package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/route"
)

var router *mux.Router
var db *sql.DB

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

// Article 对应一条文章数据, 用以存储从数据库里读出来的文章数据
type Article struct {
	Title, Body string
	ID          int64
}

// Delete 方法用以从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("delete from articles where id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// 删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// 编辑文章
func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3 - 40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// 更新文章
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	_, err := getArticleByID(id)
	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")

		}
	} else {
		// 4. 未出现错误

		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			// 4.2 表单验证通过，更新数据
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			// 更新成功，跳转到文章详情页
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没做任何更改！")
			}
		} else {

			// 4.3 表单验证不通过，显示理由
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			tmpl.Execute(w, data)
		}

	}
}

// 删除文章
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")

		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
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

func main() {

	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	// router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	// router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取URL示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL:", homeURL)
	// articlesURL, _ := router.Get("articles.show").URL("id", "23")
	// fmt.Println("articlesURL:", articlesURL)

	http.ListenAndServe(":3001", removeTrailingSlash(router))

}
