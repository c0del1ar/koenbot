package helpers

import "os"

var Public = false

func BotName() string {
	return os.Getenv("Name_Bot")
}

func Owner() string {
	return os.Getenv("Owner_Number")
}
