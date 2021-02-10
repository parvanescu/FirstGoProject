package repository

import (
	"ExGabi/payload"
	"ExGabi/queries"
	"ExGabi/response"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct{
	client mongo.Client
}


func (r *Repository) GetAllItems() (*[]response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	cursor,err :=itemCollection.Find(context.TODO(),bson.D{})
	if err != nil{
		return nil,err
	}
	defer cursor.Close(context.TODO())

	allItems := new([]response.Item)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		return nil,err
	}
	return allItems,nil
}
func (r *Repository) GetItemByTitle(item *payload.Item) (*response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	foundItem := &response.Item{}
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"title",item.Title}}).Decode(foundItem)
	if err !=nil {
		return nil,err
	}
	return foundItem,nil
}
func (r *Repository) GetItemByDescription(item *payload.Item) (*[]response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	cursor,err :=itemCollection.Find(context.TODO(),
		bson.D{{"description",bson.D{
			{"$regex",item.Description},
			{"$options","i"}}}})
	if err != nil{
		return nil,err
	}
	defer cursor.Close(context.TODO())

	allItems := new([]response.Item)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		return nil,err
	}
	return allItems,nil
}
func (r *Repository) GetItemById(id primitive.ObjectID) (*response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	item := &response.Item{}
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(item)
	if err !=nil {
		return nil,err
	}
	return item,nil
}


func (r *Repository) GetAllUsers() (*[]response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{{
		"$lookup":bson.M{
			"from" : "Items",
			"localField": "_id",
			"foreignField": "userId",
			"as": "items",
		}}}
	cursor,err :=itemCollection.Aggregate(context.TODO(),query)
	if err != nil{
		return nil, err
	}
	defer cursor.Close(context.TODO())

	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil{
		return nil, err
	}

	return users,nil
}
func (r *Repository) GetUserById(id primitive.ObjectID) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{
			"$lookup":bson.M{
				"from" : "Items",
				"localField": "_id",
				"foreignField": "userId",
				"as": "items",
			}}}
	cursor,err:= itemCollection.Aggregate(context.TODO(),query)
	if err !=nil{
		return &response.User{}, err
	}
	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil {
		return &response.User{}, err
	}
	if len(*users) == 0{
		return nil,errors.New("no user with this id was found")
	}
	return &(*users)[0],nil
}
func (r *Repository) GetUserByCredentials(user *payload.User) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"email": user.Email,
			"password": user.Password,
		}},
		{
			"$lookup":bson.M{
				"from" : "Items",
				"localField": "_id",
				"foreignField": "userId",
				"as": "items",
			}}}
	cursor,err:= itemCollection.Aggregate(context.TODO(),query)
	if err !=nil{
		return &response.User{}, err
	}
	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil {
		return &response.User{}, err
	}
	if len(*users) == 0{
		return nil,errors.New("no user with this id was found")
	}
	return &(*users)[0],nil
}
func (r *Repository) GetUserByEmail(user *payload.User) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"email": user.Email,
		}},
		{
			"$lookup":bson.M{
				"from" : "Items",
				"localField": "_id",
				"foreignField": "userId",
				"as": "items",
			}}}
	cursor,err:= itemCollection.Aggregate(context.TODO(),query)
	if err !=nil{
		return &response.User{}, err
	}
	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil {
		return &response.User{}, err
	}
	if len(*users) == 0{
		return nil,errors.New("no user with this id was found")
	}
	return &(*users)[0],nil
}


func New(client *mongo.Client)queries.IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{*client}
	return &repo
}
