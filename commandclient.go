package line2go

import (
	"line2go/linethrift"
	"line2go/thrift"
)

func (client *IcecreamClient) GetLastOpRevision() (r int64, err error) {
	r, err = client.CommandClient.GetLastOpRevision()
	if err != nil {
		return
	}
	client.opRevision = r

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetProfile() (profile *line.Profile, err error) {

	profile, err = client.CommandClient.GetProfile()

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})
	return
}

func (client *IcecreamClient) GetAllContactIDs() (list []string, err error) {
	list, err = client.CommandClient.GetAllContactIds()

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetContacts(contactIDs []string) (r []*line.Contact, err error) {

	r, err = client.CommandClient.GetContacts(contactIDs)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetContact(contactID string) (r *line.Contact, err error) {

	r, err = client.CommandClient.GetContact(contactID)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetGroupIdsJoined() (r []string, err error) {

	r, err = client.CommandClient.GetGroupIdsJoined()

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetGroupIdsInvited() (r []string, err error) {

	r, err = client.CommandClient.GetGroupIdsInvited()

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetGroups(groupsID []string) (groups []*line.Group, err error) {

	groups, err = client.CommandClient.GetGroups(groupsID)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) GetGroup(groupID string) (group *line.Group, err error) {

	group, err = client.CommandClient.GetGroup(groupID)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

// Default value is start = 1, messageBoxCount = 50
func (client *IcecreamClient) GetMessageBoxCompactWrapUpList(start int32, messageBoxCount int32) (r *line.MessageBoxWrapUpList, err error) {

	r, err = client.CommandClient.GetMessageBoxCompactWrapUpList(start, messageBoxCount)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}

func (client *IcecreamClient) SendTextMessage(id string, text string) (r *line.Message, err error) {

	msg := &line.Message{
		To:          id,
		Text:        text,
		ContentType: line.ContentType_NONE,
	}
	r, err = client.CommandClient.SendMessage(0, msg)

	client.CommandClient.headerConfig.Do(func() {
		client.cx_ls_header = client.CommandClient.Transport.(*thrift.THttpClient).GetResponse().Header.Get("X-LS")
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	})

	return
}
