package main

import (
	"html/template"
	"github.com/gin-gonic/gin"

	"project/workflow/singals"
	"project/workflow/models"
	_ "project/workflow/models"
	_ "project/workflow/middleware"
	. "project/workflow/middleware"
	"project/workflow/utils"
	"project/workflow/views"
)

func main() {
	//  信号监听
	go singals.ListenSingal()

	// Creates a router without any middleware by default
	router := gin.New()

	router.Use(gin.Logger())

	router.SetFuncMap(template.FuncMap{
		"getUserById": models.GetUserById,
		"getTypeById": models.GetTypeById,
		"getStateById": models.GetStateById,
		"timeFormat": utils.TimeFormat,
		"strToInt": utils.StrToInt,
		"getTranUsersById": models.GetTranUsersById,
		"perm": models.Perm,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/login", views.LoginApi)
	router.GET("/logout", views.LogoutApi)

	router.Use(Authenticate())
	router.GET("/index", views.Index)
	router.POST("/cms/create", views.CreateTask)
	router.GET("/cms/:objId/detail", views.TaskDetail)
	router.POST("/cms/tran", views.TranAction)

	router.Run(":8200")
}