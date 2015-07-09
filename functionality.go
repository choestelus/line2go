package main

import "line"

func (client *IcecreamClient) GetLastOpRevision() (r int64, err error) {
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}
	r, err = client.CommandClient.GetLastOpRevision()
	if err != nil {
		return
	}
	client.setCommandState()
	client.opRevision = r
	return
}

func (client *IcecreamClient) GetProfile() (profile *line.Profile, err error) {
	if client.commandClientState == true {
		SetHeaderForClientReuse(client.CommandClient, client.cx_ls_header)
	} else {
		SetHeaderForClientInit(client.CommandClient, client.authToken, client.userAgent, client.x_line_application)
	}

	profile, err = client.CommandClient.GetProfile()
	client.setCommandState()

	return
}
