package main

import (
	"bank/database"
	"bank/routes"
)

func main() {

	pg_db := database.Connect()
	defer pg_db.Close()

	schemaCreationErr := database.CreateSchema(pg_db)

	if schemaCreationErr != nil {
		panic(schemaCreationErr)
	}

	routes.CreateRoutes()

}
