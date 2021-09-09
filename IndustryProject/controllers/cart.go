package controllers

import (
	"IndustryProject/database"
	"IndustryProject/models"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("mysession"))

func cartIndex(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	myUser := getUser(w, r)

	session, _ := store.Get(r, "mysession")

	strCart := session.Values["cart"].(string)
	var cart []models.Orders
	err := json.Unmarshal([]byte(strCart), &cart)
	if err != nil {
		fmt.Println(err)
	}

	totalPrice := total(cart)
	data := struct {
		User      models.User
		Cart      []models.Orders
		ClientMsg string
		Total     float64
	}{
		myUser,
		cart,
		clientMsg,
		totalPrice,
	}
	tpl.ExecuteTemplate(w, "cart.gohtml", data)
}

func addToCart(w http.ResponseWriter, r *http.Request) {

	inputs := r.URL.Query()["item_id"]
	item_id := inputs[0]
	ItemID, _ := strconv.ParseInt(item_id, 10, 0)
	item, _ := database.GetItem(int(ItemID))

	session, _ := store.Get(r, "mysession")
	cart := session.Values["cart"]
	if cart == nil {
		var cart []models.Orders
		cart = append(cart, models.Orders{
			Item:     item,
			Quantity: 1,
		})
		bytesCart, _ := json.Marshal(cart)
		session.Values["cart"] = string(bytesCart)
	} else {
		strCart := session.Values["cart"].(string)
		var cart []models.Orders
		json.Unmarshal([]byte(strCart), &cart)
		index := exist(int(ItemID), cart)
		if index == -1 {
			cart = append(cart, models.Orders{
				Item:     item,
				Quantity: 1,
			})
		} else {
			cart[index].Quantity = cart[index].Quantity + 1
		}
		bytesCart, _ := json.Marshal(cart)
		session.Values["cart"] = string(bytesCart)
	}
	session.Save(r, w)
	http.Redirect(w, r, "cart", http.StatusSeeOther)
}

func removeFromCart(w http.ResponseWriter, r *http.Request) {
	inputs := r.URL.Query()["item_id"]
	item_id := inputs[0]
	ItemID, _ := strconv.ParseInt(item_id, 10, 0)

	session, _ := store.Get(r, "mysession")
	strCart := session.Values["cart"].(string)
	var cart []models.Orders
	json.Unmarshal([]byte(strCart), &cart)

	index := exist(int(ItemID), cart)
	cart = remove(cart, index)

	bytesCart, _ := json.Marshal(cart)
	session.Values["cart"] = string(bytesCart)
	session.Save(r, w)
	http.Redirect(w, r, "cart", http.StatusSeeOther)
}

func checkOut(w http.ResponseWriter, r *http.Request) {
	myUser := getUser(w, r)
	clientMsg := ""

	session, _ := store.Get(r, "mysession")
	strCart := session.Values["cart"].(string)
	var cart []models.Orders
	json.Unmarshal([]byte(strCart), &cart)
	// generate a unique key for order primary key using current time and user id.
	t := time.Now()
	key := string(myUser.Email_Id) + t.Format("2006-01-02 15:04:06")
	// create a entry in order table
	// hard coded the userID , later get it from session management.
	userID := 1
	err := database.CreateOrder(key, userID, time.Now())
	if err != nil {
		fmt.Println(err)
		clientMsg = "Server Error"
		data := struct {
			User      models.User
			ClientMsg string
		}{
			myUser,
			clientMsg,
		}
		tpl.ExecuteTemplate(w, "cart.gohtml", data)
		return
	}
	// for each item in the cart , create a order details entry in the database.
	for _, item := range cart {
		err = database.CreateOrderDetails(key, item.Item.Item_ID, int(item.Quantity))
		if err != nil {
			fmt.Println(err)
			clientMsg = "Server Error"
			data := struct {
				User      models.User
				ClientMsg string
			}{
				myUser,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "cart.gohtml", data)
			return
		}
	}
	// remove everything from the cart.
	cart = nil
	bytesCart, _ := json.Marshal(cart)
	session.Values["cart"] = string(bytesCart)
	session.Save(r, w)

	clientMsg = "Your order have been placed"
	data := struct {
		User      models.User
		ClientMsg string
	}{
		myUser,
		clientMsg,
	}
	tpl.ExecuteTemplate(w, "thankyoupage.gohtml", data)
}

func exist(id int, cart []models.Orders) int {
	for i := 0; i < len(cart); i++ {
		if cart[i].Item.Item_ID == id {
			return i
		}
	}
	return -1
}

func total(cart []models.Orders) float64 {
	total := 0.0
	for _, item := range cart {
		total += item.Item.Item_Price * float64(item.Quantity)
	}
	total = math.Round(total*100) / 100
	return total
}

func remove(cart []models.Orders, index int) []models.Orders {
	return append(cart[:index], cart[index+1:]...)
}
