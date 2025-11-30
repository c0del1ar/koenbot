package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"koenbot/src/libs"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name: "(xnxx|xnxxdl)",
		As:   []string{"xnxx"},
		Description: "Download XNXX video. Usage:\n" +
			"> .xnxx search query | url\n\n" +
			"Example:\n" +
			"> .xnxx crt\nOr\n" +
			"> .xnxx https://www.xnxx.com/video-abc123/example_video",
		Tags:     "downloader",
		IsPrefix: true,
		IsQuerry: true,
		IsWaitt:  true,
		After: func(client *libs.NewClientImpl, m *libs.IMessage) {
			if strings.Contains(m.QuotedMsg.GetStanzaID(), "DLINK") {
				m.Reply(">//< Processing your choice " + m.PushName + "-san")
				pattern := regexp.MustCompile(`[0-9]`)

				if pattern.MatchString(m.Body) {
					idqueryStr := m.Body
					idquery, err := strconv.Atoi(idqueryStr)
					if err != nil {
						m.Reply("Invalid input: please enter a valid number")
						return
					}
					apiURL := "http://10.77.0.23:41000/info?url=" + XnSearchCache["res"][idquery].URL

					xninfo, err := getXnInfo(apiURL)
					if err != nil {
						hLog.Error(fmt.Sprintf("error in getXnInfo():  %v", err))
						m.Reply("Something went wrong from our API")
						return
					}

					dlink := xninfo.DownloadLink
					// remove query params like ?ui=... to get plain mp4 url
					if idx := len(dlink); idx > 0 {
						if qidx := regexp.MustCompile(`\?ui\=.*`).FindStringIndex(dlink); qidx != nil {
							dlink = dlink[:qidx[0]]
						}
					}

					// Download mp4 file dari direct link
					resp, err := http.Get(dlink)
					if err != nil {
						hLog.Error(fmt.Sprintf("error download video:  %v", err))
						m.Reply("Gagal download video")
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						m.Reply(fmt.Sprintf("Gagal download video (status: %s)", resp.Status))
						return
					}

					tmpFile, err := os.CreateTemp("", "xnxx-*.mp4")
					if err != nil {
						hLog.Error(fmt.Sprintf("error membuat file sementara: %v", err))
						m.Reply("Gagal menyiapkan penyimpanan video")
						return
					}
					defer func() {
						tmpFile.Close()
						os.Remove(tmpFile.Name())
					}()

					if _, err := io.Copy(tmpFile, resp.Body); err != nil {
						hLog.Error(fmt.Sprintf("error menyimpan video: %v", err))
						m.Reply("Gagal menyimpan video")
						return
					}

					if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
						hLog.Error(fmt.Sprintf("error reset pointer file: %v", err))
						m.Reply("Gagal membaca ulang video")
						return
					}

					bytes, err := io.ReadAll(tmpFile)
					if err != nil {
						hLog.Error(fmt.Sprintf("error membaca file:  %v", err))
						m.Reply("Gagal membaca video")
						return
					}

					ok, err := client.SendDocument(m.From, bytes, xninfo.Title+".mp4",
						"Here is your video: ~\n\n"+
							"*Title:* "+xninfo.Title,
						nil)
					if err != nil {
						hLog.Error(fmt.Sprintf("error send video:  %v", err))
						m.Reply("Error to sending video.")
						return
					}

					hLog.Info(fmt.Sprintf("video sent: %v", ok))
				}
			}
		},
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			var apiURL string
			query := m.Querry

			if len(query) >= 4 && query[:4] == "http" {
				// execute if query given is link
				var xRegexp = regexp.MustCompile(`^(https?://)?(www\.)?(xnxx\.com)/.+`)
				apiURL = "http://10.77.0.23:41000/info?url=" + query // 51.79.230.125 for external nat

				isXRegexp := func(url string) bool {
					return xRegexp.MatchString(url)
				}

				if !isXRegexp(query) {
					m.Reply("This is not XnxX valid url")
					return
				}

				xninfo, err := getXnInfo(apiURL)
				if err != nil {
					hLog.Error(fmt.Sprintf("error in getXnInfo():  %v", err))
					m.Reply("Something went wrong from our API")
					return
				}

				dlink := xninfo.DownloadLink
				// remove query params like ?ui=... to get plain mp4 url
				if idx := len(dlink); idx > 0 {
					if qidx := regexp.MustCompile(`\?ui\=.*`).FindStringIndex(dlink); qidx != nil {
						dlink = dlink[:qidx[0]]
					}
				}

				// Download mp4 file dari direct link
				resp, err := http.Get(dlink)
				if err != nil {
					hLog.Error(fmt.Sprintf("error download video:  %v", err))
					m.Reply("Gagal download video")
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					m.Reply(fmt.Sprintf("Gagal download video (status: %s)", resp.Status))
					return
				}

				tmpFile, err := os.CreateTemp("", "xnxx-*.mp4")
				if err != nil {
					hLog.Error(fmt.Sprintf("error membuat file sementara: %v", err))
					m.Reply("Gagal menyiapkan penyimpanan video")
					return
				}
				defer func() {
					tmpFile.Close()
					os.Remove(tmpFile.Name())
				}()

				if _, err := io.Copy(tmpFile, resp.Body); err != nil {
					hLog.Error(fmt.Sprintf("error menyimpan video: %v", err))
					m.Reply("Gagal menyimpan video")
					return
				}

				if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
					hLog.Error(fmt.Sprintf("error reset pointer file: %v", err))
					m.Reply("Gagal membaca ulang video")
					return
				}

				bytes, err := io.ReadAll(tmpFile)
				if err != nil {
					hLog.Error(fmt.Sprintf("error membaca file:  %v", err))
					m.Reply("Gagal membaca video")
					return
				}

				ok, err := client.SendDocument(m.From, bytes, xninfo.Title+".mp4",
					"Here is your video: ~\n\n"+
						"*Title:* "+xninfo.Title,
					nil)
				if err != nil {
					hLog.Error(fmt.Sprintf("error send video:  %v", err))
					m.Reply("Error to sending video.")
					return
				}

				hLog.Info(fmt.Sprintf("video sent: %v", ok))

				return
			}

			apiURL = "http://10.77.0.23:41000/search?query=" + query // 51.79.230.125 for external nat
			xnresults, err := getXnSearch(apiURL, m)
			XnSearchCache["res"] = xnresults.Results
			if err != nil {
				hLog.Error(fmt.Sprintf("error in getXnSearch():  %v", err))
				m.Reply("Sorry >//<. There's an error getting search data")
				return
			}

			//var params []string
			var listms string

			// reply send message func
			for i, r := range xnresults.Results {
				listms += fmt.Sprintf("%d. %s\n", i, r.Title)
				//params = append(params, r.URL)
			}

			m.Reply("Please choose video bellows:\n\n"+listms, whatsmeow.SendRequestExtra{
				ID: client.GenerateMessageID("DLINK"),
			})

			// interactive btn message func
			// Get jsonList from getXnSearch result
			/* var buttons []*waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton

					for i, r := range xnresults.Results {
						params := fmt.Sprintf(`{
				"id":"xndl_%d",
				"index":%d,
				"url": %q,
				"title": %q
			}`, i, i, r.URL, r.Title)

						buttons = append(buttons, &waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
							Name:             proto.String("xnxx_select"),
							ButtonParamsJSON: proto.String(params),
						})
					}

					// optional params global
					params := map[string]any{
						"type":  "xnxx_search",
						"count": len(xnresults.Results),
					}

					client.SendInteractiveMessage(
						m.From,
						buttons,
						params,
						"Choose video to download below:",
						nil,
					) */

		},
	})

	//libs.AddButtonHandler("xnxx_select", XnxxButtonHandler)
}

var XnSearchCache = map[string][]XnResult{}

type XnSearch struct {
	Status  string     `json:"status"`
	Query   string     `json:"query"`
	Page    int        `json:"page"`
	Results []XnResult `json:"results"`
}

type XnResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type XnInfo struct {
	DownloadLink string `json:"dlink"`
	Title        string `json:"title"`
	Thumbnail    string `json:"thumbnail"`
}

type APIResponse struct {
	Status string `json:"status"`
	Info   XnInfo `json:"info"`
}

func getXnSearch(apiUrl string, m *libs.IMessage) (XnSearch, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return XnSearch{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return XnSearch{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return XnSearch{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var results XnSearch
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return XnSearch{}, err
	}

	// Simpan ke cache
	cacheKey := m.Message.Chat.String()
	XnSearchCache[cacheKey] = results.Results

	// Buat interactive list
	var rows []map[string]string
	for i, r := range results.Results {
		rows = append(rows, map[string]string{
			"id":    fmt.Sprintf("xndl_%d", i),
			"title": r.Title,
		})
	}

	jsonList, _ := json.Marshal(map[string]interface{}{
		"sections": []map[string]interface{}{
			{
				"title": "XnxX video search results",
				"rows":  rows,
			},
		},
	})
	_ = jsonList

	return results, nil
}

func getXnInfo(apiUrl string) (XnInfo, error) {
	var xninfo XnInfo

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return xninfo, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return xninfo, fmt.Errorf("failed to make HTTP request to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return xninfo, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return xninfo, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if apiResponse.Info.DownloadLink == "" {
		return xninfo, fmt.Errorf("MP4 URL not found in JSON response")
	}

	return apiResponse.Info, nil
}

func XnxxButtonHandler(client *libs.NewClientImpl, info *libs.IMessage, params map[string]any) {

	url := params["url"].(string)
	title := params["title"].(string)

	apiURL := "http://51.79.230.125:41000/info?url=" + url

	// Ambil info download
	xninfo, err := getXnInfo(apiURL)
	if err != nil {
		client.SendText(info.From, "Failed to fetch video info", nil)
		return
	}

	// Download file
	bytes, err := client.GetBytes(xninfo.DownloadLink)
	if err != nil {
		client.SendText(info.From, "Failed to download video", nil)
		return
	}

	// Kirim video
	client.SendVideo(
		info.From,
		bytes,
		"Here is your video:\n"+title,
		nil,
	)
}
