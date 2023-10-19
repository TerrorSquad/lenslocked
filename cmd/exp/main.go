package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
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

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
		    user_id INT NOT NULL,
		    amount INT NOT NULL,
		    description TEXT
		);
	`)

	if err != nil {
		panic(err)
	}

	fmt.Println("Created tables successfully")
	//
	//name := "Alice Calhoun"
	//email := "alice2@example.com"
	// This is bad! Don't do this!
	// SQL injection vulnerability
	//name = "',''); DROP TABLE users; --"
	//query := fmt.Sprintf(`
	//	INSERT INTO users (name, email)
	//   	VALUES ('%s', '%s');`, name, email)
	//fmt.Println("Executing query: " + query)
	//_, err = db.Exec(query)

	//row := db.QueryRow(`
	//	INSERT INTO users (name, email)
	//	VALUES ($1, $2) RETURNING id;`, name, email)
	//var id int
	//err = row.Scan(&id)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("User created. ID: ", id)

	id := 1
	row := db.QueryRow(`SELECT name, email FROM users WHERE id = $1;`, id)
	var name, email string
	err = row.Scan(&name, &email)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("No user found")
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("User:", name, email)

	userID := 1
	for i := 1; i <= 5; i++ {
		amount := i * 100
		description := fmt.Sprintf("Fake order #%d", i)
		_, err = db.Exec(`
			INSERT INTO orders (user_id, amount, description)
			VALUES ($1, $2, $3);`, userID, amount, description)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Created orders successfully")
}
