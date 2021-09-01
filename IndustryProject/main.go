package main

import (
	"IndustryProject/controllers"
	"IndustryProject/database"
	"net/http"
)

func main() {
	database.InitDB()
	controllers.ServerHTTPStarter()
	http.ListenAndServe(":8080", nil)
}
