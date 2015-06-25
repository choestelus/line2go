package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
	"io/ioutil"
	"line"
	"log"
)

func main() {
	fmt.Fprintf(ioutil.Discard, "")
	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	var err error
	var talkClient *line.TalkServiceClient
	var httpTransport thrift.TTransport

	ident := "choestelus@gmail.com"
	pwd := "suchlinemuchwow@443"

	httpTransport, _ = thrift.NewTHttpPostClient("https://" + LineThriftServer + LineLoginPath)

	// type assertion
	httpTrans := httpTransport.(*thrift.THttpClient)

	// set specific header
	httpTrans.SetHeader("User-Agent", AppUserAgent)
	httpTrans.SetHeader("X-Line-Application", LineApplication)
	httpTrans.SetHeader("connection", "keep-alive")

	transportFactory := thrift.NewTTransportFactory()
	wrappedhttpTransport := transportFactory.GetTransport(httpTrans)

	talkClient = line.NewTalkServiceClientFactory(wrappedhttpTransport, thrift.NewTCompactProtocolFactory())

	result, err := talkClient.
		LoginWithIdentityCredentialForCertificate(line.IdentityProvider_LINE, ident, pwd, true, "127.0.0.1", AppUserAgent, "Pidgin")
	if err != nil {
		log.Println("Error logging in: ", err)
	}

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)
	printLoginResult(result)
}
