package main

import (
	"pethelp-backend/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	app.NewApp().Run()
}
