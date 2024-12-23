package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/greeneg/todoer/controllers"
)

func PublicRoutes(g *gin.RouterGroup, i *controllers.TodoerService) {
	g.GET("/health", i.GetHealth) // service health
}

func PrivateRoutes(g *gin.RouterGroup, i *controllers.TodoerService) {
	// todo related routes
	g.GET("/todo", i.GetTodos)          // get todos
	g.GET("/todo/:id", i.GetTodoById)	// get todo by its Id
	g.POST("/todo", i.CreateTodo)       // create a new todo
	g.DELETE("/todo/:id", i.DeleteTodo) // trash a todo entry
	g.PUT("/todo/:id/:status", i.UpdateTodo)    // replace todo status
	// user related routes
	g.GET("/user/id/:id", i.GetUserById)            // get user by id
	g.GET("/user/name/:name", i.GetUserByUserName)  // get user by username
	g.GET("/user/:name/status", i.GetUserStatus)    // get whether a user is locked or not
	g.GET("/users", i.GetUsers)                     // get users
	g.POST("/user", i.CreateUser)                   // create new user
	g.PATCH("/user/:name", i.ChangeAccountPassword) // update a user password
	g.PATCH("/user/:name/status", i.SetUserStatus)  // lock a user
	g.DELETE("/user/:name", i.DeleteUser)           // trash a user
}
