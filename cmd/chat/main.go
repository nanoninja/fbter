// Copyright 2022 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"

	"github.com/fbuster/chat"
)

func main() {
	var logger = log.Default()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:3000")
	if err != nil {
		logger.Fatalln(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Fatalln(err)
	}

	c := chat.New(logger)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			logger.Println(err)
			continue
		}
		go c.Handle(conn)
	}
}
