package main

import "line"

func (client *IcecreamClient) GetLastOpRevision() (r int64, err error) {
	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetLastOpRevision()
	if err != nil {
		return
	}
	client.opRevision = r

	// begin after section
	client.setCommandState()
	// end after section

	return
}

func (client *IcecreamClient) GetProfile() (profile *line.Profile, err error) {
	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	profile, err = client.CommandClient.GetProfile()

	// begin after section
	client.setCommandState()
	// end after section

	return
}

func (client *IcecreamClient) GetAllContactIDs() (list []string, err error) {
	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section
	list, err = client.CommandClient.GetAllContactIds()
	// begin after section
	client.setCommandState()
	// end after section
	return
}

func (client *IcecreamClient) GetContacts(contactIDs []string) (r []*line.Contact, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetContacts(contactIDs)

	// begin after section
	client.setCommandState()
	// end after section
	return
}

func (client *IcecreamClient) GetContact(contactID string) (r *line.Contact, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetContact(contactID)

	// begin after section
	client.setCommandState()
	// end after section
	return
}

func (client *IcecreamClient) GetGroupIdsJoined() (r []string, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetGroupIdsJoined()

	// begin after section
	client.setCommandState()
	// end after section

	return
}

func (client *IcecreamClient) GetGroupIdsInvited() (r []string, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetGroupIdsInvited()

	// begin after section
	client.setCommandState()
	// end after section

	return
}

func (client *IcecreamClient) GetGroups(groupsID []string) (groups []*line.Group, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section
	groups, err = client.CommandClient.GetGroups(groupsID)
	// begin after section
	client.setCommandState()
	// end after section
	return
}

func (client *IcecreamClient) GetGroup(groupID string) (group *line.Group, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section
	group, err = client.CommandClient.GetGroup(groupID)
	// begin after section
	client.setCommandState()
	// end after section
	return
}

// Default value is start = 1, messageBoxCount = 50
func (client *IcecreamClient) GetMessageBoxCompactWrapUpList(start int32, messageBoxCount int32) (r *line.MessageBoxWrapUpList, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	r, err = client.CommandClient.GetMessageBoxCompactWrapUpList(start, messageBoxCount)
	// begin after section
	client.setCommandState()
	// end after section
	return
}

func (client *IcecreamClient) SendTextMessage(id string, text string) (r *line.Message, err error) {

	// begin before section
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	// end before section

	msg := &line.Message{
		To:          id,
		Text:        text,
		ContentType: line.ContentType_NONE,
	}
	r, err = client.CommandClient.SendMessage(0, msg)

	// begin after section
	client.setCommandState()
	// end after section
	return
}
