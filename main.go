package main

import (
	"fmt"
	"github.com/Nepooomuk/ToDoList/model"
	"github.com/Nepooomuk/ToDoList/redisclient"
	"github.com/garyburd/redigo/redis"
	"github.com/kataras/iris"
	"log"
)

var pool = redisclient.NewPool()

func main() {

	c := pool.Get()
	defer c.Close()

	redis.Strings(c.Do("SET", "message", "Redis is up and running."))
	msg, err := redis.String(c.Do("GET", "message"))
	if err != nil {
		log.Print(err)
	}
	fmt.Println(msg)

	app := iris.New()

	app.Get("/task/:id", taskGetHandler)
	app.Get("/task", taskGetHandler)
	app.Post("/task/:id", taskPostHandler)
	app.Delete("/task/:id", taskDeleteHandler)

	app.Listen(":8080")
}

func taskPostHandler(ctx *iris.Context) {
	c := pool.Get()
	defer c.Close()

	task := &model.Task{}
	if err := ctx.ReadJSON(&task); err != nil {
		log.Print(err)
		ctx.JSON(500, err.Error())
	} else {
		msg, err := redis.Strings(c.Do("SET", task.ID, task.Name))
		if err != nil {
			log.Print(err)
		}
		log.Print(msg)
	}
	ctx.JSON(iris.StatusOK, task)
}

func taskGetHandler(ctx *iris.Context) {
	c := pool.Get()
	defer c.Close()

	keys, err := redis.Strings(c.Do("KEYS", "*"))
	if err != nil {
		log.Print(err)
	}

	tasksRepo := make([]string, 0)
	for _, values := range keys {
		value, err := redis.String(c.Do("GET", values))
		if err != nil {
			log.Print(err)
		}
		tasksRepo = append(tasksRepo, value)
		log.Print(value)
	}

	ctx.JSON(iris.StatusOK, map[string]interface{}{
		"allTasks": &tasksRepo,
	})
}

func taskDeleteHandler(ctx *iris.Context) {
	taskID, err := ctx.ParamInt("id")
	if err != nil {
		log.Print(err)
	}

	c := pool.Get()
	defer c.Close()

	c.Do("DEL", taskID)
	ctx.JSON(iris.StatusOK, taskID)
}
