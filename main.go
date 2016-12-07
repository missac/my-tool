package main

import (
	"dbmgm"
	"github.com/gin-gonic/gin"
	"muslog"
	"router"
)

func main() {
	muslog.InitLog(muslog.DEBUG, "", muslog.LevelTrace)
	muslog.Info("init log with DEBUG and trace level")
	dbmgm.InitDB("./mustang.db")
	dbmgm.ListStreamingTask()
	r := gin.Default()
	r.LoadHTMLFiles("test.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "test.html", nil)
	})
	router.NewRouter().Register(r)

	muslog.Info("server start at port 9998")
	r.Run(":9998")
}
