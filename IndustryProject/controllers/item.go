package controllers

import (
	"IndustryProject/database"
	"IndustryProject/models"
	"fmt"
	"net/http"
	"strconv"
)

func AddItem(w http.ResponseWriter, r *http.Request) {
	// check if the user is logged in , and whether he is a restuarant owner
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	// get the user information using session management
	myUser := getUser(w, r)

	if r.Method == http.MethodPost {
		// Process form
		name := r.FormValue("name")
		price := r.FormValue("price")
		calories := r.FormValue("calories")
		fats := r.FormValue("fats")
		sugar := r.FormValue("sugar")

		// check if there is empty field , return a clientMsg if there is any empty field
		if name == "" || price == "" {
			fmt.Println("Field cannot be empty")
			clientMsg = "Field cannot empty"
			data := struct {
				myUser    models.User
				ClientMsg string
			}{
				myUser,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "addItem.gohtml", data)
			return
		}

		priceFloat, err := strconv.ParseFloat(price, 32)

		if err != nil {
			fmt.Println("Enter a float value for price")
			clientMsg = "Enter a float value for price"
			data := struct {
				myUser    models.User
				ClientMsg string
			}{
				myUser,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "addItem.gohtml", data)
			return
		}

		if calories == "" && fats == "" && sugar == "" {
			err = database.AddItem2(1, name, priceFloat)
			if err != nil {
				fmt.Println(err)
				clientMsg = "Internal server error at database"
			} else {
				fmt.Println("You successfully create a item")
				clientMsg = "You have successfully created a new menu item"
			}
		} else {
			caloriesFloat := 0.0
			if calories != "" {
				caloriesFloat, err = strconv.ParseFloat(calories, 32)
				if err != nil {
					fmt.Println("Enter a float value for calories")
					clientMsg = "Enter a float value for calories"
					data := struct {
						myUser    models.User
						ClientMsg string
					}{
						myUser,
						clientMsg,
					}
					tpl.ExecuteTemplate(w, "addItem.gohtml", data)
					return
				}
			}

			fatsFloat := 0.0
			if fats != "" {
				fatsFloat, err = strconv.ParseFloat(fats, 32)
				if err != nil {
					fmt.Println("Enter a float value for fats")
					clientMsg = "Enter a float value for fats"
					data := struct {
						myUser    models.User
						ClientMsg string
					}{
						myUser,
						clientMsg,
					}
					tpl.ExecuteTemplate(w, "addItem.gohtml", data)
					return
				}
			}

			sugarFloat := 0.0
			if sugar != "" {
				sugarFloat, err = strconv.ParseFloat(sugar, 32)
				if err != nil {
					fmt.Println("Enter a float value for sugar level")
					clientMsg = "Enter a float value for sugar level"
					data := struct {
						myUser    models.User
						ClientMsg string
					}{
						myUser,
						clientMsg,
					}
					tpl.ExecuteTemplate(w, "addItem.gohtml", data)
					return
				}
			}
			// currently hardcoded restaurant id , need to get the restaurant id from sesssion management
			// call database package and add the item into the database
			err = database.AddItem(1, name, priceFloat, caloriesFloat, fatsFloat, sugarFloat)

			if err != nil {
				fmt.Println(err)
				clientMsg = "Internal server error at database"
			} else {
				fmt.Println("You successfully create a item")
				clientMsg = "You have successfully created a new menu item"
			}
		}
		data := struct {
			myUser    models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "addItem.gohtml", data)
	}
	tpl.ExecuteTemplate(w, "addItem.gohtml", nil)
}
