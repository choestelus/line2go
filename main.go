package main

import (
	"fmt"
	"io/ioutil"
	"line"
	"log"
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/fatih/color"
)

var (
	token     string
	Line_X_LS string
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
	loginTrans.SetHeader("Connection", "Keep-Alive")

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

	// Initialize new client from received authtoken
	token = result.GetAuthToken()
	commandTrans.SetHeader("X-Line-Access", token)
	commandTrans.SetHeader("User-Agent", AppUserAgent)
	commandTrans.SetHeader("X-Line-Application", LineApplication)
	commandTrans.SetHeader("Connection", "Keep-Alive")

	wrappedCommandTransport := transportFactory.GetTransport(commandTrans)
	commandClient := line.NewTalkServiceClientFactory(wrappedCommandTransport, thrift.NewTCompactProtocolFactory())

	lastOpRevision, err := commandClient.GetLastOpRevision()
	if err != nil {
		log.Fatalln("Error GetLastOpRevision: ", err)
	}

	// proof that GetHeader is useless to get header from HTTP response messages.
	// fmt.Printf("\ntest get-header: [%v]\n", loginClient.Transport.(*thrift.THttpClient).GetHeader("X-Line-Application"))
	// fmt.Printf("\ntest get-header: [%v]\n", loginClient.Transport.(*thrift.THttpClient).GetHeader("X-Lcr"))

	// Workaround: use this instead
	// disclamier: น่าจะ Non thread-safe if there are more than 1 go-routine call using loginClient instance

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)

	printLoginResult(result)
	fmt.Printf("GetLastOpRevision = %v\n", greenBold(strconv.FormatInt(lastOpRevision, 10)))
	Line_X_LS = commandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	fmt.Printf("\nX-LS: [%v]\n", cyanBold(Line_X_LS))

	// TODO: add X-LS request header, remove other unused headers
	commandClient.Transport.(*thrift.THttpClient).DelHeader("X-Line-Access")
	commandClient.Transport.(*thrift.THttpClient).DelHeader("User-Agent")
	commandClient.Transport.(*thrift.THttpClient).DelHeader("X-Line-Application")
	commandClient.Transport.(*thrift.THttpClient).DelHeader("Connection")

	commandClient.Transport.(*thrift.THttpClient).SetHeader("X-LS", Line_X_LS)

	// TODO: GetProfile
	profile, err := commandClient.GetProfile()
	if err != nil {
		log.Fatalln("Error GetProfile: ", err)
	}
	log.Printf("profile: [%v]\n", cyanBold(profile.String()))

	// TODO: GetAllContactIds
	allContactIDs, err := commandClient.GetAllContactIds()
	if err != nil {
		log.Fatalln("Error GetAllContactIds: ", err)
	}
	for index, element := range allContactIDs {
		fmt.Printf("#%v: [%v]\n", index, element)
	}

	// TODO: handle pinverfication request
	if result.GetTypeA1() == line.LoginResultType_REQUIRE_DEVICE_CONFIRM {
		log.Fatalf("error: need pin verification; not handle yet")
		// Code here:
		// create blank http request
	}
}
