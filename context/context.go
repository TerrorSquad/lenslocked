package context

import (
	"context"
	"github.com/terrorsquad/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	value := ctx.Value(userKey)
	user, ok := value.(*models.User)
	if !ok {
		// The most likely reason is that the value is nil because nothing has been set in the context.
		return nil
	}
	return user
}
