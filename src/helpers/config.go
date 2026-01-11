package helpers

import (
	"context"
	"os"
)

var Public = false

func BotName() string {
	return os.Getenv("Name_Bot")
}

func Owner() string {
	return os.Getenv("Owner_Number")
}

func CtxB() context.Context {
	return context.Background()
}
