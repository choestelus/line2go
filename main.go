package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"line"
	"log"
	"os"
	"strconv"
)

type ThriftClient struct {
	*line.TalkServiceClient
}

var _ = tls.Config{}

func main() {
	var err error
	var transport *thrift.StreamTransport
	var talkClient *line.TalkServiceClient
	var message bytes.Buffer

	ident := "choestelus@gmail.com"
	pwd := ""

	dev0, _ := os.Open("/dev/zero")

	transport = thrift.NewStreamTransport(dev0, &message)
	if err != nil {
		log.Fatalln("Error opening transport connection: ", err)
	}

	transportFactory := thrift.NewTTransportFactory()
	wrappedTransport := transportFactory.GetTransport(transport)

	talkClient = line.NewTalkServiceClientFactory(wrappedTransport, thrift.NewTCompactProtocolFactory())

	result, err := talkClient.LoginWithIdentityCredentialForCertificate(line.IdentityProvider_LINE, ident, pwd, true, "127.0.0.1", "goLINE", "emp")
	if err != nil {
		log.Println("Error logging in: ", err)
	}

	log.Printf("body message : %T [%v]\n", message, message.String())
	fmt.Fprintf(ioutil.Discard, "result: %v", result)

	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	req := gorequest.New() //.TLSClientConfig(&tls.Config{})
	resp, body, errs := req.Post("http://"+LineThriftServer+LineLoginPath).
		Set("connection", "keep-alive").
		Set("Content-Type", "application/x-thrift").
		Set("User-Agent", "purple-line (LINE for libpurple/pidgin)").
		Set("X-Line-Application", "DESKTOPWIN\t3.2.1.83\tWINDOWS\t5.1.2600-XP-x64").
		SendString(message.String()).
		End()
	if errs != nil {
		log.Fatalf("Fatal : %v", errs)
	}

	log.Printf("status : %v code = [%v]\n", greenBold(resp.Status), greenBold(strconv.Itoa(resp.StatusCode)))
	log.Printf("resp = [%v]\n", resp.Header)
	fmt.Printf("%v %v\nlength = [%v]\n", greenBold("respTLS: "), resp.TLS, resp.ContentLength)
	fmt.Printf("%v%v\n", greenBold("------body------\n"), body)

}
