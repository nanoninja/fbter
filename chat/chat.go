// Copyright 2022 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const MsgDateLayout = "15:04:05"

type User struct {
	Name     string
	Addr     net.Addr
	Outgoing chan Message
}

func NewUser() User {
	return User{
		Outgoing: make(chan Message),
	}
}

type Message struct {
	Username string
	Text     string
	Date     time.Time
}

func NewMessage(username string, text string) Message {
	return Message{
		Username: username,
		Text:     text,
		Date:     time.Now(),
	}
}

type Chat struct {
	join   chan User
	leave  chan User
	input  chan Message
	users  map[string]User
	logger *log.Logger
}

func New(logger *log.Logger) *Chat {
	c := &Chat{
		join:   make(chan User),
		leave:  make(chan User),
		input:  make(chan Message),
		users:  make(map[string]User),
		logger: logger,
	}
	if c.logger == nil {
		c.logger = log.Default()
	}

	go c.Run()

	return c
}

func (c *Chat) Handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Println(err)
		}
	}()
	scanner := bufio.NewScanner(conn)
	user, err := c.Login(conn, scanner)
	if err != nil {
		return
	}
	defer c.Logout(user)
	c.Join(user)

	go c.Read(user, conn)

	for m := range user.Outgoing {
		c.WriteMessage(conn, m)
	}
}

func (c *Chat) Join(user User) {
	c.join <- user
}

func (c *Chat) Emit(m Message) {
	c.input <- m
}

func (c *Chat) Broadcast(m Message) {
	for _, user := range c.users {
		if m.Username != user.Name {
			user.Outgoing <- m
		}
	}
}

func (c *Chat) Read(u User, rw io.ReadWriteCloser) {
	r := bufio.NewReader(rw)

	for {
		text, err := r.ReadString('\n')
		if err == io.EOF {
			c.logger.Println(err)
			c.Logout(u)
			return
		}
		if err != nil {
			c.logger.Println(err)
			return
		}
		if strings.TrimSpace(text) == "quit" {
			c.Logout(u)
			rw.Close()
			return
		}
		c.Emit(NewMessage(u.Name, text))
	}
}

func (c *Chat) Logout(user User) {
	c.leave <- user
}

func (c *Chat) Login(conn net.Conn, scanner *bufio.Scanner) (user User, err error) {
	user = NewUser()

	for {
		_, err = c.WriteString(conn, "Enter your name: ")
		if err != nil {
			break
		}
		scanner.Scan()

		if scanner.Text() == "" {
			continue
		} else if _, ok := c.users[scanner.Text()]; ok {
			_, err = c.WriteString(conn, "Username already exists \n")
			if err != nil {
				c.logger.Println(err)
				return
			}
			continue
		} else {
			user.Addr = conn.RemoteAddr()
			user.Name = scanner.Text()
			break
		}
	}
	return
}

func (c *Chat) Add(u User) {
	c.users[u.Name] = u
}

func (c *Chat) Remove(u User) {
	delete(c.users, u.Name)
}

func (c *Chat) Run() {
	for {
		select {
		case user := <-c.join:
			c.Add(user)

			go c.Emit(NewMessage("bot", fmt.Sprintf("%s joined \n", user.Name)))

		case user := <-c.leave:
			c.Remove(user)

			go c.Emit(NewMessage("bot", fmt.Sprintf("%s left \n", user.Name)))

		case m := <-c.input:
			c.Broadcast(m)
			c.LogMessage(os.Stdout, m)
		}
	}
}

func (c *Chat) LogMessage(w io.Writer, m Message) (int, error) {
	f := fmt.Sprintf("[%s] %s> %s", m.Date.Format(MsgDateLayout), m.Username, m.Text)

	if u, ok := c.users[m.Username]; ok {
		f = fmt.Sprintf("[%s %s] %s> %s", u.Addr, m.Date.Format(MsgDateLayout), u.Name, m.Text)
	}
	return c.WriteString(os.Stdout, f)
}

func (c *Chat) WriteString(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s)
}

func (c *Chat) WriteMessage(w io.Writer, m Message) (int, error) {
	s := fmt.Sprintf("%s> %s: %s", m.Username, m.Date.Format(MsgDateLayout), m.Text)
	return c.WriteString(w, s)
}
