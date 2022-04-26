package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/olafszymanski/user-cms/graph/generated"
	"github.com/olafszymanski/user-cms/graph/model"
	"github.com/olafszymanski/user-cms/postgres"
)

func (r *mutationResolver) CreateUser(ctx context.Context, user model.NewUser) (*model.User, error) {
	return postgres.Database.CreateUser(&user)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, user model.UpdateUser) (*model.User, error) {
	return postgres.Database.UpdateUser(&user)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (bool, error) {
	deleted := true
	err := postgres.Database.DeleteUser(id)
	if err != nil {
		deleted = false
	}
	return deleted, err
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	return postgres.Database.User(id)
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return postgres.Database.Users()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
