package main

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func ResponseMessage(msg string, client *whatsmeow.Client, r *events.Message) error {
	if msg == ".menu" {
		client.SendMessage(context.Background(), r.Info.Chat, &waProto.Message{
			Conversation: proto.String(textMenu),
		})
	}
	return nil
}

func ResponseJoinedGroup(client *whatsmeow.Client, r *events.JoinedGroup) error {
	client.SendMessage(context.Background(), r.JID, &waProto.Message{
		Conversation: proto.String("Hello! Welcome to the group!"),
	})
	return nil
}
