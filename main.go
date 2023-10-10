package main

import (
	"kanban-board/config"
	"kanban-board/routes"

	"github.com/labstack/echo/v4"
)


func main() {
	// load config
	var cfg = config.InitConfig()

	// open mysql connection
	dbMysql := config.InitMyqlConn(cfg)

	// migrate db
	config.Migrate(dbMysql)

	e := echo.New()
	// init routing
	routes.InitRouter(e, dbMysql)
	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8000"))
}
