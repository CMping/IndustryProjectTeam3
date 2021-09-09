package database

import (
	"IndustryProject/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// AddItem implements the sql operations to add a new item into the database.
func AddItem(Restaurant_ID int, Item_Name string, Item_Price float64, Calories float64, Fats float64, Sugar_Level float64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Items (Restaurant_ID_fk, Item_Name, Item_Price, Calories, Fats, Sugar_Level) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Restaurant_ID, Item_Name, Item_Price, Calories, Fats, Sugar_Level)
	if err != nil {
		return err
	}
	return nil
}

// AddItem2 implements the sql operations to add a new item into the database when the user didnt input any nutritional value.
func AddItem2(Restaurant_ID int, Item_Name string, Item_Price float64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Items (Restaurant_ID_fk, Item_Name, Item_Price) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Restaurant_ID, Item_Name, Item_Price)
	if err != nil {
		return err
	}
	return nil
}

// GetMenu implements the sql operations to retrieve all the item under the Restaurant_ID and return a menu(array of item) if no error.
func GetMenu(Restaurant_ID int) ([]models.Item, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()
	var menu []models.Item

	results, err := DB.Query("SELECT Item_ID, Item_Name, Item_Price, Calories, Fats,"+
		"Sugar_Level FROM Items WHERE Restaurant_ID_fk = ?", Restaurant_ID)

	if err != nil {
		return nil, err
	} else {
		// TODO: (nit) unnecessary else
		for results.Next() {
			var item models.Item
			err := results.Scan(&item.Item_ID, &item.Item_Name, &item.Item_Price,
				&item.Calories, &item.Fats, &item.Sugar_Level)
			if err != nil {
				return nil, err
			}
			menu = append(menu, item)
		}
		return menu, nil
	}
}

// GetItem implements the sql operations to retrieve all the item information under the Item_ID and return a struct(item) if no error.
func GetItem(Item_ID int) (models.Item, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var item models.Item

	results, err := DB.Query("SELECT Item_ID, Restaurant_ID_fk, Item_Name, Item_Price, Calories, Fats,"+
		"Sugar_Level FROM Items WHERE Item_ID = ?", Item_ID)

	if err != nil {
		return item, err
	} else {
		// TODO: (nit) unnecessary else
		for results.Next() {
			err := results.Scan(&item.Item_ID, &item.Restaurant_ID, &item.Item_Name, &item.Item_Price,
				&item.Calories, &item.Fats, &item.Sugar_Level)
			if err != nil {
				return item, err
			}
		}
		return item, nil
	}
}

// UpdateItem implements the sql operations to update the item details.
func UpdateItem(Item_ID int, Restaurant_ID int, Item_Name string, Item_Price float64, Calories float64, Fats float64, Sugar_Level float64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Items SET Item_Name=?,Item_Price=?, Calories=?, Fats=?, Sugar_Level=? WHERE Item_ID=? AND Restaurant_ID_fk=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Item_Name, Item_Price, Calories, Fats, Sugar_Level, Item_ID, Restaurant_ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem implements the sql operations to delete the recipient using repID and RecipientID.
func DeleteItem(Item_ID int, Restaurant_ID int) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Items WHERE Item_ID=? AND Restaurant_ID_fk=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Item_ID, Restaurant_ID)
	if err != nil {
		return err
	}
	return nil
}
