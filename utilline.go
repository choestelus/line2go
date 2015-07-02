package main

import (
	"fmt"
	"line"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/fatih/color"
)

var cyanBold = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
var greenBold = color.New(color.FgGreen).Add(color.Bold).SprintFunc()

func LoginLine(ident string, ptPassword string) (*line.LoginResult_, error) {
	loginURL := HTTPPrefix + LineThriftServer + LineLoginPath
	loginTransport, err := thrift.NewTHttpPostClient(loginURL)
	if err != nil {
		log.Fatalln("Error Creating HTTP Client: ", err)
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

func printLoginResult(result *line.LoginResult_) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("token: [%v]\n", cyanBold(result.GetAuthToken()))
	fmt.Printf("certificate: [%v]\n", cyanBold(result.GetCertificate()))
	fmt.Printf("pincode: [%v]\n", cyanBold(result.GetPinCode()))
	fmt.Printf("loginResult: [%v]\n", cyanBold(result.GetTypeA1().String()))
	fmt.Printf("verifier: [%v]\n", cyanBold(result.GetVerifier()))
}
