package userRepository

import (
	"ExGabi/interfaces"
	"ExGabi/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection mongo.Collection
	ctx context.Context
}

func (u UserRepository) Add(user model.User) {
	panic("implement me")
}

func (u UserRepository) Delete(id primitive.ObjectID) model.ToDoItem {
	panic("implement me")
}

func (u UserRepository) Update(id primitive.ObjectID, newItem model.ToDoItem) model.ToDoItem {
	panic("implement me")
}

func (u UserRepository) GetAll() []model.ToDoItem {
	panic("implement me")
}

func New(collection *mongo.Collection,context context.Context)interfaces.IRepository{
	fmt.Println("Repo initialized")
	repo := UserRepository{*collection,context}
	return &repo
}
