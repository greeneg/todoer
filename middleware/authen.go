package middleware

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/greeneg/todoer/globals"
	"github.com/greeneg/todoer/helpers"
	"github.com/greeneg/todoer/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func processAuthorizationHeader(authHeader string) (string, string) {
	// split the header value at the space
	encodedString := strings.Split(authHeader, " ")

	// remove base64 encoding
	decodedString, _ := base64.StdEncoding.DecodeString(encodedString[1])

	// now lets return both the
	authValues := strings.Split(string(decodedString), ":")

	return authValues[0], authValues[1]
}

func AuthCheck(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		log.Println("INFO: No session found. Attempting to check for authentication headers")
		baHeader := c.GetHeader("Authorization")
		if baHeader == "" {
			log.Println("ERROR: No authentication header found. Aborting")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
			c.Abort()
			return
		}
		// otherwise, lets process that header
		username, password := processAuthorizationHeader(baHeader)
		authStatus := helpers.CheckUserPass(username, password)
		if authStatus {
			session.Set(globals.UserKey, username)
			if err := session.Save(); err != nil {
				c.IndentedJSON(http.StatusInternalServerError,
					gin.H{"error": "failed to save user session"})
				// session saving is not fatal, so allow them to proceed
			}
			log.Println("INFO: Authenticated")
		} else {
			log.Println("ERROR: Authentication failed. Aborting")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
			c.Abort()
			return
		}
	} else {
		userString := fmt.Sprintf("%v", user)
		log.Println("INFO: Session found: User: " + userString)
		log.Println("INFO: Checking if user is locked or not...")
		user, err := model.GetUserByUserName(userString)
		if err != nil {
			log.Println("ERROR: " + string(err.Error()))
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unable to authenticate: " + err.Error()})
			c.Abort()
			return
		}
		status := helpers.CheckIsNotLocked(user)
		if status {
			log.Println("INFO: Authenticated")
		} else {
			log.Println("WARN: User '" + userString + "' is locked!")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
			c.Abort()
			return
		}
	}
	c.Next()
}
