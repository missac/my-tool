package core

import (
	"dbmgm"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"muslog"
	"os"
	"os/exec"
	"path/filepath"
	"tool"
)

type StreamingMgm struct {
	AppName            string
	TaskType           string // streming or sql or ml
	TaskID             string // identify task in databse
	TaskDes            string // streming or sql or ml
	Master             string // spark  yarn  messos  or local
	DeployMode         string // client or cluster
	ClassName          string // for java and scala
	Jars               string // jars path for java and scala
	PyFile             string //for python
	DriverMem          string
	ExecutorMem        string
	DriverCores        string // spark standalone and yarn with cluster mode
	TotalExecutorCores string // spark standalone and messos only
	ExecutorCores      string //spark standalone and yarn only
	ExecutorNums       string //yarn only

	Parallelism           string //spark.default.parallelism num- executors * executor-cores的2~3倍较为合适
	StorgeMemoryFraction  string //spark.storage.memoryFraction default 0.6
	ShuffleMemoryFraction string // spark.shuffle.memoryFraction default 0.2

	DataSource   string // kafka, flume, socket
	CodeType     string // java, python scala
	ScheduleMode string // TODO default

	ScheduleSeq string
	FilePath    string
	ExecSeq     string
}

func NewStreamingMgm() *StreamingMgm {
	return &StreamingMgm{}
}

func (s *StreamingMgm) GenTemplate() (string, error) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	print(dir)
	path := "/mnt/hgfs/code/trunk/engine/mustang/code/template/" +
		s.CodeType + "/" + s.DataSource + "_template.py"
	muslog.Info("get template from file" + path)
	return s.genTemplate(path)
}

func (s *StreamingMgm) SubmitTask() (string, error) {
	var cmd = "/opt/spark/bin/spark-submit"
	cmd += " --name " + s.AppName

	switch {
	case s.CodeType == LanPython:
		cmd += " --py-files " + s.FilePath
	case s.CodeType == LanJava:
	case s.CodeType == LanScala:
		cmd += " --jars " + s.Jars
	default:
		errors.New("code type not support")
	}

	if s.DriverMem != "1G" {
		cmd += " --driver-memory " + s.DriverMem
	}
	switch s.Master {
	case StandAlone:
		cmd += " --master " + s.Master + "://192.168.1.241:7077"
		if s.DeployMode == DeployCluster && s.DriverCores != "1" {
			cmd += " --driver-cores " + s.DriverCores
		}
		if s.TotalExecutorCores != "0" {
			cmd += " --total-executor-cores " + s.TotalExecutorCores
		}
		if s.ExecutorCores != "0" {
			cmd += " --executor-cores " + s.ExecutorCores
		}
	case Yarn:
		// TODO --master
		if s.ExecutorCores != "0" {
			cmd += " --executor-cores " + s.ExecutorCores
		}
		if s.DeployMode == DeployCluster && s.DriverCores != "1" {
			cmd += " --driver-cores " + s.DriverCores
		}
		if s.ExecutorNums != "2" {
			cmd += " --num-executors " + s.ExecutorNums
		}
	case Messos:
		cmd += " -- master messos://192.168.1.241 "
		if s.TotalExecutorCores != "0" {
			cmd += " --total-executor-cores " + s.TotalExecutorCores
		}
	default:
		muslog.Error("master type not support:" + s.Master)
		return "", errors.New("master not support")
	}

	if s.Parallelism != "0" {
		cmd += " --conf spark.default.parallelism=" + s.Parallelism
	}
	cmd += " --conf spark.storage.memoryFraction=" + s.StorgeMemoryFraction
	cmd += " --conf spark.shuffle.memoryFraction=" + s.ShuffleMemoryFraction
	muslog.Info("spark submit sequence: " + cmd)
	s.ExecSeq = cmd
	var execTimes = 0
	if s.ScheduleMode == ScheduleNow {
		execTimes = 1
	}
	s.TaskID = tool.GenUUID()
	muslog.Info("add task with taskid: " + s.TaskID)

	err := dbmgm.AddTask(s.AppName, s.TaskID, s.TaskType, s.TaskDes, s.CodeType, s.ScheduleMode, execTimes, s.DataSource, cmd, s.FilePath, s.DeployMode, "notset")
	if err != nil {
		muslog.Error("add task faild: " + err.Error())
		return "", err
	}
	return s.TaskID, nil
}

func (s *StreamingMgm) StartJob() (io.ReadCloser, int, error) {
	muslog.Info("start job")
	a := `/opt/spark/bin/spark-submit /opt/spark/examples/src/main/python/streaming/network_wordcount.py 127.0.0.1 9999`
	//	a := "ping 127.0.0.1"
	//a := `/opt/spark/bin/spark-submit --master spark://127.0.0.1:6066 --deploy-mode cluster --class org.apache.spark.examples.streaming.NetworkWordCount /opt/spark/lib/spark-examples-1.6.2-hadoop2.2.0.jar 127.0.0.1 9999 2>&1`
	cmd := exec.Command("/bin/sh", "-c", a)
	muslog.Info("start task with cmd: " + s.ExecSeq)
	//cmd.Stdout = os.Stdout
	r, err := cmd.StdoutPipe()
	if err != nil {
		muslog.Error(err)
		return nil, 0, err
	}
	//out, err := cmd.Output()
	//stdoutPipe, _ := cmd.StdoutPipe()
	cmd.Start()
	muslog.Info(fmt.Sprintf("pidof task: %d", cmd.Process.Pid))
	return r, cmd.Process.Pid, nil
}

func (s *StreamingMgm) genTemplate(file string) (string, error) {
	fi, err := os.Open(file)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		muslog.Error(err)
		return "", err
	}
	return string(fd), nil
}
