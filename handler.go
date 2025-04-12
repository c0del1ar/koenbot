package main

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			err := ResponseMessage(messageBody, client, v)
			if err != nil {
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String(err.Error()),
				})
			}
		case *events.JoinedGroup:
			if ResponseJoinedGroup(client, v) != nil {
				client.SendMessage(context.Background(), v.JID, &waProto.Message{
					Conversation: proto.String("Error processing your request."),
				})
			}
		}
	}
}
