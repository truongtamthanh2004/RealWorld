package main

import (
	"Netlfy/database"
	"Netlfy/routes"
)

func main() {
	database.ConnectDB()
	
	r := routes.SetupRouter()
	r.Run(":8085")
}
