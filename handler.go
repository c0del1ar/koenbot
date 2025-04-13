package main

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			err := ResponseMessage(messageBody, client, v)
			if err != nil {
				fmt.Println("Error processing message:", err)
			}
		case *events.JoinedGroup:
			err := ResponseJoinedGroup(client, v)
			if err != nil {
				fmt.Println("Error processing joined group:", err)
			}
		case *events.GroupInfo:
			err := ResponseGroupUpdate(client, v)
			if err != nil {
				fmt.Println("Error processing group update")
			}
		}
	}
}
