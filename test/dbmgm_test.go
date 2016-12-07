package gotest

import (
	"dbmgm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initDB(t *testing.T) {
	dbmgm.InitDB("./mustang.db")
}

func Test_addTask(t *testing.T) {
	assert.NotEqual(t, a, b, "a not equal with b")
}

func Test_listStreamingTask(t *testing.T) {
	res, _ := dbmgm.ListStreamingTask()
	assert.NotEmpty(t, res)
}
