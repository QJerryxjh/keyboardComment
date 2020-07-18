package router

import (
	"keyboardComment/controller/managerController"
)

func ManagerRouter(baseRoutePath string) {
	r := Router.Group("/" + baseRoutePath)
	{
		r.GET("/list", managerController.GetManagers)
	}
}
