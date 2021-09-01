package models

import (
	"encoding/json"
	"time"
)

type Orders struct {
	Id         string          `json:"id"`
	Create_dtm time.Time       `json:"create_dtm"`
	Order_id   string          `json:"order_id"`
	Phone      string          `json:"phone"`
	Name       string          `json:"name"`
	Address    string          `json:"address"`
	Menu       json.RawMessage `json:"menu"`
	Total_item int             `json:"total_item"`
	Pay        int             `json:"pay"`
}

type Query struct {
	Phone string `json:"phone"`
	Date  string `json:"date"`
}

type Item struct {
	Item_ID       int
	Restaurant_ID int
	Item_Name     string
	Item_Price    float64
	Calories      float64
	Fats          float64
	Sugar_Level   float64
}

type User struct {
	Email_Id   string
	Password   []byte
	First_Name string
	Last_Name  string
	Age        string
}

type Restaurant struct {
	Restaurant_ID int
}
