package main

import (
	//This is an alias for the standard library context package.
	stdctx "context"
	"fmt"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/models"
)

func main() {
	ctx := stdctx.Background()
	user := models.User{
		Email: "g.ninkovic@angeltech.rs",
	}
	ctx = context.WithUser(ctx, &user)

	retreivedUser := context.User(ctx)

	fmt.Println(retreivedUser)
}
