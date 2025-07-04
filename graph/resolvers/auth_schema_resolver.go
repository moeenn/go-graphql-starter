package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.73

import (
	"api/graph"
	"api/graph/gmodel"
	"context"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input gmodel.LoginInput) (*gmodel.LoginResponse, error) {
	return r.Service.Login(ctx, input)
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
