package main

import (
	"bytes"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
)

// suppress compilation error
var (
	_ = fmt.Fprintf
	_ = http.DefaultMaxHeaderBytes
)

func createHTTPPacket(x_ls string) *http.Request {
	var b bytes.Buffer
	req, err := http.NewRequest("POST", LineThriftServer, &b)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("X-LS", x_ls)
	req.Header.Add("X-Line-Application", LineApplication)
	return req
}

func createHTTPPacket2() {
	req := gorequest.New()
	req.End()
}
