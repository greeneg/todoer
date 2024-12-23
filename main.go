package main

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/greeneg/todoer/controllers"
	"github.com/greeneg/todoer/globals"
	"github.com/greeneg/todoer/helpers"
	"github.com/greeneg/todoer/middleware"
	"github.com/greeneg/todoer/model"
	"github.com/greeneg/todoer/routes"
)

//	@title		Todoer
//	@version	0.0.1
//	@description	An API for managing todos

//	@contact.name	Gary Greene
//	@contact.url	https://github.com/greeneg/todoer

//	@securityDefinitions.basic	BasicAuth

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/api/v1

// @schemas	http https
func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	helpers.FatalCheckError(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := globals.Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	helpers.FatalCheckError(err)
	err = json.Unmarshal(jsonContent, &config)
	helpers.FatalCheckError(err)

	// create an app object that contains our routes and the configuration
	TodoerService := new(controllers.TodoerService)
	TodoerService.AppPath = appdir
	TodoerService.ConfigPath = configDir
	TodoerService.ConfStruct = config

	err = model.ConnectDatabase(TodoerService.ConfStruct.DbPath)
	helpers.FatalCheckError(err)

	// This will ensure that the angular files are served correctly
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	// some defaults for using session support
	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	// API
	public := r.Group("/api/v1")
	routes.PublicRoutes(public, TodoerService)

	private := r.Group("/api/v1")
	private.Use(middleware.AuthCheck)
	routes.PrivateRoutes(private, TodoerService)

	// Front-end stuff

	// swagger doc
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tcpPort := strconv.Itoa(TodoerService.ConfStruct.TcpPort)
	tlsTcpPort := strconv.Itoa(TodoerService.ConfStruct.TLSTcpPort)
	tlsPemFile := TodoerService.ConfStruct.TLSPemFile
	tlsKeyFile := TodoerService.ConfStruct.TLSKeyFile
	if TodoerService.ConfStruct.UseTLS {
		r.RunTLS(":"+tlsTcpPort, tlsPemFile, tlsKeyFile)
	} else {
		r.Run(":" + tcpPort)
	}
}
