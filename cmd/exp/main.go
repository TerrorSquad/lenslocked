package main

import (
	"database/sql"
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

	name := "Bob Calhoun"
	email := "bob@example.com"
	// This is bad! Don't do this!
	// SQL injection vulnerability
	//name = "',''); DROP TABLE users; --"
	//query := fmt.Sprintf(`
	//	INSERT INTO users (name, email)
	//   	VALUES ('%s', '%s');`, name, email)
	//fmt.Println("Executing query: " + query)
	//_, err = db.Exec(query)
	
	_, err = db.Exec(`
		INSERT INTO users (name, email)
		VALUES ($1, $2);`, name, email)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted user successfully")
}
