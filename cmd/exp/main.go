package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
	Age  int
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Doe",
		Bio:  `<script>alert('Hello')</script>`,
		Age:  30,
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
