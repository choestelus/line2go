package main

import (
	"fmt"
	"io/ioutil"
	"line2go/linethrift"
	"line2go/thrift"

	"log"
	"os"
	"strconv"
	"strings"
)

var ()

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

	lastOpRevision, err := sherbet.GetLastOpRevision()
	if err != nil {
		log.Fatalln("Error GetLastOpRevision: ", err)
	}

	prettyResult := fmt.Sprint(greenBold(result.String()))
	log.Printf("Type: [%T], result: %v\n", result, prettyResult)

	printLoginResult(result)
	fmt.Printf("GetLastOpRevision = %v\n", greenBold(strconv.FormatInt(lastOpRevision, 10)))
	fmt.Printf("sherbet.opRevision = %v\n", greenBold(strconv.FormatInt(sherbet.opRevision, 10)))

	Line_X_LS := sherbet.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
	fmt.Printf("\nX-LS: [%v]\n", cyanBold(Line_X_LS))

	profile, err := sherbet.GetProfile()
	if err != nil {
		log.Fatalln("Error GetProfile: ", err)
	}
	log.Printf("profile: [%v]\n", cyanBold(profile.String()))

	allContactIDs, err := sherbet.GetAllContactIDs()
	if err != nil {
		log.Fatalln("Error GetAllContactIds: ", err)
	}
	for index, element := range allContactIDs {
		if index == len(allContactIDs)-1 {
			fmt.Fprintf(os.Stdout, "#%v: [%v]\n", index, element)
		}
	}

	contacts, err := sherbet.GetContacts(allContactIDs)
	if err != nil {
		log.Fatalln("Error GetContacts: ", err)
	}
	fmt.Println(greenBold("contacts: "), contacts[len(contacts)-3:len(contacts)-1])

	for index, element := range contacts {
		if strings.Contains(element.GetDisplayName(), "iko") {
			fmt.Printf("\n\n #%v ID of %v It's : [%v]\n\n", index, element.GetDisplayName(), element.GetMid())
		}
	}

	contact, err := sherbet.GetContact(allContactIDs[0])
	if err != nil {
		log.Fatalln("Error Get Contact: ", err)
	}
	fmt.Println("contact: [", contact.GetStatus(), contact.GetDisplayName(), "]")

	groupsJoined, err := sherbet.GetGroupIdsJoined()
	if err != nil {
		log.Fatalln("Error GetGroupIdsJoined(): ", err)
	}
	fmt.Println(greenBold("Groups Joined: "), groupsJoined)

	groupsInvited, err := sherbet.GetGroupIdsInvited()
	if err != nil {
		log.Fatalln("Error GetGroupIdsInvited(): ", err)
	}
	fmt.Println(greenBold("Groups Invited: "), groupsInvited)
	groups, err := sherbet.GetGroups(groupsJoined)
	for index, element := range groups {
		if index == len(groups)-1 {
			fmt.Printf("#%v: [%v]\n", cyanBold(index), element)
		}
	}
	group, err := sherbet.GetGroup(groupsJoined[0])
	fmt.Println(greenBold("single group: "), group.String())

	msgboxl, err := sherbet.GetMessageBoxCompactWrapUpList(1, 50)
	if err != nil {
		log.Fatalln("Error GetMessageBoxCompactWrapUpList: ", err)
	}
	fmt.Println(greenBold("msgboxl: "), "[", msgboxl.GetMessageBoxWrapUpList()[0].String(), "]")

	fmt.Println("--------------------------------------------------------------------------------")
	try_fetch, err := sherbet.FetchOperations()
	if err != nil {

		if err.Error() == "HTTP Response code: 410" {
			log.Printf("410 gone: re-requesing...\n")
		} else {
			log.Printf("Error [%v]\n", err.Error())
		}
	}

	// Try to fetch
	for i := 0; i < 10; i++ {
		fmt.Println("--------------------------------------------------------------------------------")
		try_fetch, err = sherbet.FetchOperations()
		if err != nil {
			if err.Error() == "HTTP Response code: 410" {
				log.Printf("410 gone: re-requesing...\n")
			} else {
				log.Printf("Error [%v]\n", err.Error())
			}
		}
		fmt.Printf("Fetch Result: [%v]\n", try_fetch)
	}
	// rmsg, err := sherbet.SendTextMessage("ue2af231f5fe993dda7051b816d072c2c", "สวัสดี ภาษา go ก็รับ Unicode นะ :=")
	// if err != nil {
	// 	log.Fatalln("Error Sending Message", err)
	// }
	// fmt.Printf("Message sent\nID: [%v]\nothers: [%v]\n", cyanBold(rmsg.GetId()), cyanBold(rmsg.GetCreatedTime()))
}
