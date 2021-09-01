package main

import (
	"IndustryProject/controllers"
	"IndustryProject/database"
	"net/http"
)

/*
var tpl *template.Template

var mapUsers = map[string]user{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = user{"admin", bPassword, "admin", "admin", "admin"}
}
*/
func main() {
	database.InitDB()
	controllers.ServerHTTPStarter()
	http.ListenAndServe(":8080", nil)
}
