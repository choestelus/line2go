package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
	// "github.com/parnurzeal/gorequest"
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
	var httpTransport thrift.TTransport
	var message bytes.Buffer

	ident := "choestelus@gmail.com"
	pwd := "suchlinemuchwow@443"

	dev0, _ := os.Open("/dev/zero")

	httpTransport, _ = thrift.NewTHttpPostClient("https://" + LineThriftServer + LineLoginPath)

	// type assertion
	// httptrans.SetHeader("User-Agent", "purple-line (LINE for libpurple/pidgin)")
	httpTrans := httpTransport.(*thrift.THttpClient)
	httpTrans.SetHeader("User-Agent", "purple-line (LINE for libpurple/pidgin)")
	httpTrans.SetHeader("X-Line-Application", "DESKTOPWIN\t3.2.1.83\t\tWINDOWS\t5.1.2600-XP-x64")
	httpTrans.SetHeader("connection", "keep-alive")

	transport = thrift.NewStreamTransport(dev0, &message)
	if err != nil {
		log.Fatalln("Error opening transport connection: ", err)
	}

	transportFactory := thrift.NewTTransportFactory()

	wrappedTransport := transportFactory.GetTransport(transport)
	wrappedhttpTransport := transportFactory.GetTransport(httpTrans)
	var _ = wrappedTransport

	talkClient = line.NewTalkServiceClientFactory(wrappedhttpTransport, thrift.NewTCompactProtocolFactory())

	result, err := talkClient.LoginWithIdentityCredentialForCertificate(line.IdentityProvider_LINE, ident, pwd, true, "127.0.0.1", "purple-line (LINE for libpurple/Pidgin)", "Pidgin")
	if err != nil {
		log.Println("Error logging in: ", err)
	}

	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	log.Printf("body message : %T [%v]\n", message, message.String())
	fmt.Fprintf(ioutil.Discard, "result: %v", result)
	log.Printf("result: %v\n", result)

	yellowBold := color.New(color.FgYellow).Add(color.Bold).SprintfFunc()

	// req := gorequest.New() //.TLSClientConfig(&tls.Config{})
	// resp, body, errs := req.Post("http://"+LineThriftServer+LineLoginPath).
	// 	Set("connection", "keep-alive").
	// 	Set("Content-Type", "application/x-thrift").
	// 	Set("User-Agent", "purple-line (LINE for libpurple/pidgin)").
	// 	Set("X-Line-Application", "DESKTOPWIN\t3.2.1.83\t\tWINDOWS\t5.1.2600-XP-x64").
	// 	Send(message.Bytes()).
	// 	End()
	// if errs != nil {
	// 	log.Fatalf("Fatal : %v", errs)
	// }

	log.Printf("status : %v code = [%v]\n", yellowBold("50x"), yellowBold(strconv.Itoa(509)))
	log.Printf("resp = [%v]\n", "placehold header")
	fmt.Printf("%v %v\nlength = [%v]\n", greenBold("respTLS: "), "placehold TLS", "placehold contentlength")
	fmt.Printf("%v%v\n", greenBold("------body------\n"), "placehold body")

}
