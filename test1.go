package main

import (
	"os"
	"os/exec"
	//	"time"
)

func run() {
	println("in run")
	//a := `ping 127.0.0.1`
	//cmd := exec.Command("/bin/sh", "-c", a)
	cmd := exec.Command("ping", "127.0.0.1")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//if err := cmd.Wait(); err != nil {
	//	panic(err.Error())
	//}
	println("run exit")
}

func main() {
	run()
	println("in main")
}
