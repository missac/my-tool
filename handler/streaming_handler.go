package handler

import (
	"bufio"
	"core"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"mime/multipart"
	"muslog"
	"net/http"
	"os"
	"time"
)

var m *StreamingHandler

type StreamingHandler struct {
	wsTaskMap map[string]io.ReadCloser
}

func NewStreamingHandler() *StreamingHandler {
	once.Do(func() {
		wtm := make(map[string]io.ReadCloser, 0)
		m = &StreamingHandler{wtm}
	})
	return m
}

func (sh *StreamingHandler) TestFunc(ctx *gin.Context) {
	//	mgm := core.NewStreamingMgm()
	//	go mgm.StartJob()
	//	ctx.Status(http.StatusOK)
	muslog.Info("handler request for test")
	sh.wsHandler(ctx.Writer, ctx.Request, false, "aaaa")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (sh *StreamingHandler) wsHandler(w http.ResponseWriter, r *http.Request, stopOnClose bool, taskId string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	var ri io.ReadCloser
	var pid int
	if _, ok := sh.wsTaskMap[taskId]; ok {
		ri, _ = sh.wsTaskMap[taskId]
		println("job already exist")
	} else {
		mgm := core.NewStreamingMgm()
		ri, pid, _ = mgm.StartJob()
		sh.wsTaskMap[taskId] = ri
	}

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
	if stopOnClose {
		//kill the runing task
		println(pid)
		//clear wstaskmap
	}
}

func (sh *StreamingHandler) GetTemplate(ctx *gin.Context) {
	sourceType := ctx.Query("datasource")
	srcType := ctx.Query("codetype")
	muslog.Info(fmt.Sprintf("handle Request for get template with sourceType: %s and CodeType %s", sourceType, srcType))
	mgm := core.NewStreamingMgm()
	mgm.DataSource = sourceType
	mgm.CodeType = srcType
	template, err := mgm.GenTemplate()
	if err != nil {
		muslog.Error(err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, template)
}

func (sh *StreamingHandler) Submit(ctx *gin.Context) {
	mgm := core.NewStreamingMgm()

	//	mgm.AppName = ctx.DefaultPostForm("appname", "streaming")
	mgm.Master = ctx.DefaultPostForm("master", core.StandAlone)
	mgm.DeployMode = ctx.DefaultPostForm("deploymode", core.DeployClient)
	//mgm.TaskDes  = ctx.DefaultPostForm("taksdes", "streaming task")
	//mgm.TaskType = ctx.DefaultPostForm("takstype", "streaming")
	//	mgm.ClassName = ctx.PostForm("classname")
	//	mgm.DriverMem = ctx.DefaultPostForm("drivermem", "1G")
	//	mgm.ExecutorMem = ctx.DefaultPostForm("executormem", "1G")
	//	mgm.DriverCores = ctx.DefaultPostForm("dirvercores", 1)
	//	mgm.TotalExecutorCores = ctx.DefaultPostForm("totalexecutorcores", 0)
	//	mgm.ExecutorCores = ctx.DefaultPostForm("executorcores", 0)
	//	mgm.ExecutorNums = ctx.DefaultPostForm("executornums", 2)
	//
	//
	//	mgm.DataSource = ctx.PostForm("datasource")
	//	mgm.CodeType = ctx.PostForm("codetype")
	mgm.Parallelism = ctx.DefaultPostForm("parallelism", "0")
	mgm.StorgeMemoryFraction = ctx.DefaultPostForm("storgememoryfraction", "0.6")
	mgm.ShuffleMemoryFraction = ctx.DefaultPostForm("shuffermemoryfraction", "0.2")

	mgm.ScheduleMode = ctx.DefaultPostForm("schedulemode", core.ScheduleNow)
	//useWebsock := ctx.DefaultPostForm("", core.ScheduleNow)

	// upload pyfile or jars
	file, header, err := ctx.Request.FormFile("upload")
	print(file)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fileName := header.Filename
	fileType := ctx.DefaultPostForm("type", "py")
	path, err := sh.saveFile(fileType, fileName, file)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	print(path)
	// mgm.FilePath = path
	taskId, err := mgm.SubmitTask()
	muslog.Info("submit job with:" + taskId)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if mgm.ScheduleMode == core.ScheduleNow && mgm.DeployMode == core.DeployClient {
		//sh.wsHandler(ctx.Writer, ctx.Request, false)
	}
}

func (sh *StreamingHandler) saveFile(fileType string, fileName string, content multipart.File) (string, error) {
	var path string
	switch fileType {
	case "py":
		path = "/mnt/hgfs/code/trunk/engine/mustang/code/user/python/"
	case "java":
		path = "/mnt/hgfs/code/trunk/engine/mustang/code/user/java/"
	case "scala":
		path = "/mnt/hgfs/code/trunk/engine/mustang/code/user/scala/"
	case "jar":
	default:
		return "", errors.New("file type not support")
	}
	secs := time.Now().Unix()
	name := fmt.Sprintf("%s_%d", fileName, secs)
	filePath := path + name + "." + fileType
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, content)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
