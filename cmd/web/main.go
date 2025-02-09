package main

import (
	"golang-rest-api/infrastructure"
	"golang-rest-api/internal/api"
)

func main() {
	// DB Connection
	db := infrastructure.NewDbConnection()
	defer db.Close()

	// HTTP
	routes := api.InitRoutes(db)
	routes.Run(":8080")
}
