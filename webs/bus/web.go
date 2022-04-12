package bus

import (
	"log"
	"net/http"
	"fmt"
	_ "log"
	"strings"
)

// bus执行的方法类型
// 处理逻辑的方法
type HandlerFunc func(c *Context)



type RouterGroup struct {
	perfix      string        // 前缀
	middlewares []HandlerFunc // 中间件执行的方法们
	parent      *RouterGroup  //父亲分组
	//  在初始化的时候就是传入的本身的bus
	bus         *Bus          // all groups share a Engine instance  也就是说 调用的都是同一个bus多对象的方法

}

// server
type Bus struct {
	router       *router
	*RouterGroup // 继承
	Groups       []*RouterGroup
}

// 创建一个busserver
func NewBus() *Bus {
	// 先初始化一个bus对象
	bus := &Bus{router: newRouter()}
	// 把包含bus对象的rg赋值给bus自己的   存的是不在group里的那些方法
	bus.RouterGroup = &RouterGroup{bus: bus}
	// 再把自己的rg变量付给自己的groups里面
	bus.Groups = []*RouterGroup{bus.RouterGroup}

	return bus
}

// 创建新的group
func (group *RouterGroup) Group(perfix string) *RouterGroup {
	bus := group.bus
	newGroup := &RouterGroup{
		perfix: group.perfix + perfix,
		parent: group,
		bus:    bus,
	}
	bus.Groups = append(bus.Groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.perfix + comp
	log.Printf("Rout %4s - %s", method, pattern)
	group.bus.router.addRoute(method, pattern, handler)
}

// 每个请求都先落到这个方法，接口里的方法
func (b *Bus) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("请求到了")
	// 且前用var声明后可以直接用
	// 存放这个请求要用的中间件list
	var middlewares []HandlerFunc

	for _,group := range b.Groups{
		// 简单通过 URL 的前缀来判断用哪些中间件
		if strings.HasPrefix(req.URL.Path,group.perfix){
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, req)
	// 中间件list付给本次的ctx
	c.handlers = middlewares
	b.router.handle(c)
}

// 绑定 方法和path
// func (b *Bus) AddRouter(method string, pattern string, handler HandlerFunc) {
// 	b.router.addRoute(method, pattern, handler)
// }

// 让请求和group绑定
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (b *Bus) Run(addr string) (err error) {
	return http.ListenAndServe(addr, b)
}

// 可以在调用的时候传入多个参数
func (group *RouterGroup) UseMiddleware(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares, middlewares...)
}
