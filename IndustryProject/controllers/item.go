package controllers

import (
	"IndustryProject/database"
	"IndustryProject/models"
	"fmt"
	"net/http"
	"strconv"
)

// testing 123
func addItem(w http.ResponseWriter, r *http.Request) {
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
				User      models.User
				ClientMsg string
			}{
				myUser,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "addItem.gohtml", data)
			return
		}

		priceFloat, err := strconv.ParseFloat(price, 32)

		if err != nil || priceFloat <= 0 {
			fmt.Println("Enter a positive float value for price")
			clientMsg = "Enter a positive float value for price"
			data := struct {
				User      models.User
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
				if err != nil || caloriesFloat < 0 {
					fmt.Println("Enter a positive float value for calories")
					clientMsg = "Enter a positive float value for calories"
					data := struct {
						User      models.User
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
				if err != nil || fatsFloat < 0 {
					fmt.Println("Enter a positive float value for fats")
					clientMsg = "Enter a positive float value for fats"
					data := struct {
						User      models.User
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
				if err != nil || sugarFloat < 0 {
					fmt.Println("Enter a positive float value for sugar level")
					clientMsg = "Enter a positive float value for sugar level"
					data := struct {
						User      models.User
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
			restaurantID := 1
			// call database package and add the item into the database
			err = database.AddItem(restaurantID, name, priceFloat, caloriesFloat, fatsFloat, sugarFloat)

			if err != nil {
				fmt.Println(err)
				clientMsg = "Internal server error at database"
			} else {
				fmt.Println("You successfully create a item")
				clientMsg = "You have successfully created a new menu item"
			}
		}
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "addItem.gohtml", data)
		return
	}
	tpl.ExecuteTemplate(w, "addItem.gohtml", nil)
}

func showMenu(w http.ResponseWriter, r *http.Request) {

	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// need to implement 2 different website , 1 is for customer another for restaurant owner.

	myUser := getUser(w, r)

	clientMsg := ""

	inputs := r.URL.Query()["restaurant_id"]

	restaurant_id := inputs[0]

	restaurantIDInt, err := strconv.ParseInt(restaurant_id, 10, 0)

	if err != nil {
		fmt.Println(err)
		clientMsg = "Internal server error"
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "menu.gohtml", data)
		return
	}

	menu, err := database.GetMenu(int(restaurantIDInt))
	if err != nil {
		fmt.Println(err)
		clientMsg = "database error"
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "menu.gohtml", data)
		return
	}
	data := struct {
		User       models.User
		Restaurant models.User
		ClientMsg  string
		Menu       []models.Item
	}{
		myUser,
		myUser,
		clientMsg,
		menu,
	}
	tpl.ExecuteTemplate(w, "menu.gohtml", data)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(w, r)

	clientMsg := ""

	inputs := r.URL.Query()["item_id"]

	item_id := inputs[0]

	Item_ID, err := strconv.ParseInt(item_id, 10, 0)

	if err != nil {
		fmt.Println(err)
		clientMsg = "Server error at converting int"
		updateError(w, myUser, int(Item_ID), clientMsg)
		return
	}

	// process update form
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		price := r.FormValue("price")
		calories := r.FormValue("calories")
		fats := r.FormValue("fats")
		sugar := r.FormValue("sugar")

		// check if there is empty field , return a clientMsg if there is any empty field
		if name == "" || price == "" {
			fmt.Println("Field cannot be empty")
			clientMsg = "Field cannot empty"
			updateError(w, myUser, int(Item_ID), clientMsg)
			return
		}
		priceFloat, err := strconv.ParseFloat(price, 32)

		if err != nil || priceFloat <= 0 {
			fmt.Println("Enter a positive float value for price")
			clientMsg = "Enter a positive float value for price"
			updateError(w, myUser, int(Item_ID), clientMsg)
			return
		}

		caloriesFloat := 0.0
		if calories != "" {
			caloriesFloat, err = strconv.ParseFloat(calories, 32)
			if err != nil || caloriesFloat < 0 {
				fmt.Println("Enter a positive float value for calories")
				clientMsg = "Enter a positive float value for calories"
				updateError(w, myUser, int(Item_ID), clientMsg)
				return
			}
		}

		fatsFloat := 0.0
		if fats != "" {
			fatsFloat, err = strconv.ParseFloat(fats, 32)
			if err != nil || fatsFloat < 0 {
				fmt.Println("Enter a positive float value for fats")
				clientMsg = "Enter a positive float value for fats"
				updateError(w, myUser, int(Item_ID), clientMsg)
				return
			}
		}

		sugarFloat := 0.0
		if sugar != "" {
			sugarFloat, err = strconv.ParseFloat(sugar, 32)
			if err != nil || sugarFloat < 0 {
				fmt.Println("Enter a positive float value for sugar level")
				clientMsg = "Enter a positive float value for sugar level"
				updateError(w, myUser, int(Item_ID), clientMsg)
				return
			}
		}

		// currently hard coded the restaurant id , later will get it by session
		restaurantID := 1
		err = database.UpdateItem(int(Item_ID), restaurantID, name, priceFloat, caloriesFloat, fatsFloat, sugarFloat)

		if err != nil {
			fmt.Print(err)
			clientMsg = "Error updating item"
		} else {
			fmt.Printf("Succesfully updated item id: %d\n", Item_ID)
			clientMsg = "Item updated"
			item, err := database.GetItem(int(Item_ID))

			if err != nil {
				fmt.Println(err)
				clientMsg = "Server Error at getting new item details"
				data := struct {
					User      models.User
					ClientMsg string
				}{
					myUser,
					clientMsg,
				}
				tpl.ExecuteTemplate(w, "getItem.gohtml", data)
				return
			}

			data := struct {
				User      models.User
				Item      models.Item
				ClientMsg string
			}{
				myUser,
				item,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "getItem.gohtml", data)
			return
		}
	}

	// Get the current item details in the database and send back to the user
	item, err := database.GetItem(int(Item_ID))
	if err != nil {
		updateError(w, myUser, int(Item_ID), "Error at retrieving item details")
		return
	}
	data := struct {
		User      models.User
		Item      models.Item
		ClientMsg string
	}{
		myUser,
		item,
		clientMsg,
	}
	tpl.ExecuteTemplate(w, "updateItem.gohtml", data)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	myUser := getUser(w, r)

	clientMsg := ""

	inputs := r.URL.Query()["item_id"]

	item_id := inputs[0]

	Item_ID, err := strconv.ParseInt(item_id, 10, 0)

	// currently hard coded the restaurant id , later get it from session management
	RestaurantID := 1

	if err != nil {
		fmt.Println(err)
		clientMsg = "Internal server error"
		deleteRouting(w, myUser, RestaurantID, clientMsg)
		return
	}

	err = database.DeleteItem(int(Item_ID), RestaurantID)

	if err != nil {
		fmt.Println(err)
		clientMsg = "Internal server error"
		deleteRouting(w, myUser, RestaurantID, clientMsg)
		return
	} else {
		fmt.Printf("Succesfully deleted item id: %d\n", Item_ID)
		clientMsg = "Successfully deleted a item"
		deleteRouting(w, myUser, RestaurantID, clientMsg)
		return
	}
}

func getItem(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// need to implement 2 different website , 1 is for customer another for restaurant owner.

	myUser := getUser(w, r)

	clientMsg := ""

	inputs := r.URL.Query()["item_id"]

	item_id := inputs[0]

	ItemID, err := strconv.ParseInt(item_id, 10, 0)

	if err != nil {
		fmt.Println(err)
		clientMsg = "Internal server error"
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "getItem.gohtml", data)
		return
	}

	item, err := database.GetItem(int(ItemID))
	if err != nil {
		fmt.Println(err)
		clientMsg = "database error"
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "getItem.gohtml", data)
		return
	}
	data := struct {
		User       models.User
		Restaurant models.User
		ClientMsg  string
		Item       models.Item
	}{
		myUser,
		myUser,
		clientMsg,
		item,
	}
	tpl.ExecuteTemplate(w, "getItem.gohtml", data)
}

func updateError(w http.ResponseWriter, user models.User, itemID int, errMsg string) {
	item, err := database.GetItem(int(itemID))
	if err != nil {
		fmt.Println(err)
		errMsg = "database error"
	}
	data := struct {
		User      models.User
		Item      models.Item
		ClientMsg string
	}{
		user,
		item,
		errMsg,
	}
	tpl.ExecuteTemplate(w, "updateItem.gohtml", data)
}

func deleteRouting(w http.ResponseWriter, user models.User, restaurantID int, errMsg string) {
	menu, err := database.GetMenu(restaurantID)
	if err != nil {
		fmt.Println(err)
		errMsg = "database error"
	}
	data := struct {
		Restaurant models.User
		Menu       []models.Item
		ClientMsg  string
	}{
		user,
		menu,
		errMsg,
	}
	tpl.ExecuteTemplate(w, "menu.gohtml", data)
}
