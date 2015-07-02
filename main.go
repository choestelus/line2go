package main

import (
	"fmt"
	"io/ioutil"
	"line"
	"log"
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	token     string
	Line_X_LS string
)

func main() {
	fmt.Fprintf(ioutil.Discard, "")
	var err error

	ident := "choestelus@gmail.com"
	pwd := "suchlinemuchwow@443"

	// type assertion

	result, err := LoginLine(ident, pwd)
	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}
	// TODO: handle pinverfication request
	if result.GetTypeA1() == line.LoginResultType_REQUIRE_DEVICE_CONFIRM {
		log.Fatalf("error: need pin verification; not handle yet")
		// Code here:
		// create blank http request
	}

	// Initialize new client from received authtoken
	token = result.GetAuthToken()

	// commandClient := line.NewTalkServiceClientFactory(wrappedCommandTransport, thrift.NewTCompactProtocolFactory())
	commandClient := GetCommandClient(token)

	lastOpRevision, err := commandClient.GetLastOpRevision()
	if err != nil {
		log.Fatalln("Error GetLastOpRevision: ", err)
	}

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)

	printLoginResult(result)
	fmt.Printf("GetLastOpRevision = %v\n", greenBold(strconv.FormatInt(lastOpRevision, 10)))

	// Get X-LS value from response header
	Line_X_LS = commandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	fmt.Printf("\nX-LS: [%v]\n", cyanBold(Line_X_LS))

	SetHeaderForClientReuse(commandClient, Line_X_LS)

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
		fmt.Fprintf(ioutil.Discard, "#%v: [%v]\n", index, element)
	}

	// TODO: GetMessageBoxCompactWrapUpList
	wrapuplist, err := commandClient.GetMessageBoxCompactWrapUpList(1, 50)
	fmt.Printf("%v\n", wrapuplist.String())
}
