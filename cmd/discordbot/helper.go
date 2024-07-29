package main

import "fmt"

func generateInviteLink(clientID string) string {
	const linkFormat = "https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot"
	return fmt.Sprintf(linkFormat, clientID)
}
