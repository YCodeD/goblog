package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/requests"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"

	"gorm.io/gorm"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// fmt.Printf("%#v", article)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		// 4. 读取成功,显示文章

		view.Render(w, view.D{
			"Article": article,
		}, "articles.show")

		// // 4.0 设置模板相对路径
		// viewDir := "resources/views"

		// // 4.1 所有布局模板文件 Slice
		// files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		// logger.LogError(err)

		// // 4.2 在 Slice 里新增我们的目标文件
		// newFiles := append(files, viewDir+"/articles/show.gohtml")

		// // 4.3 解析模板文件
		// tmpl, err := template.New("show.gohtml").
		// 	Funcs(template.FuncMap{
		// 		"RouteName2URL": route.Name2URL,
		// 	}).ParseFiles(newFiles...)
		// logger.LogError(err)

		// // 4.4 渲染模板，将所有文章的数据传输进去
		// tmpl.ExecuteTemplate(w, "app", article)
	}
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 执行查询语句，返回一个结果集
	articles, err := article.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// ---  2. 加载模板 ---

		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index")

		// // 2.0 设置模板相对路径
		// viewDir := "resources/views"

		// // 2.1 所有布局模板文件 Slice
		// files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		// logger.LogError(err)

		// // 2.2 在 Slice 里新增我们的目标文件
		// newFiles := append(files, viewDir+"/articles/index.gohtml")

		// // 2.3 解析模板文件
		// tmpl, err := template.ParseFiles(newFiles...)
		// logger.LogError(err)

		// // 2.4 渲染模板，将所有文章的数据传输进去
		// tmpl.ExecuteTemplate(w, "app", articles)
	}

}

// ArticlesFormData 创建博文表单数据
// type ArticlesFormData struct {
// 	Title, Body string
// 	Article     article.Article
// 	Errors      map[string]string
// }

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	view.Render(w, view.D{}, "articles.create", "articles._form_field")

	// storeURL := route.Name2URL("articles.store")
	// data := ArticlesFormData{
	// 	Title:  "",
	// 	Body:   "",
	// 	URL:    storeURL,
	// 	Errors: nil,
	// }

	// tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	// tmpl.Execute(w, data)
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body:  r.PostFormValue("body"),
	}

	errors := requests.ValidateArticleForm(_article)

	// 检查是否有错误
	if len(errors) == 0 {

		_article.Create()

		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
            http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "文章创建失败，请联系管理员")
		}
	} else {

		view.Render(w, view.D{
			"Article": _article,
			"Errors": errors,
		}, "articles.create", "articles._form_field")

		// 对 errors 的传参到html中进行渲染，使用到标准库html/template
		// storeURL := route.Name2URL("articles.store")

		// data := ArticlesFormData{
		// 	Title:  title,
		// 	Body:   body,
		// 	URL:    storeURL,
		// 	Errors: errors,
		// }

		// tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		// if err != nil {
		// 	panic(err)
		// }

		// tmpl.Execute(w, data)
	}
}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  view.D{},
		}, "articles.edit", "articles._form_field")
	}
}

// Update 更新文章
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)
	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		_article.Title = r.PostFormValue("title")
		_article.Body = r.PostFormValue("body")

		errors := requests.ValidateArticleForm(_article)

		if len(errors) == 0 {

			// 4.2 表单验证通过，更新数据
			RowsAffected, err := _article.Update()

			if err != nil {
				// 数据库错误
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			// 更新成功，跳转到文章详情页
			if RowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没做任何更改！")
			}
		} else {

			// 4.3 表单验证不通过，显示理由
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  errors,
			}, "articles.edit", "articles._form_field")

		}

	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouterVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		rowsAffected, err := _article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
