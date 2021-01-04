package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"ExGabi/graph/generated"
	"ExGabi/graph/model"
	"ExGabi/payload"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) AddItem(ctx context.Context, item model.NewItem) (*model.Item, error) {
	payloadUserId,err:= primitive.ObjectIDFromHex(item.UserID)
	if err!=nil{
		return nil,err
	}
	newItem := payload.Item{
		Title:       item.Title,
		Description: item.Description,
		UserId:      payloadUserId,
		Token:       item.Token,
	}
	responseItem,err:=r.UseCase.AddItem(&newItem)
	if err!=nil{
		return nil, err
	}

	modelResponseItem:=model.Item{
		ID:          responseItem.Id.Hex(),
		Title:       responseItem.Title,
		Description: responseItem.Description,
		UserID:      responseItem.UserId.Hex(),
		Token:       responseItem.Token,
	}
	return &modelResponseItem,nil

}

func (r *mutationResolver) DeleteItem(ctx context.Context, item model.NewItem) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateItem(ctx context.Context, item model.NewItem) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, user model.NewUser) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, user model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Register(ctx context.Context, user model.NewUser) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, user model.NewUser) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItemByID(ctx context.Context, item model.NewItem) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllItems(ctx context.Context) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserByID(ctx context.Context, user model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

