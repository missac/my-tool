package router

import (
	"github.com/gin-gonic/gin"
	"handler"
)

type StreamingRouter struct {
	name string
	desc string
}

func newStreamingRouter(name string, desc string) RouterComponent {
	return &StreamingRouter{
		name: name,
		desc: desc,
	}
}
func (s *StreamingRouter) registerRouter(r *gin.Engine) {
	group := r.Group("/" + s.name)
	group.GET("/test", handler.NewStreamingHandler().TestFunc)
	group.GET("/template", handler.NewStreamingHandler().GetTemplate)
	group.POST("/submit", handler.NewStreamingHandler().Submit)
}
