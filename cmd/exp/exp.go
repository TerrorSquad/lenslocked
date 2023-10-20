package main

import (
	"fmt"
	"github.com/terrorsquad/lenslocked/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
		    name TEXT,
		    email TEXT NOT NULL UNIQUE
		);
	
-- 		CREATE TABLE IF NOT EXISTS orders (
-- 			id SERIAL PRIMARY KEY,
-- 		    user_id INT NOT NULL,
-- 		    amount INT NOT NULL,
-- 		    description TEXT
-- 		);
	`)

	if err != nil {
		panic(err)
	}
	//
	fmt.Println("Created tables successfully")
	//
	name := "Alice Calhoun"
	email := "alice2@example.com"
	// This is bad! Don't do this!
	// SQL injection vulnerability
	//name = "',''); DROP TABLE users; --"
	//query := fmt.Sprintf(`
	//	INSERT INTO users (name, email)
	//   	VALUES ('%s', '%s');`, name, email)
	//fmt.Println("Executing query: " + query)
	//_, err = db.Exec(query)

	row := db.QueryRow(`
		INSERT INTO users (name, email)
		VALUES ($1, $2) RETURNING id;`, name, email)
	var id int
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}

	fmt.Println("User created. ID: ", id)

	//id := 1
	//row := db.QueryRow(`SELECT name, email FROM users WHERE id = $1;`, id)
	//var name, email string
	//err = row.Scan(&name, &email)
	//if errors.Is(err, sql.ErrNoRows) {
	//	fmt.Println("No user found")
	//}
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("User:", name, email)

	//userID := 1
	//for i := 1; i <= 5; i++ {
	//	amount := i * 100
	//	description := fmt.Sprintf("Fake order #%d", i)
	//	_, err = db.Exec(`
	//		INSERT INTO orders (user_id, amount, description)
	//		VALUES ($1, $2, $3);`, userID, amount, description)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//fmt.Println("Created orders successfully")

	//type Order struct {
	//	ID          int
	//	userID      int
	//	Amount      int
	//	Description string
	//}
	//
	//var orders []Order
	//userID := 1
	//rows, err := db.Query(`
	//	SELECT id, amount, description
	//	FROM orders
	//	WHERE user_id = $1;`, userID)
	//if err != nil {
	//	panic(err)
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var o Order
	//	o.userID = userID
	//	err = rows.Scan(&o.ID, &o.Amount, &o.Description)
	//	if err != nil {
	//		panic(err)
	//	}
	//	orders = append(orders, o)
	//}
	//err = rows.Err()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Orders:", orders)

}
