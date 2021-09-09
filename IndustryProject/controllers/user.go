package controllers

import (
	"IndustryProject/models"
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

// TODO: (important) avoid using global maps, it's not concurrent safe
var mapUsers = map[string]models.User{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = models.User{"admin", bPassword, "admin", "admin", "admin"}
}
func ServerHTTPStarter() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/additem", addItem)
	http.HandleFunc("/getitem", getItem)
	http.HandleFunc("/updateitem", updateItem)
	http.HandleFunc("/deleteitem", deleteItem)
	http.HandleFunc("/menu", showMenu)
	http.HandleFunc("/cart", cartIndex)
	http.HandleFunc("/addtocart", addToCart)
	http.HandleFunc("/removefromcart", removeFromCart)

	http.Handle("/favicon.ico", http.NotFoundHandler())
}
func Index(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

func Signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	var myUser models.User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		emailid := req.FormValue("emailid")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		age := req.FormValue("age")
		if emailid != "" {
			// check  exist/ taken
			if _, ok := mapUsers[emailid]; ok {
				http.Error(res, "Email already taken", http.StatusForbidden)
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			mapSessions[myCookie.Value] = emailid

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = models.User{emailid, bPassword, firstname, lastname, age}
			mapUsers[emailid] = myUser
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

func Login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		emailid := req.FormValue("emailid")
		password := req.FormValue("password")
		// check if user exist with username
		myUser, ok := mapUsers[emailid]
		if !ok {
			http.Error(res, "emailid and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "emailid and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = emailid
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func getUser(res http.ResponseWriter, req *http.Request) models.User {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}

	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var myUser models.User
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username]
	return ok
}
