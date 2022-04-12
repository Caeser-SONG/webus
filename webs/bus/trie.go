package bus

import (
	"strings"
	"fmt"
)

// HTTP请求的路径恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。
// 我们通过树结构查询，如果中间某一层的节点都不满足条件，那么就说明没有匹配到的路由，查询结束。

// 路由节点
type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的部分
	children []*node // 子节点
	isWild   bool    // 是否精准
}

// 第一个匹配成功的节点 用于插入url
func (n *node) matchChild(part string) *node {
	// 遍历子节点
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点 ，用于查找
func (n *node) matchChlidren(part string) []*node {
	// var nodes []*node
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 对于路由来说，最重要的当然是注册与匹配了。
// 开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler

// trie的插入
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果高度是
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// trie 树的查找
func (n *node) search(parts []string, height int) *node{
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		if n.pattern == ""{
			return nil
		}
		return n
	}
	part := parts[height]
	children:= n.matchChlidren(part)
	for _,child := range children{
		result := child.search(parts, height + 1)
		if result != nil {
			 return result
		}
	}
	return nil

}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}
