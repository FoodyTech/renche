package main

import (
	"github.com/FoodyTech/renche/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		panic(err)
	}
}
