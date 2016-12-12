package core

import (
	"fmt"
	"io"
	"muslog"
	"os/exec"
)

type TaskMgm struct {
}

func NewTaskMgm() *TaskMgm {
	return &TaskMgm{}
}

func (tm *TaskMgm) StartTask(execSeq string) (io.ReadCloser, int, error) {
	muslog.Info("start task with cmd: " + execSeq)
	a := `/opt/spark/bin/spark-submit /opt/spark/examples/src/main/python/streaming/network_wordcount.py 127.0.0.1 9999`
	//	a := "ping 127.0.0.1"
	//a := `/opt/spark/bin/spark-submit --master spark://127.0.0.1:6066 --deploy-mode cluster --class org.apache.spark.examples.streaming.NetworkWordCount /opt/spark/lib/spark-examples-1.6.2-hadoop2.2.0.jar 127.0.0.1 9999 2>&1`
	cmd := exec.Command("/bin/sh", "-c", a)
	//cmd.Stdout = os.Stdout
	r, err := cmd.StdoutPipe()
	if err != nil {
		muslog.Error(err)
		return nil, 0, err
	}
	//out, err := cmd.Output()
	//stdoutPipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		muslog.Error(err)
		return nil, 0, err
	}
	muslog.Info(fmt.Sprintf("pidof task: %d", cmd.Process.Pid))
	return r, cmd.Process.Pid, nil
}

func (tm *TaskMgm) StartClusterTask(execSeq string) (out string, int, error) {
	muslog.Info("start task with cmd: " + execSeq)
	//a := `/opt/spark/bin/spark-submit /opt/spark/examples/src/main/python/streaming/network_wordcount.py 127.0.0.1 9999`
	//	a := "ping 127.0.0.1"
	a := `/opt/spark/bin/spark-submit --master spark://127.0.0.1:6066 --deploy-mode cluster --class org.apache.spark.examples.streaming.NetworkWordCount /opt/spark/lib/spark-examples-1.6.2-hadoop2.2.0.jar 127.0.0.1 9999 2>&1`
	cmd := exec.Command("/bin/sh", "-c", a)
	//cmd.Stdout = os.Stdout
	out, err := cmd.Output()
	//stdoutPipe, _ := cmd.StdoutPipe()
	if err != nil {
		muslog.Error(err)
		return nil, 0, err
	}
	return string(out), cmd.Process.Pid, nil
}
