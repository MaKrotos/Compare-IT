package main

import (
	"backend/internal/app"
)

func main() {
	application := app.New()
	application.Run()
}