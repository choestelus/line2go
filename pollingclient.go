package main

import "line2go/linethrift"

func (client *IcecreamClient) FetchOperations(localRev int64, count int32) (op []*line.Operation, err error) {

	// begin before section
	if client.pollingClientState == true {
		SetHeaderForClientReuse(client.PollingClient, client.px_ls_header)
	} else {
		SetHeaderForClientInit(client.PollingClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section
	op, err = client.PollingClient.FetchOperations(client.opRevision, client.fetchCount)

	// begin after section
	client.setPollingState()
	// end after section
	return
}
