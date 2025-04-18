package commands

import (
	"fmt"
	"koenbot/src/libs"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     "(join)",
		As:       []string{"join"},
		Tags:     "group",
		IsPrefix: true,
		IsOwner:  true,
		IsQuerry: true,
		IsWaitt:  true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			gid, err := client.WA.JoinGroupWithLink(m.Querry)
			if err != nil {
				m.Reply("koenbot tidak bisa gabung ke group itu.")
			} else {
				resp, _ := client.WA.GetGroupInfo(gid)
				m.Reply(fmt.Sprintf("koenbot berhasil gabung ke group %s", resp.Name))
			}
		},
	})
}
