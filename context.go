package main

import (
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bslatkin/go-signin-with-twitter"
)

type Context struct {
	RequestId  string
	ScreenName string
}

func (c *Context) Log(format string, args ...interface{}) {
	log.Printf(c.RequestId+": "+format, args...)
}

type Handler func(*Context, http.ResponseWriter, *http.Request)

func Register(url string, h Handler) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		c := Context{
			RequestId: uuid.New(),
		}
		if u, _ := signin.GetUserInfo(r); u != nil {
			c.ScreenName = u.ScreenName
		}
		h(&c, w, r)
	})
}
