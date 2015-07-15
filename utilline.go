package line2go

import (
	"line2go/linethrift"
	"line2go/thrift"
	"sync"

	"log"
)

type IcecreamClient struct {
	// A bunch of clients
	CommandClient *ConfigurableClient
	LoginClient   *ConfigurableClient
	PollingClient *ConfigurableClient

	// URLs & HTTPS
	useHTTPS   string
	loginURL   string
	commandURL string
	pollingURL string

	// Returned variables
	authToken  string
	opRevision int64

	//polling related
	fetchCount         int32
	fetchResultChannel chan *FetchResult

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

// Type FetchResult is intended to encapsulate result from FetchOperations() for using within channels
type FetchResult struct {
	ops []*line.Operation
	err error
}

type ConfigurableClient struct {
	*line.TalkServiceClient
	headerConfig *sync.Once
}

func (this *IcecreamClient) newLoginClient() (client *ConfigurableClient) {
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
	preclient := line.NewTalkServiceClientFactory(wrappedLoginTrans, thrift.NewTCompactProtocolFactory())

	client = &ConfigurableClient{preclient, new(sync.Once)}
	return
}

func (this *IcecreamClient) newCommandClient() (client *ConfigurableClient) {
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
	preclient := line.NewTalkServiceClientFactory(wrappedCommandTrans, thrift.NewTCompactProtocolFactory())

	client = &ConfigurableClient{preclient, new(sync.Once)}
	return
}

func (this *IcecreamClient) newPollingClient() (client *ConfigurableClient) {
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
	preclient := line.NewTalkServiceClientFactory(wrappedPollingTrans, thrift.NewTCompactProtocolFactory())

	client = &ConfigurableClient{preclient, new(sync.Once)}
	return
}

// Create new instance of client
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

	client.fetchResultChannel = make(chan *FetchResult)
	client.LoginClient = client.newLoginClient()
	client.CommandClient = client.newCommandClient()
	client.PollingClient = client.newPollingClient()

	return
}

// Login with registered E-mail and password
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

// Return current OpRevision value
func (client *IcecreamClient) GetLocalOpRevision() int64 {
	return client.opRevision
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

func SetHeaderForClientReuse(client *ConfigurableClient, x_ls_header string) {
	// Add X-LS request header, remove other unused headers
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Access")
	client.Transport.(*thrift.THttpClient).DelHeader("User-Agent")
	client.Transport.(*thrift.THttpClient).DelHeader("X-Line-Application")
	client.Transport.(*thrift.THttpClient).DelHeader("Connection")
	client.Transport.(*thrift.THttpClient).DelHeader("X-LS")

	client.Transport.(*thrift.THttpClient).SetHeader("X-LS", x_ls_header)
}

func SetHeaderForClientInit(client *ConfigurableClient, authToken string, userAgent string, x_line_application string) {
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
