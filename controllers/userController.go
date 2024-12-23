package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/greeneg/todoer/model"
	"github.com/gin-gonic/gin"
)

// CreateUser Register a user for authentication and authorization
//
//	@Summary		Register user
//	@Description	Add a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.ProposedUser	true	"User Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user [post]
func (g *TodoerService) CreateUser(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		var json model.ProposedUser
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateUser(json)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// ChangeAccountPassowrd Change an account's password
//
//	@Summary		Change password
//	@Description	Change password
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Param			changePassword	body	model.PasswordChange	true	"Password data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name} [patch]
func (g *TodoerService) ChangeAccountPassword(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		username := c.Param("name")
		var json model.PasswordChange
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.ChangeAccountPassword(username, json.OldPassword, json.NewPassword)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User '" + username + "' has changed their password"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User password could not be updated!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteUser Remove a user for authentication and authorization
//
//	@Summary		Delete user
//	@Description	Delete a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name} [delete]
func (g *TodoerService) DeleteUser(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		username := c.Param("name")
		status, err := model.DeleteUser(username)
		if err != nil {
			log.Println("ERROR: Cannot delete user: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User " + username + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUserStatus Retrieve the active status of a user. Can be either 'enabled' or 'locked'
//
//	@Summary		Retrieve a user's active status. Can be either 'enabled' or 'locked'
//	@Description	Retrieve a user's active status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.UserStatusMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name}/status [get]
func (g *TodoerService) GetUserStatus(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		username := c.Param("name")
		status, err := model.GetUserStatus(username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to get the user " + username + " status: " + string(err.Error())})
			return
		}

		if status != "" {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User status: " + status, "userStatus": status})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to retrieve user status"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// SetUserStatus Set the active status of a user. Can be either 'enabled' or 'locked'
//
//	@Summary		Set a user's active status. Can be either 'enabled' or 'locked'
//	@Description	Set a user's active status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.User.UserName	true	"User Data"
//	@Param			name	path	string	true "User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.UserStatusMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name}/status [patch]
func (g *TodoerService) SetUserStatus(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		username := c.Param("name")
		var json model.UserStatus
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.SetUserStatus(username, json)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User '" + username + "' has been " + json.Status})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUsers Retrieve list of all users
//
//	@Summary		Retrieve list of all users
//	@Description	Retrieve list of all users
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	model.UsersList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/users [get]
func (g *TodoerService) GetUsers(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		users, err := model.GetUsers()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		safeUsers := make([]SafeUser, 0)
		for _, user := range users {
			safeUser := SafeUser{}
			safeUser.Id = user.Id
			safeUser.UserName = user.UserName
			safeUser.CreationDate = user.CreationDate

			safeUsers = append(safeUsers, safeUser)
		}

		if users == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": safeUsers})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUserById Retrieve a user by their Id
//
//	@Summary		Retrieve a user by their Id
//	@Description	Retrieve a user by their Id
//	@Tags			user
//	@Produce		json
//	@Param			id	path int true "User ID"
//	@Success		200	{object}	SafeUser
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/id/{id} [get]
func (g *TodoerService) GetUserById(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("id"))
		ent, err := model.GetUserById(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		// don't return the password hash
		safeUser := new(SafeUser)
		safeUser.Id = ent.Id
		safeUser.UserName = ent.UserName
		safeUser.CreationDate = ent.CreationDate

		if ent.UserName == "" {
			strId := strconv.Itoa(id)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with user id " + strId})
		} else {
			c.IndentedJSON(http.StatusOK, safeUser)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUserByName Retrieve a user by their UserName
//
//	@Summary		Retrieve a user by their UserName
//	@Description	Retrieve a user by their UserName
//	@Tags			user
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Success		200	{object}	SafeUser
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/name/{name} [get]
func (g *TodoerService) GetUserByUserName(c *gin.Context) {
	_, authed := g.GetUserId(c)
	if authed {
		username := c.Param("name")
		ent, err := model.GetUserByUserName(username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		// don't return the password hash
		safeUser := new(SafeUser)
		safeUser.Id = ent.Id
		safeUser.UserName = ent.UserName
		safeUser.CreationDate = ent.CreationDate

		if ent.UserName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with user name " + username})
		} else {
			c.IndentedJSON(http.StatusOK, safeUser)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
