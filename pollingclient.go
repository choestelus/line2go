package main

import "line2go/linethrift"

func (client *IcecreamClient) FetchOperations() (op []*line.Operation, err error) {

	// begin before section
	if client.pollingClientState == true {
		SetHeaderForClientReuse(client.PollingClient, client.px_ls_header)
	} else {
		SetHeaderForClientInit(client.PollingClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section
	op, err = client.PollingClient.FetchOperations(client.opRevision, client.fetchCount)
	if err != nil {
		return
	}
	for _, element := range op {
		if client.opRevision < element.GetRevision() {
			client.opRevision = element.GetRevision()
		}
	}

	// begin after section
	client.setPollingState()
	// end after section
	return
}
