package commands

import (
	"koenbot/src/libs"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     "(sc|source)",
		As:       []string{"sc"},
		Tags:     "main",
		IsPrefix: true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			m.Reply("https://github.com/c0del1ar/koenbot\n\n_Origin from https://github.com/fckvania/MaoGo_\n\n_Free Not For Sell_")
		},
	})
}
