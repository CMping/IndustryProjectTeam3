package models

type Orders struct {
	Item     Item  `json:"item"`
	Quantity int64 `json:"quantity"`
}

// TODO: (nit) golang does not use snake_case, only camelCase
type Item struct {
	Item_ID       int     `json:"item_id"`
	Restaurant_ID int     `json:"restaurant_id"`
	Item_Name     string  `json:"item_name"`
	Item_Price    float64 `json:"item_price"`
	Calories      float64 `json:"calories"`
	Fats          float64 `json:"fats"`
	Sugar_Level   float64 `json:"sugar_level"`
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
