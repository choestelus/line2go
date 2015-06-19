package main

import (
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"line"
	"log"
)

var _ = thrift.ZERO
var _ = fmt.Printf

type ThriftClient struct {
	*line.TalkServiceClient
}

func main() {
	// var err error
	var transport thrift.TTransport
	var talkClient *line.TalkServiceClient
	// var talkClient ThriftClient
	urlPath := LineThriftServer + ":443" // + LineLoginPath

	ident := "choestelus@gmail.com"
	pwd := ""
	cfg := new(tls.Config)
	log.Printf("openning socket to: %v", urlPath)
	transport, err2 := thrift.NewTSSLSocket(urlPath, cfg)
	if err2 != nil {
		log.Fatalln("Error opening secure socket : ", err2)
	}

	transportFactory := thrift.NewTTransportFactory()
	transport = transportFactory.GetTransport(transport)
	if err := transport.Open(); err != nil {
		log.Fatalln("could not open connection", err)
	}
	defer transport.Close()

	talkClient = line.NewTalkServiceClientFactory(transport, thrift.NewTCompactProtocolFactory())

	result, err := talkClient.LoginWithIdentityCredentialForCertificate(line.IdentityProvider_LINE, ident, pwd, true, "127.0.0.1", "goLINE", "emp")
	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}

	log.Printf("result: %v", result)
}
