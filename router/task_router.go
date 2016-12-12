package router

import (
	"github.com/gin-gonic/gin"
	"handler"
)

type TaskRouter struct {
	name string
	desc string
}

func newTaskRouter(name string, desc string) RouterComponent {
	return &TaskRouter{
		name: name,
		desc: desc,
	}
}
func (t *TaskRouter) registerRouter(r *gin.Engine) {
	group := r.Group("/" + t.name)
	group.GET("/status", handler.NewTaskHandler().Status)
	group.GET("/list", handler.NewTaskHandler().List)
	group.DELETE("/:taskid", handler.NewTaskHandler().Delete)
	group.POST("/update/:taskid", handler.NewTaskHandler().Update)
	group.POST("/start/:taskid", handler.NewTaskHandler().Start)
	group.GET("/ws/:taskid", handler.NewTaskHandler().WSocket)
	group.POST("/stop/:taskid", handler.NewTaskHandler().Stop)
}
