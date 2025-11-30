package libs

import (
	"context"
	"encoding/json"
	"koenbot/src/helpers"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type IHandler struct {
	Container *store.Device
}

func NewHandler(container *sqlstore.Container) *IHandler {
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}
	return &IHandler{
		Container: deviceStore,
	}
}

func (h *IHandler) Client(jbot ...bool) *whatsmeow.Client {
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(h.Container, clientLog)
	client.AddEventHandler(h.RegisterHandler(client, jbot...))
	return client
}

func (h *IHandler) RegisterHandler(client *whatsmeow.Client, jbot ...bool) func(evt interface{}) {
	return func(evt interface{}) {
		sock := NewClient(client)
		switch v := evt.(type) {
		case *events.Message:
			if ir := v.Message.GetInteractiveResponseMessage(); ir != nil {
				nf := ir.GetNativeFlowResponseMessage()
				if nf != nil {
					var data map[string]any
					json.Unmarshal([]byte(nf.GetParamsJSON()), &data)

					name := nf.GetName() // tombol bernama apa

					if handler, ok := ButtonHandlers[name]; ok {
						msg := NewSmsg(v, sock, jbot...)
						go handler(sock, msg, data)
						return
					}
				}
			}
			// ==== END OF BUTTON HANDLER ====

			// ==== GREETING HANDLER ====
			if v.Message.GetConversation() != "" {
				text := strings.ToLower(v.Message.GetConversation())
				from := v.Info.Chat

				switch text {
				case "p", "pp", "p!", "p?":
					sock.SendText(from, "P juga kak 😸", nil)
					return

				case "hi", "hai", "hii", "halo", "hello", "helo":
					sock.SendText(from, "Hay uwu>\\\\<! Can i help u? 😊", nil)
					return
				}
			}
			// ==== END OF GREETING HANDLER ====

			m := NewSmsg(v, sock, jbot...)
			if !helpers.Public && !m.IsOwner {
				return
			}
			// Read message
			sock.WA.MarkRead(context.Background(), []string{m.StanzaId}, time.Now(), m.From, m.Sender)
			// Get command
			go Get(sock, m)
			return
		case *events.StreamError:
			helpers.ErrorLogger.Printf("stream error: code=%s raw=%v", v.Code, v.Raw)
		case *events.ConnectFailure:
			helpers.ErrorLogger.Printf("connect failure: %s (message=%s)", v.Reason.String(), v.Message)
		case *events.LoggedOut:
			helpers.ErrorLogger.Printf("logged out: on_connect=%v reason=%s", v.OnConnect, v.Reason.String())
		case *events.TemporaryBan:
			helpers.ErrorLogger.Printf("temporary ban: %s", v.String())
		case events.PermanentDisconnect:
			helpers.ErrorLogger.Printf("permanent disconnect: %s", v.PermanentDisconnectDescription())
		case *events.Connected, *events.PushNameSetting:
			if len(client.Store.PushName) == 0 {
				return
			}
			client.SendPresence(context.Background(), types.PresenceAvailable)
		}
	}
}

var ButtonHandlers = map[string]func(*NewClientImpl, *IMessage, map[string]any){}

func AddButtonHandler(name string, handler func(*NewClientImpl, *IMessage, map[string]any)) {
	ButtonHandlers[name] = handler
}
