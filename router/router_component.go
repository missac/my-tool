package router

import (
	"github.com/gin-gonic/gin"
)

type RouterComponent interface {
	registerRouter(r *gin.Engine)
}

type Router struct {
	children []RouterComponent
}

func NewRouter() *Router {
	r := &Router{}
	r.init()
	return r
}

func (rou *Router) init() {
	rou.add(newStreamingRouter("streaming", "streaming"))
	rou.add(newTaskRouter("task", "task"))
}

func (rou *Router) add(m RouterComponent) {
	rou.children = append(rou.children, m)
}
func (rou *Router) Register(r *gin.Engine) {
	for _, v := range rou.children {
		v.registerRouter(r)
	}
}
