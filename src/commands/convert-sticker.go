package commands

import (
	"context"
	"koenbot/src/libs"
	"koenbot/src/libs/api"
	"koenbot/src/typings"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     "(s|sticker)",
		As:       []string{"sticker"},
		Tags:     "convert",
		IsPrefix: true,
		IsMedia:  true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			data, _ := client.WA.Download(context.Background(), m.Media)

			s := api.StickerApi(&typings.Sticker{
				File: data,
				Tipe: func() typings.MediaType {
					if m.IsImage || m.IsQuotedImage || m.IsQuotedSticker {
						return typings.IMAGE
					} else {
						return typings.VIDEO
					}
				}(),
			}, &typings.MetadataSticker{
				Author:    "Sally-chan",
				Pack:      "Sally-chan Technolojia",
				KeepScale: true,
				Removebg:  "false",
				Circle: func() bool {
					if m.Querry == "-c" {
						return true
					} else {
						return false
					}
				}(),
			})

			client.SendSticker(m.From, s.Build(), m.ID)

		},
	})
}
