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
