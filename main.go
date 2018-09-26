package main

import (
	"github.com/mohammedajao/rest-api/app"
	"github.com/mohammedajao/rest-api/app/routes"
)

func main() {
	a := app.App{}
	a.Initialize()
	routes.Init(&a)
	a.Run("8000")
}
