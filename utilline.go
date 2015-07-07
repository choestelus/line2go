package main

import (
	"fmt"
	"line"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"

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

	// Returned headers
	x_ls_header        string
	x_line_application string

	// Client specified user-agent
	userAgent string

	// remember state of command client (first use or not), deal with header
	commandClientState bool
	pollingClientState bool
}

type IcecreamService interface {
	SetHTTPS(bool)
	Login(ident string, pwd string) error
	GetProfile() (line.Profile, error)
	GetAllContactIDs() ([]string, error)
	GetAllGroups() ([]string, error)
	GetMessageHistory(id string) ([]string, error)
	GetAuthToken() string
	GetCertificate() string
	GetOpRevision() int64
	GetX_LSHeader() (string, error)
}

func (this *IcecreamClient) getLoginClient() (client *line.TalkServiceClient) {
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

func (this *IcecreamClient) getCommandClient() (client *line.TalkServiceClient) {
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

func (this *IcecreamClient) getPollingClient() (client *line.TalkServiceClient) {
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
		useHTTPS:   HTTPPrefix,
		commandURL: LineThriftServer + LineCommandPath,
		loginURL:   LineThriftServer + LineLoginPath,
		pollingURL: LineThriftServer + LinePollPath,
		userAgent:  AppUserAgent,
	}
	// TODO: initialize polling service client
	client.LoginClient = client.getLoginClient()

	//commandClient
	client.CommandClient = client.getCommandClient()

	//pollingClient
	client.PollingClient = client.getPollingClient()

	return
}

var cyanBold = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
var greenBold = color.New(color.FgGreen).Add(color.Bold).SprintFunc()

// Login to LINE Server using E-mail as identifier and plaintext password
func LoginLine2(ident string, ptPassword string, client *line.TalkServiceClient) (*line.LoginResult_, error) {
	// Parameters:
	//  - IdentityProvider
	//  - Identifier
	//  - Password
	//  - KeepLoggedIn
	//  - AccessLocation
	//  - SystemName
	//  - Certificate
	result, err := client.LoginWithIdentityCredentialForCertificate(
		line.IdentityProvider_LINE,
		ident,
		ptPassword,
		true,
		"127.0.0.1",
		AppUserAgent,
		"")
	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}

	return result, err
}

// Login to LINE Server using E-mail as identifier and plaintext password
func LoginLine(ident string, ptPassword string) (*line.LoginResult_, error) {
	loginURL := HTTPPrefix + LineThriftServer + LineLoginPath
	loginTransport, err := thrift.NewTHttpPostClient(loginURL)
	if err != nil {
		log.Fatalln("Error Creating Login HTTP Client: ", err)
	}
	loginTrans := loginTransport.(*thrift.THttpClient)

	// set specific header
	loginTrans.SetHeader("User-Agent", AppUserAgent)
	loginTrans.SetHeader("X-Line-Application", LineApplication)
	loginTrans.SetHeader("Connection", "Keep-Alive")

	wrappedLoginTrans := thrift.NewTTransportFactory().GetTransport(loginTrans)

	loginClient := line.NewTalkServiceClientFactory(wrappedLoginTrans, thrift.NewTCompactProtocolFactory())

	// Parameters:
	//  - IdentityProvider
	//  - Identifier
	//  - Password
	//  - KeepLoggedIn
	//  - AccessLocation
	//  - SystemName
	//  - Certificate
	result, err := loginClient.LoginWithIdentityCredentialForCertificate(
		line.IdentityProvider_LINE,
		ident,
		ptPassword,
		true,
		"127.0.0.1",
		AppUserAgent,
		"")
	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}

	return result, err
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

// Returns Command Client to use with first command in initialization Sequence
func GetCommandClient(authToken string) *line.TalkServiceClient {
	commandURL := HTTPPrefix + LineThriftServer + LineCommandPath
	commandTransport, err := thrift.NewTHttpPostClient(commandURL)
	if err != nil {
		log.Fatalln("Error Creating Command HTTP Client: ", err)
	}
	commandTrans := commandTransport.(*thrift.THttpClient)
	commandTrans.SetHeader("X-Line-Access", authToken)
	commandTrans.SetHeader("User-Agent", AppUserAgent)
	commandTrans.SetHeader("X-Line-Application", LineApplication)
	commandTrans.SetHeader("Connection", "Keep-Alive")

	wrappedCommandTrans := thrift.NewTTransportFactory().GetTransport(commandTrans)
	commandClient := line.NewTalkServiceClientFactory(wrappedCommandTrans, thrift.NewTCompactProtocolFactory())

	return commandClient
}

func SetHeaderForClientReuse(client *line.TalkServiceClient, x_ls_header string) {
	// Add X-LS request header, remove other unused headers
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Access")
	client.Transport.(*thrift.THttpClient).DelHeader("User-Agent")
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Application")
	client.Transport.(*thrift.THttpClient).DelHeader("Connection")
	client.Transport.(*thrift.THttpClient).SetHeader("X-LS", x_ls_header)
}
