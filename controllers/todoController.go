package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/greeneg/todoer/model"
	"github.com/gin-gonic/gin"
)

// CreateTodo Register a new todo
//
//	@Summary		Register todo
//	@Description	Add a new todo
//	@Tags			todo
//	@Accept			json
//	@Produce		json
//	@Param			todo	body	model.ProposedTodo	true	"Todo Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/todo [post]
func (t *TodoerService) CreateTodo(c *gin.Context) {
	_, authed := t.GetUserId(c)
	if authed {
		var json model.ProposedTodo
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateTodo(json)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Todo has been created"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteTodo Remove a todo
//
//	@Summary		Delete todo
//	@Description	Delete a todo
//	@Tags			todo
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Todo Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/todo/{id} [delete]
func (t *TodoerService) DeleteTodo(c *gin.Context) {
	_, authed := t.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("id"))
		status, err := model.DeleteTodo(id)
		if err != nil {
			log.Println("ERROR: Cannot delete todo: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove todo! " + string(err.Error())})
			return
		}

		if status {
			idString := strconv.Itoa(id)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User " + idString + " has been removed"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetTodos Retrieve list of all todos
//
//	@Summary		Retrieve list of todos
//	@Description	Retrieve list of all todos
//	@Tags			todo
//	@Produce		json
//	@Success		200	{object}	model.TodoList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/todo [get]
func (t *TodoerService) GetTodos(c *gin.Context) {
	_, authed := t.GetUserId(c)
	if authed {
		todos, err := model.GetTodos()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		if todos == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": todos})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetTodoById Retrieve a todo by its Id
//
//	@Summary		Retrieve a todo by its Id
//	@Description	Retrieve a todo by its Id
//	@Tags			todo
//	@Produce		json
//	@Param			id	path int true "Todo ID"
//	@Success		200	{object}	model.Todo
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/todo/{id} [get]
func (t *TodoerService) GetTodoById(c *gin.Context) {
	_, authed := t.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("id"))
		ent, err := model.GetTodoById(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		if ent.Description == "" {
			strId := strconv.Itoa(id)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with todo id " + strId})
		} else {
			c.IndentedJSON(http.StatusOK, ent)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// UpdateTodo	Update the status of a todo
//
//	@Summary	Update the status of a todo
//	@Description	Updates the status field of a todo
//	@Tags		todo
//	@Accept		json
//	@Produce	json
//	@Param		id	path int true "Todo ID"
//	@Param		status	path string true "Todo Status"
//	@Success	200	{object}	model.SuccessMsg
//	@Failure	400	{object}	model.FailureMsg
//	@Router		/todo/{id}/{status} [put]
func (t *TodoerService) UpdateTodo(c *gin.Context) {
	_, authed := t.GetUserId(c)
	if authed {
		// first, _get_ the Todo, then update it with the data
		id, _ := strconv.Atoi(c.Param("id"))
		ent, err := model.GetTodoById(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}
		// now validate that the status string is one we know
		status := c.Param("status")
		statusId, err := model.GetStatusByName(status)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}
		// now update ent with the new status
		ent, err = model.UpdateTodo(ent.Id, statusId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		} else {
			c.IndentedJSON(http.StatusOK, ent)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
