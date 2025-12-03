package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"koenbot/src/libs"
	"net/http"
	"strings"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     "(fb|ig|twitter|pinterest|likee|threads|terabox)",
		As:       []string{"fb", "ig", "twitter", "pinterest", "likee", "threads", "terabox"},
		Tags:     "downloader",
		IsPrefix: true,
		IsQuerry: true,
		IsWaitt:  true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			apiUrl := "https://697d7432c3abec333cee5ab7a95e8774.aryakun.id/api/download/"
			platformAliases := map[string]string{
				"fb":      "fb",
				"ig":      "ig",
				"twitter": "twt",
				"likee":   "lk",
				"threads": "th",
				"terabox": "tb",
			}

			// To convert ".fb" to "fb", we can remove the leading dot (if any)
			command := strings.TrimPrefix(strings.ToLower(m.Command), ".")
			platformKey := platformAliases[command]

			apiUrl += platformKey

			jsonData, err := json.Marshal(map[string]string{"url": m.Querry})
			if err != nil {
				m.Reply(err.Error())
				return
			}
			response, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				m.Reply("Sorry, something went wrong from our API")
				hLog.Error(fmt.Sprintf("error in fetching api:  %v", err))
				return
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				m.Reply(err.Error())
				hLog.Error(fmt.Sprintf("error in reading body:  %v", err))
				return
			}
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				m.Reply(err.Error())
				hLog.Error(fmt.Sprintf("error in unmarshalling body:  %v", err))
				return
			}
			if data["status"] == false {
				m.Reply(data["msg"].(string))
				hLog.Error(fmt.Sprintf("error in status:  %v", data["msg"]))
				return
			}

			dataRespon := data["data"].(map[string]interface{})
			var videoUrl string

			if platformKey == "fb" {
				for _, quality := range []string{"high", "medium", "low"} {
					if val, ok := dataRespon[quality].(string); ok && val != "" {
						videoUrl = val
						break
					}
				}
			} else if platformKey == "ig" {
				videoArr, ok := dataRespon["video"].([]interface{})
				if ok && len(videoArr) > 0 {
					if url, valid := videoArr[0].(string); valid && url != "" {
						videoUrl = url
					}
				}
			} else if platformKey == "yt" {
				videoUrl = dataRespon["video_hd"].(string)
			} else if platformKey == "twt" {
				videoUrl = dataRespon["HD"].(string)
			} else {
				videoUrl = dataRespon["video"].(string)
			}

			bytes, err := client.GetBytes(videoUrl)
			if err != nil {
				m.Reply("Sorry senpai, I'm not able to fetch the video")
				hLog.Error(fmt.Sprintf("error in fetching video:  %v", err))
				return
			}
			ok, err := client.SendVideo(m.From, bytes, "Here is your video "+m.PushName+"-san", m.ID)
			if err != nil {
				hLog.Error(fmt.Sprintf("error in sending document:  %v", err))
				return
			}
			hLog.Info(fmt.Sprintf("document sent: %v", ok))
		},
	})
}
