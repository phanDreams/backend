package main

import (
	"pethelp-backend/internal/app"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.NewApp()).Run()
}
