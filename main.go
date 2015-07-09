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
	sherbet := NewIcecreamClient()

	fmt.Printf("this is sherbet %T:\n[%v]\n", sherbet, sherbet)
	ident := "choestelus@gmail.com"
	pwd := "suchlinemuchwow@443"

	result, err := sherbet.Login(ident, pwd)
	if err != nil {
		log.Fatalln("Error logging in: ", err)
	}
	// TODO: handle pinverfication request
	if result.GetTypeA1() == line.LoginResultType_REQUIRE_DEVICE_CONFIRM {
		log.Fatalf("error: need pin verification; not handle yet")
	}

	// Initialize new client from received authtoken
	// token = result.GetAuthToken()

	// commandClient := GetCommandClient(token)

	// lastOpRevision, err := commandClient.GetLastOpRevision()
	lastOpRevision, err := sherbet.GetLastOpRevision()
	if err != nil {
		log.Fatalln("Error GetLastOpRevision: ", err)
	}

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)

	printLoginResult(result)
	fmt.Printf("GetLastOpRevision = %v\n", greenBold(strconv.FormatInt(lastOpRevision, 10)))
	fmt.Printf("sherbet.opRevision = %v\n", greenBold(strconv.FormatInt(sherbet.opRevision, 10)))

	// Get X-LS value from response header
	// Line_X_LS = commandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	Line_X_LS = sherbet.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	fmt.Printf("\nX-LS: [%v]\n", cyanBold(Line_X_LS))

	// profile, err := commandClient.GetProfile()
	profile, err := sherbet.GetProfile()
	if err != nil {
		log.Fatalln("Error GetProfile: ", err)
	}
	log.Printf("profile: [%v]\n", cyanBold(profile.String()))

	// allContactIDs, err := commandClient.GetAllContactIds()
	// if err != nil {
	// 	log.Fatalln("Error GetAllContactIds: ", err)
	// }
	// for index, element := range allContactIDs {
	// 	fmt.Fprintf(ioutil.Discard, "#%v: [%v]\n", index, element)
	// }

	// wrapuplist, err := commandClient.GetMessageBoxCompactWrapUpList(1, 50)
	// fmt.Fprintf(ioutil.Discard, "%v\n", wrapuplist.String())
}
