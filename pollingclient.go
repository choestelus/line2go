package line2go

import (
	"line2go/linethrift"
	"line2go/thrift"
)

func (client *IcecreamClient) FetchOperations() (op []*line.Operation, err error) {

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
	client.PollingClient.headerConfig.Do(func() {
		client.px_ls_header = client.PollingClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.PollingClient, client.px_ls_header)
	})
	// end after section
	return
}
