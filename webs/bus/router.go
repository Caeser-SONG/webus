package bus

import (
	"fmt"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	handlers map[string]HandlerFunc // 存储每种请求方式的 HandlerFunc
	roots    map[string]*node       // 存储每种请求方式的Trie 树根节点
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

// 拆分路径  返回节点数组
func parsePattern(pattern string) []string {
	// 拆分路径
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range vs {
		if part != "" {
			parts = append(parts, part)
			// 这里为啥判断？？？
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	// 如果不存在则新建一个根节点
	// if _,ok := r.roots[method];!ok {
	// 	r.roots[method] = &node{}
	// }
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 起始位置是根节点
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)

	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				fmt.Println(params)
				fmt.Println(parts)
				fmt.Println(searchParts)
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 路由到方法后执行方法
func (r *router) handle(c *Context) {
	// 获取参数
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 把要执行的方法加到中间件方法后面
		c.handlers = append(c.handlers,r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func (c *Context){
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 方法在next里执行
	c.Next()
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}
