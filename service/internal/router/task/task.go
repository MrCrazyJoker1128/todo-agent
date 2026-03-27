package task

import (
	"service/internal/controller/task"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRoutes(group *ghttp.RouterGroup) {
	taskController := task.NewTaskController()
	group.Group("/task", func(group *ghttp.RouterGroup) {
		group.POST("/create", taskController.CreateTask)
		group.DELETE("/delete", taskController.DeleteTask)
		group.PUT("/update", taskController.UpdateTask)
		group.GET("/list", taskController.GetTaskList)
	})
}
