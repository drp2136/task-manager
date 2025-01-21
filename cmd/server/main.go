package main

import (
	"log"
	"net/http"
	"task-manager/internal/model"
	"task-manager/internal/router"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {
	// db, err := gorm.Open("postgres", "postgres:password@/task_manager?charset=utf8&parseTime=True&loc=Local&sslmode=disable")
	db, err := gorm.Open("postgres", "user=postgres dbname=task_manager sslmode=disable host=localhost port=5432 password=password")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&model.Task{})

	// taskService := &service.TaskService{DB: db}
	taskService := service.NewTaskService(db)

	r := gin.Default()

	router.SetupRouter(r, taskService)

	log.Fatal(http.ListenAndServe(":8080", r))
}
