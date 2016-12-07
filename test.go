// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"muslog"
	"time"
)

func main() {
	nw := time.Now().Format("2006-01-02 15:04:05")
	println(nw)
	muslog.InitLog("debug", "", muslog.LevelTrace)
	muslog.Trace("test")
}
