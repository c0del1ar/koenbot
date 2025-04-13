package commands

import (
	"koenbot/src/helpers"
	"koenbot/src/libs"
	"strconv"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     "(setmode|mode)",
		As:       []string{"setmode"},
		Tags:     "owner",
		IsPrefix: true,
		IsOwner:  true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			if m.Querry == "public" {
				helpers.Public = true
				m.Reply("Public Mode: " + strconv.FormatBool(helpers.Public))
			} else {
				helpers.Public = false
				m.Reply("Public Mode: " + strconv.FormatBool(helpers.Public))
			}
		},
	})
}
