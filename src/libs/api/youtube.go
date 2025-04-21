package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"koenbot/src/typings"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

var youtubeRegex = regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/.+`)

func IsYoutubeURL(url string) bool {
	return youtubeRegex.MatchString(url)
}

func parseSeconds(s string) string {
	seconds, err := strconv.ParseFloat(s, 64)
	if nil != err {
		return ""
	}
	return fmt.Sprintf("%v", seconds)
}

func YoutubeDL(uri string) (typings.YoutubeInfos, error) {
	if !IsYoutubeURL(uri) {
		return typings.YoutubeInfos{}, errors.New("url invalid")
	}
	ydl := ytdlp.New().
		PrintJSON().
		NoProgress().
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		NoPlaylist().
		Continue().
		ProgressFunc(100*time.Millisecond, func(prog ytdlp.ProgressUpdate) {
			fmt.Printf( //nolint:forbidigo
				"%s @ %s [eta: %s] :: %s\n",
				prog.Status,
				prog.PercentString(),
				prog.ETA(),
				prog.Filename,
			)
		}).
		Output("%(extractor)s - %(title)s.%(ext)s")

	req, err := ydl.Run(context.TODO(), uri)
	if err != nil {
		return typings.YoutubeInfos{}, err
	}

	f, err := os.Create("results.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	if err = enc.Encode(req); err != nil {
		panic(err)
	}

	log.Println("wrote results to results.json")

	var data map[string]interface{}
	err = json.Unmarshal(*req.OutputLogs[0].JSON, &data)
	if err != nil {
		return typings.YoutubeInfos{}, err
	}

	info := typings.YoutubeInfo{
		Title:    data["title"].(string),
		Duration: data["duration"].(float64),
		Author:   data["uploader"].(string),
	}
	link := typings.YoutubeLinks{}

	formats, ok := data["requested_formats"].([]interface{})
	if !ok {
		return typings.YoutubeInfos{}, errors.New("invalid formats structure")
	}

	for _, dat := range formats {
		a := dat.(map[string]interface{})
		if a["ext"].(string) != "mp3" {
			continue
		}
		link.Audio = append(link.Audio, typings.YoutubeAV{
			Size:    a["filesize"].(string),
			Format:  a["format"].(string),
			Quality: a["quality"].(string),
			Url: func() (string, error) {
				return data["filename"].(string), nil
			},
		})
	}

	for _, dat := range formats {
		a := dat.(map[string]interface{})
		if a["ext"].(string) != "mp4" {
			continue
		}

		var sizefile string
		if fs, ok := a["filesize"]; ok && fs != nil {
			// filesize is usually a float64 (bytes), convert to readable format or string
			switch v := fs.(type) {
			case float64:
				sizefile = fmt.Sprintf("%.0f", v) // raw byte size as string
			default:
				sizefile = fmt.Sprintf("%v", v) // fallback for any other type
			}
		} else {
			sizefile = "Unknown"
		}

		link.Video = append(link.Video, typings.YoutubeAV{
			Size:   sizefile,
			Format: a["format"].(string),
			Quality: func() string {
				q, ok := a["quality"]
				if !ok || q == nil {
					return "Unknown"
				}
				switch v := q.(type) {
				case string:
					return v
				case float64:
					return fmt.Sprintf("%.0f", v)
				default:
					return fmt.Sprintf("%v", v)
				}
			}(),
			Url: func() (string, error) {
				return data["filename"].(string), nil
			},
		})
	}

	log.Println("success download")
	log.Println("title: ", info.Title)

	return typings.YoutubeInfos{
		Info: info,
		Link: link,
	}, nil
}

func Download(url string, formatID string) (string, error) {
	// Generate a temp filename based on formatID
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("yt_%s.tmp", formatID))

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("download failed with status: " + resp.Status)
	}

	// Write to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	// Download success
	return url, nil
}
