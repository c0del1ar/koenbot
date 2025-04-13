package main

import (
	"context"
	"fmt"
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
		Conversation: proto.String("Hi all\n\n> I'm KoenChan, a bot created by AryaKun.\n> I'm here to assist you with various tools.\n\n_- KoenChan -_"),
	})
	return nil
}

func ResponseGroupUpdate(client *whatsmeow.Client, r *events.GroupInfo) error {
	ctx := context.Background()
	if join := r.Join; join != nil {
		user := join[0].User
		message := fmt.Sprintf("Hey, @~%s! Welcome to the group ma bro!\n\n_- KoenChan -_", user)
		_, err := client.SendMessage(ctx, r.JID, &waProto.Message{
			Conversation: proto.String(message),
		})
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}
	if leave := r.Leave; leave != nil {
		user := leave[0].User
		message := fmt.Sprintf("Ah, see u later @~%s!\n\n_- KoenChan -_", user)
		_, err := client.SendMessage(ctx, r.JID, &waProto.Message{
			Conversation: proto.String(message),
		})
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}
	return nil
}
