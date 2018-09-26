package routes

import (
	"github.com/mohammedajao/rest-api/app"
	"github.com/mohammedajao/rest-api/app/controllers"
)

func Init(a *app.App) {
	UserController := &controllers.UserController{a.DB}
	UserController.Init(a)
}
