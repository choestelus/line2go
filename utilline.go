package main

import (
	"fmt"
	"line2go/linethrift"
	"line2go/thrift"

	"log"

	"github.com/fatih/color"
)

type IcecreamClient struct {
	// A bunch of clients
	CommandClient *line.TalkServiceClient
	LoginClient   *line.TalkServiceClient
	PollingClient *line.TalkServiceClient

	// URLs & HTTPS
	useHTTPS   string
	loginURL   string
	commandURL string
	pollingURL string

	// Returned variables
	authToken  string
	opRevision int64

	fetchCount int32

	// Returned headers
	cx_ls_header       string
	px_ls_header       string
	x_line_application string

	// Client specified user-agent
	userAgent string

	// remember state of command client (first use or not), deal with header
	commandClientState bool
	pollingClientState bool
}

func (this *IcecreamClient) NewLoginClient() (client *line.TalkServiceClient) {
	// Assuming URL is sanitized
	loginURL := this.useHTTPS + this.loginURL
	loginTransport, err := thrift.NewTHttpPostClient(loginURL)
	if err != nil {
		log.Fatalln("error creating loginTransport: ", err)
	}

	loginTrans := loginTransport.(*thrift.THttpClient)

	loginTrans.SetHeader("User-Agent", AppUserAgent)
	loginTrans.SetHeader("X-Line-Application", LineApplication)
	loginTrans.SetHeader("Connection", "Keep-Alive")

	wrappedLoginTrans := thrift.NewTTransportFactory().GetTransport(loginTrans)
	client = line.NewTalkServiceClientFactory(wrappedLoginTrans, thrift.NewTCompactProtocolFactory())

	return client
}

func (this *IcecreamClient) NewCommandClient() (client *line.TalkServiceClient) {
	// Assuming URL is sanitized
	commandURL := this.useHTTPS + this.commandURL
	commandTransport, err := thrift.NewTHttpPostClient(commandURL)
	if err != nil {
		log.Fatalln("error creating commandTransport: ", err)
	}

	commandTrans := commandTransport.(*thrift.THttpClient)
	commandTrans.SetHeader("X-Line-Access", this.authToken)
	commandTrans.SetHeader("User-Agent", this.userAgent)
	commandTrans.SetHeader("X-Line-Application", this.x_line_application)
	commandTrans.SetHeader("Connection", "Keep-Alive")

	wrappedCommandTrans := thrift.NewTTransportFactory().GetTransport(commandTrans)
	client = line.NewTalkServiceClientFactory(wrappedCommandTrans, thrift.NewTCompactProtocolFactory())

	return
}

func (this *IcecreamClient) NewPollingClient() (client *line.TalkServiceClient) {
	// Assuming URL is sanitized
	pollingURL := this.useHTTPS + this.pollingURL
	pollingTransport, err := thrift.NewTHttpPostClient(pollingURL)
	if err != nil {
		log.Fatalln("error creating pollingTransport: ", err)
	}

	pollingTrans := pollingTransport.(*thrift.THttpClient)
	pollingTrans.SetHeader("X-Line-Access", this.authToken)
	pollingTrans.SetHeader("User-Agent", this.userAgent)
	pollingTrans.SetHeader("X-Line-Application", this.x_line_application)
	pollingTrans.SetHeader("Connection", "Keep-Alive")

	wrappedPollingTrans := thrift.NewTTransportFactory().GetTransport(pollingTrans)
	client = line.NewTalkServiceClientFactory(wrappedPollingTrans, thrift.NewTCompactProtocolFactory())

	return
}

func NewIcecreamClient() (client *IcecreamClient) {
	// TODO: decoupling from global constant
	client = &IcecreamClient{
		useHTTPS:           HTTPPrefix,
		commandURL:         LineThriftServer + LineCommandPath,
		loginURL:           LineThriftServer + LineLoginPath,
		pollingURL:         LineThriftServer + LinePollPath,
		userAgent:          AppUserAgent,
		x_line_application: LineApplication,
		fetchCount:         DefaultFetchCount,
	}

	client.LoginClient = client.NewLoginClient()
	client.CommandClient = client.NewCommandClient()
	client.PollingClient = client.NewPollingClient()

	return
}

var cyanBold = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
var greenBold = color.New(color.FgGreen).Add(color.Bold).SprintFunc()

func (client *IcecreamClient) Login(ident string, ptpwd string) (result *line.LoginResult_, err error) {
	// Parameters:
	//  - IdentityProvider
	//  - Identifier
	//  - Password
	//  - KeepLoggedIn
	//  - AccessLocation
	//  - SystemName
	//  - Certificate
	result, err = client.LoginClient.LoginWithIdentityCredentialForCertificate(
		line.IdentityProvider_LINE,
		ident,
		ptpwd,
		true,
		"127.0.0.1",
		AppUserAgent,
		"")
	if err != nil || result.GetTypeA1() != line.LoginResultType_SUCCESS {
		return
	}
	log.Println("token from login: ", result.GetAuthToken())
	client.authToken = result.GetAuthToken()
	return
}

func (client *IcecreamClient) GetAuthToken() string {
	return client.authToken
}

func setState(state *bool) {
	if *state != true {
		*state = true
	}
}

func (client *IcecreamClient) setCommandState() {
	client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	setState(&client.commandClientState)
}

func (client *IcecreamClient) setPollingState() {
	client.px_ls_header = client.PollingClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	setState(&client.pollingClientState)
}

// Print formatted *line.LoginResult_
func printLoginResult(result *line.LoginResult_) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("token: [%v]\n", cyanBold(result.GetAuthToken()))
	fmt.Printf("certificate: [%v]\n", cyanBold(result.GetCertificate()))
	fmt.Printf("pincode: [%v]\n", cyanBold(result.GetPinCode()))
	fmt.Printf("loginResult: [%v]\n", cyanBold(result.GetTypeA1().String()))
	fmt.Printf("verifier: [%v]\n", cyanBold(result.GetVerifier()))
	fmt.Println("--------------------------------------------------------------------------------")
}

func SetHeaderForClientReuse(client *line.TalkServiceClient, x_ls_header string) {
	// Add X-LS request header, remove other unused headers
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Access")
	client.Transport.(*thrift.THttpClient).DelHeader("User-Agent")
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Application")
	client.Transport.(*thrift.THttpClient).DelHeader("Connection")
	client.Transport.(*thrift.THttpClient).DelHeader("X-LS")

	client.Transport.(*thrift.THttpClient).SetHeader("X-LS", x_ls_header)
}

func SetHeaderForClientInit(client *line.TalkServiceClient, authToken string, userAgent string, x_line_application string) {
	client.Transport.(*thrift.THttpClient).DelHeader("X-LS")
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Access")
	client.Transport.(*thrift.THttpClient).DelHeader("User-Agent")
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Application")
	client.Transport.(*thrift.THttpClient).DelHeader("Connection")

	client.Transport.(*thrift.THttpClient).SetHeader("X-Line-Access", authToken)
	client.Transport.(*thrift.THttpClient).SetHeader("User-Agent", userAgent)
	client.Transport.(*thrift.THttpClient).SetHeader("X-Line-Application", x_line_application)
	client.Transport.(*thrift.THttpClient).SetHeader("Connection", "Keep-Alive")
}
