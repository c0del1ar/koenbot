package main

import (
	"context"
	"koenbot/libtools"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func ResponseMessage(msg string, client *whatsmeow.Client, r *events.Message) error {
	if msg == ".menu" {
		client.SendMessage(context.Background(), r.Info.Chat, &waProto.Message{
			Conversation: proto.String(textMenu),
		})
	}
	if msg == ".menuHack" {
		client.SendMessage(context.Background(), r.Info.Chat, &waProto.Message{
			Conversation: proto.String(textMenuHack),
		})
	}
	if msg == ".sticker" {
		client.SendMessage(context.Background(), r.Info.Chat, &waProto.Message{
			Conversation: proto.String("Please send an image with the caption '.sticker' to convert it into a sticker."),
		})
	}
	if r.Message.ImageMessage != nil {
		caption := strings.TrimSpace(r.Message.GetImageMessage().GetCaption())
		if caption == ".sticker" {
			return libtools.ConvertSticker(client, r)
		}
	}
	return nil
}

func ResponseJoinedGroup(client *whatsmeow.Client, r *events.JoinedGroup) error {
	client.SendMessage(context.Background(), r.JID, &waProto.Message{
		Conversation: proto.String("Hello! Welcome to the group!"),
	})
	return nil
}
