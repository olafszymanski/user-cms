package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/olafszymanski/user-cms/auth"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	user := auth.ForContext(ctx)
	if user != nil {
		if *user.Admin {
			return next(ctx)
		}
	}
	return nil, fmt.Errorf("access denied")
}
