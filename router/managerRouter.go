package router

import (
	"keyboardComment/controller/managerController"
)

func ManagerRouter(baseRoutePath string) {
	r := Router.Group("/" + baseRoutePath)
	{
		r.GET("/list", managerController.GetManagers)
		r.POST("/info", managerController.GetManagerInfo)
		r.POST("/delete", managerController.DeleteManager)
		r.POST("/create", managerController.CreateManager)
	}
}
