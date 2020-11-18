package main

import (
	"ExGabi/interfaces"
	"ExGabi/itemRepository"
	"ExGabi/useCase"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main(){
	uri := "mongodb+srv://UserToDoList:ToDoList123@firstcluster.cwp9s.mongodb.net/DB_ToDoItem?retryWrites=true&w=majority"
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}
	var repo interfaces.IRepository = itemRepository.New(client.Database("DB_ToDoItem").Collection("ToDo_Collection"),ctx)
	var uC interfaces.IUseCase = useCase.New(repo)
	//uC.AddItem("Title1","Description1")
	//uC.AddItem("Title2","Description2")
	//uC.AddItem("Title3","Description3")
	//uC.AddItem("Title4","Description4")
	//uC.AddItem("Title6","Description6")
	//uC.UpdateItem(primitive.ObjectID{5fb51acf3508c7376d72597f})
	//obj,err:=primitive.ObjectIDFromHex("5fb51acf3508c7376d72597f")
	//oldItem := uC.UpdateItem(obj,"UpdatedTitle","UpdatedDescription")
	//fmt.Println("Old item: ",oldItem)
	//uC.GetAll()
	//deletedItem := uC.DeleteItem(obj)
	//fmt.Println("Deleted item:",deletedItem)
	uC.GetAll()
	err =client.Disconnect(ctx)
	if err !=nil{
		panic(err)
	}
}
