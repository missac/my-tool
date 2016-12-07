package handler

import (
	"dbmgm"
	"github.com/gin-gonic/gin"
	"io"
	"muslog"
	"os"
	"os/exec"
)

var th *TaskHandler

type TaskHandler struct {
}

func NewTaskHandler() *TaskHandler {
	once.Do(func() {
		th = &TaskHandler{}
	})
	return th
}

func (th *TaskHandler) Status(ctx *gin.Context) {
	appId := ctx.Query("appid")
	url := "curl 127.0.0.1:8080/api/v1/applications/" + appId
	muslog.Info("get status for task" + appId)
	cmd, err := exec.Command("/bin/sh", "-c", url).Output()
	if err != nil {
		muslog.Error(err)
		ctx.JSON(400, gin.H{
			"status": "faild",
			"res":    err.Error(),
		})
		return
	}
	//ctx.JSON(200, gin.H{
	//	"status": "success",
	//	"res":    string(cmd),
	//})
	ctx.String(200, string(cmd))
}

func (th *TaskHandler) List(ctx *gin.Context) {
	list, err := dbmgm.ListStreamingTask()
	if err != nil {
		ctx.JSON(400, gin.H{
			"list": err.Error(),
		})
		return
	}
	ctx.String(200, list)
}
func (th *TaskHandler) Delete(ctx *gin.Context) {
	taskId := ctx.Param("taskid")
	ctx.JSON(200, gin.H{
		"status": taskId,
	})
}

func (th *TaskHandler) Stop(ctx *gin.Context) {

}

func (th *TaskHandler) Start(ctx *gin.Context) {

}

func (th *TaskHandler) Update(ctx *gin.Context) {
	taskId := ctx.Param("taskid")
	appName := ctx.PostForm("appnama")
	taskDes := ctx.PostForm("taskdes")
	schStrategy := ctx.PostForm("schstrategy")
	execSeq := ctx.PostForm("execseq")
	codePath := ctx.PostForm("codepath")

	content, _, err := ctx.Request.FormFile("upload")
	out, err := os.Create(codePath)
	defer out.Close()
	if err != nil {
		muslog.Error(err)
		ctx.JSON(400, gin.H{
			"status": "failed",
			"msg":    err.Error(),
		})
		return
	}
	_, err = io.Copy(out, content)
	if err != nil {
		muslog.Error(err)
		ctx.JSON(400, gin.H{
			"status": 400,
			"msg":    err.Error(),
		})
		return
	}
	err = dbmgm.UpdateStreamingTask(taskId, appName, taskDes, schStrategy, execSeq, "appid")
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"msg":    err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
	})
}
