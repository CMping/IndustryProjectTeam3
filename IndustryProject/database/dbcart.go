package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func CreateOrder(orderId string, userId int, createdDT time.Time) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Orders (OrderId, UserId_fk, CreatedDT) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(orderId, userId, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func CreateOrderDetails(orderId_fk string, itemId_fk int, quantity int) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO orderdetails (orderId_fk, itemId_fk, quantity) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(orderId_fk, itemId_fk, quantity)
	if err != nil {
		return err
	}
	return nil
}
