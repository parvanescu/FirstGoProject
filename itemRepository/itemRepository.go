package itemRepository

import (
	"ExGabi/interfaces"
	"ExGabi/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//mongodb+srv://UserToDoList:<password>@firstcluster.cwp9s.mongodb.net/<dbname>?retryWrites=true&w=majority

type ItemRepository struct{
	collection mongo.Collection
	ctx context.Context
}

func New(collection *mongo.Collection,context context.Context)interfaces.IRepository{
	fmt.Println("Repo initialized")
	repo := ItemRepository{*collection,context}
	return &repo
}


func (r *ItemRepository)Add(item model.ToDoItem){
	_,err :=r.collection.InsertOne(context.Background(),item)//add result to a list of user posts
	if err!=nil{
		panic(err)
	}

}

func (r *ItemRepository)Delete(id primitive.ObjectID)model.ToDoItem{
	//opts := options.FindOneAndDelete().SetProjection(bson.D{{"title", 1}, {"description", 1}})
	var deletedItem model.ToDoItem
	err :=r.collection.FindOneAndDelete(r.ctx,bson.D{{"_id",id}}).Decode(&deletedItem)
	if err !=nil{
		panic(err)
	}
	return deletedItem
}

func (r *ItemRepository)Update(id primitive.ObjectID,newItem model.ToDoItem)model.ToDoItem{
	var updatedItem model.ToDoItem
	err :=r.collection.FindOneAndReplace(r.ctx,bson.D{{"_id",id}},newItem).Decode(&updatedItem)
	if err !=nil{
		panic(err)
	}
	return updatedItem
}

func (r *ItemRepository)GetAll() []model.ToDoItem{

	cursor,err :=r.collection.Find(context.TODO(),bson.D{})
	if err != nil{
		panic(err)
	}
	defer cursor.Close(context.Background())

	var allItems []model.ToDoItem
	err = cursor.All(r.ctx,&allItems)
	if err != nil{
		panic(err)
	}
	return allItems
}
