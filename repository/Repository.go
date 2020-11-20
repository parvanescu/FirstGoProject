package repository

import (
	"ExGabi/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//mongodb+srv://UserToDoList:<password>@firstcluster.cwp9s.mongodb.net/<dbname>?retryWrites=true&w=majority
type IRepository interface{
	AddItem(item model.Item)
	DeleteItem(id primitive.ObjectID)model.Item
	UpdateItem(id primitive.ObjectID,newItem model.Item)model.Item
	GetAllItems() []model.Item
}

type Repository struct{
	client mongo.Client
}

func New(client *mongo.Client)IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{*client}
	return &repo
}


func (r *Repository)AddItem(item model.Item){
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	_,err :=itemCollection.InsertOne(context.TODO(),item)//add result to a list of user posts
	if err!=nil{
		panic(err)
	}
}

func (r *Repository)DeleteItem(id primitive.ObjectID)model.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	var deletedItem model.Item
	err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id}}).Decode(&deletedItem)
	if err !=nil{
		panic(err)
	}
	return deletedItem
}

func (r *Repository)UpdateItem(id primitive.ObjectID,newItem model.Item)model.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	var updatedItem model.Item
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",id}},newItem).Decode(&updatedItem)
	if err !=nil{
		panic(err)
	}
	return updatedItem
}

func (r *Repository)GetAllItems() []model.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	cursor,err :=itemCollection.Find(context.TODO(),bson.D{})
	if err != nil{
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var allItems []model.Item
	err = cursor.All(context.TODO(),&allItems)
	if err != nil{
		panic(err)
	}
	return allItems
}
