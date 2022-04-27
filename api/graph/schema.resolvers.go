package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/olafszymanski/user-cms/graph/generated"
	"github.com/olafszymanski/user-cms/graph/model"
	"github.com/olafszymanski/user-cms/postgres"
	"github.com/olafszymanski/user-cms/users"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user, err := postgres.Database.CreateUser(&users.User{
		Username: &input.Username,
		Email:    &input.Email,
		Password: &input.Password,
		Admin:    &input.Admin,
	})
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Admin:    user.Admin,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	user, err := postgres.Database.UpdateUser(&users.User{
		ID:       &input.ID,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Admin:    input.Admin,
	})
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Admin:    user.Admin,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (bool, error) {
	if err := postgres.Database.DeleteUser(id); err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	user, err := postgres.Database.User(id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Admin:    user.Admin,
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := postgres.Database.Users()
	if err != nil {
		return nil, err
	}
	var graphqlUsers []*model.User
	for _, user := range users {
		graphqlUsers = append(graphqlUsers, &model.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Admin:    user.Admin,
		})
	}
	return graphqlUsers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
