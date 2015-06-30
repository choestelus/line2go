package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
	"io/ioutil"
	"line"
	"log"
)

var (
	token string
)

func main() {
	fmt.Fprintf(ioutil.Discard, "")
	greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	var err error
	var loginClient *line.TalkServiceClient
	var loginTransport thrift.TTransport
	var commandTransport thrift.TTransport

	ident := "choestelus@gmail.com"
	pwd := "suchlinemuchwow@443"

	url := "https://" + LineThriftServer + LineLoginPath
	url2 := "https://" + LineThriftServer + LineCommandPath
	loginTransport, _ = thrift.NewTHttpPostClient(url)
	commandTransport, _ = thrift.NewTHttpPostClient(url2)

	// type assertion
	loginTrans := loginTransport.(*thrift.THttpClient)
	commandTrans := commandTransport.(*thrift.THttpClient)

	// set specific header
	loginTrans.SetHeader("User-Agent", AppUserAgent)
	loginTrans.SetHeader("X-Line-Application", LineApplication)
	loginTrans.SetHeader("connection", "keep-alive")

	transportFactory := thrift.NewTTransportFactory()
	wrappedloginTransport := transportFactory.GetTransport(loginTrans)

	loginClient = line.NewTalkServiceClientFactory(
		wrappedloginTransport,
		thrift.NewTCompactProtocolFactory())

	result, err := loginClient.LoginWithIdentityCredentialForCertificate(
		line.IdentityProvider_LINE,
		ident,
		pwd,
		true,
		"127.0.0.1",
		AppUserAgent,
		"Pidgin")

	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}
	// proof that GetHeader is useless to get header from HTTP response messages.
	//fmt.Printf("\ntest get-header: [%v]\n", loginClient.Transport.(*thrift.THttpClient).GetHeader("X-Line-Application"))
	//fmt.Printf("\ntest get-header: [%v]\n", loginClient.Transport.(*thrift.THttpClient).GetHeader("X-Lcr"))

	// Workaround: use this instead
	// disclamier: น่าจะ Non thread-safe if there are more than 1 go-routine call using loginClient instance
	fmt.Printf("\ntest get-header: [%v]\n", loginClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-Lcr"))

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)
	printLoginResult(result)

	//TODO: handle pinverfication request
	if result.GetTypeA1() == line.LoginResultType_REQUIRE_DEVICE_CONFIRM {
		log.Fatalf("error: need pin verification; not handle yet")
		// Code here:
		// create blank http request
	}
}
