package main

import (
	"fmt"
	"net/http"
	b "webs/bus"
)

func main() {
	r := b.NewBus()
	// 全局中间件
	r.UseMiddleware(b.Logger())
	r.POST("/login", func(c *b.Context) {
		c.JSON(http.StatusOK, b.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password")})

	})
	r.GET("/", func(c *b.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *b.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *b.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	// r.GET("/assets/*filepath", func(c *b.Context) {
	// 	c.JSON(http.StatusOK, b.H{"filepath": c.Param("filepath")}
	// })

	v1 := r.Group("/v1")
	v1.UseMiddleware(func (c *b.Context) {
		fmt.Println("zzzzzzzzzzzzzz")
	})
	v1.GET("/book", func(c *b.Context) {
		c.String(http.StatusOK, "book")
	})
	r.GET("/panic", func(c *b.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[1])
	})
	r.Run(":8000")
}
