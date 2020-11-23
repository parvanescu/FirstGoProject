package main

import (
	"ExGabi/repository"
	"ExGabi/useCase"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main(){
	uri := "mongodb+srv://UserToDoList:ToDoList123@firstcluster.cwp9s.mongodb.net/DB_ToDoItem?retryWrites=true&w=majority"
	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}
	var repo repository.IRepository = repository.New(client)
	var uC useCase.IUseCase = useCase.New(repo)

	//uC.UpdateItem(primitive.ObjectID{5fb51acf3508c7376d72597f})
	//obj,err:=primitive.ObjectIDFromHex("5fb51acf3508c7376d72597f")
	//oldItem := uC.UpdateItem(obj,"UpdatedTitle","UpdatedDescription")
	//fmt.Println("Old item: ",oldItem)
	//uC.GetAll()
	//deletedItem := uC.DeleteItem(obj)
	//fmt.Println("Deleted item:",deletedItem)
	//uC.GetAllItems()
	//uC.AddUser("testUser1","testPassword1")
	//uC.AddUser("testUser2","testPassword2")
	userId,err := primitive.ObjectIDFromHex("5fbaab6071a85619272966c2")
	//userId,err := primitive.ObjectIDFromHex("5fbb78fc6061ffef548e6183")
	//itemId,err := primitive.ObjectIDFromHex("5fbb7d2bdd5a74e9f935b1e8")
	//uC.AddItem(userId,"Title1","Description1")
	//uC.AddItem(userId,"Title2","Description2")
	//uC.AddItem(userId,"Title3","Description3")
	//uC.AddItem(userId,"Title4","Description4")
	//uC.AddItem(userId,"Title5","Description5")
	//uC.AddItem(userId,"Title6","Description6")
	//uC.AddItem(userId,"Title7","Description7")
	//uC.AddItem(userId,"Title8","Description8")
	//uC.AddItem(userId,"Title9","Description9")
	//uC.AddItem(userId,"Title10","Description10")
	//uC.DeleteItem(userId,itemId)
	//uC.DeleteUser(userId)
	uC.UpdateUser(userId,"updatedTestUser1")
	uC.GetAllUsers()
	err =client.Disconnect(ctx)
	if err !=nil{
		panic(err)
	}
}
