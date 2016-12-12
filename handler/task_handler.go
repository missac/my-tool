package handler

import (
	"bufio"
	"core"
	"dbmgm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"muslog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var taskOnce sync.Once
var taskHander *TaskHandler

type TaskHandler struct {
	wsTaskMap map[string]io.ReadCloser
	test      string
}

func NewTaskHandler() *TaskHandler {
	taskOnce.Do(func() {
		wtm := make(map[string]io.ReadCloser, 0)
		taskHander = &TaskHandler{wtm, "hahah"}
	})
	return taskHander
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
	//taskId := ctx.Param("taskid")

}

func (th *TaskHandler) Start(ctx *gin.Context) {
	taskId := ctx.Param("taskid")
	muslog.Info("start task: " + taskId)

	//execSeq, err := dbmgm.GetExecSeq(taskId)
	var execSeq = "aa"
	//	if err != nil {
	//		ctx.JSON(400, gin.H{
	//			"res": err.Error(),
	//		})
	//		return
	//	}

	if _, ok := th.wsTaskMap[taskId]; ok {
		msg := fmt.Sprintf("job %s already exist", taskId)
		muslog.Warning(msg)
		ctx.JSON(400, gin.H{
			"res": msg,
		})
		return
	} else {
		mgm := core.NewTaskMgm()
		ri, pid, err := mgm.StartTask(execSeq)
		if err != nil {
			ctx.JSON(400, gin.H{
				"res": err.Error(),
			})
			return
		}
		dbmgm.UpdatePid(taskId, strconv.Itoa(pid))
		th.wsTaskMap[taskId] = ri

		ctx.JSON(200, gin.H{
			"res": "success",
		})
	}
}

func (th *TaskHandler) WSocket(ctx *gin.Context) {
	taskId := ctx.Param("taskid")
	th.wsHandler(ctx.Writer, ctx.Request, taskId)
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

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (th *TaskHandler) wsHandler(w http.ResponseWriter, r *http.Request, taskId string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	if _, ok := th.wsTaskMap[taskId]; ok {
		ri, _ := th.wsTaskMap[taskId]
		sr := bufio.NewScanner(ri)
		for sr.Scan() {
			println(string(sr.Bytes()))
			err = conn.WriteMessage(websocket.TextMessage, sr.Bytes())
			if err != nil {
				println("client close websocket")
				conn.Close()
				break
			}
		}
		muslog.Info("no thing to send for websocket: " + taskId)
		muslog.Info("remove from taskmap:" + taskId)
		delete(th.wsTaskMap, taskId)
		conn.WriteMessage(websocket.TextMessage, []byte("nothing to send, close websocket"))
	} else {
		var msg = "can`t find job: " + taskId
		muslog.Info(msg)
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
		muslog.Info("websocket closed")
		conn.Close()
		return
	}
}
