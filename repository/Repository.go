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
	Add(item model.Item)
	Delete(id primitive.ObjectID)model.Item
	Update(id primitive.ObjectID,newItem model.Item)model.Item
	GetAll() []model.Item
}


type ItemRepository struct{
	collection mongo.Collection
	ctx context.Context
}

func New(collection *mongo.Collection,context context.Context)IRepository{
	fmt.Println("Repo initialized")
	repo := ItemRepository{*collection,context}
	return &repo
}


func (r *ItemRepository)Add(item model.Item){
	_,err :=r.collection.InsertOne(context.Background(),item)//add result to a list of user posts
	if err!=nil{
		panic(err)
	}

}

func (r *ItemRepository)Delete(id primitive.ObjectID)model.Item {
	//opts := options.FindOneAndDelete().SetProjection(bson.D{{"title", 1}, {"description", 1}})
	var deletedItem model.Item
	err :=r.collection.FindOneAndDelete(r.ctx,bson.D{{"_id",id}}).Decode(&deletedItem)
	if err !=nil{
		panic(err)
	}
	return deletedItem
}

func (r *ItemRepository)Update(id primitive.ObjectID,newItem model.Item)model.Item {
	var updatedItem model.Item
	err :=r.collection.FindOneAndReplace(r.ctx,bson.D{{"_id",id}},newItem).Decode(&updatedItem)
	if err !=nil{
		panic(err)
	}
	return updatedItem
}

func (r *ItemRepository)GetAll() []model.Item {

	cursor,err :=r.collection.Find(context.TODO(),bson.D{})
	if err != nil{
		panic(err)
	}
	defer cursor.Close(context.Background())

	var allItems []model.Item
	err = cursor.All(r.ctx,&allItems)
	if err != nil{
		panic(err)
	}
	return allItems
}
