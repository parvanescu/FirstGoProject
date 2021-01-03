package main

import (
	"ExGabi/repository"
	"ExGabi/useCase"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main(){
	uri := "mongodb+srv://DbUser123:Password123@cluster0.5tzzs.mongodb.net/<dbname>?retryWrites=true&w=majority"
	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}
	var repo repository.IRepository = repository.New(client)
	var uC useCase.IUseCase = useCase.New(repo)

	//DbUser123:Password123
	//uC.AddUser(payload.NewUser("testUser1","testPassword1"))
	//uC.AddUser(payload.NewUser("testUser2","testPassword2"))
	//userId,err := primitive.ObjectIDFromHex("5fbe2a3decf018f23f138077")
	userId2,err := primitive.ObjectIDFromHex("5fbe3039b99d41988d9ed066")
	//itemId1,err := primitive.ObjectIDFromHex("5fbe84e3a1eb4a237bfff6f0")
	//itemId2,err := primitive.ObjectIDFromHex("5fbe84f4a1eb4a237bfff6f1")
	//itemId3,err := primitive.ObjectIDFromHex("5fbe84f7a1eb4a237bfff6f2")
	//uC.AddItem(payload.NewItem("Title1","Description1",userId))
	//uC.AddItem(payload.NewItem("Title2","Description2",userId))
	//uC.AddItem(payload.NewItem("Title3","Description3",userId))
	//uC.AddItem(payload.NewItem("Title4","Description4",userId2))
	//uC.AddItem(payload.NewItem("Title5","Description5",userId2))
	//uC.AddItem(payload.NewItem("Title6","Description6",userId2))
	//uC.DeleteUser(userId)
	//uC.UpdateUser(userId,"updatedTestUser1")
	//uC.GetAllUsers()
	//fmt.Println(uC.GetUser(userId2))
	//user,err :=uC.UpdateUser(userId2,payload.NewUser("User2","Password2"))
	//uC.DeleteItem(userId,itemId1)
	//uC.DeleteItem(userId,itemId2)
	res,err :=uC.GetUserById(userId2)
	if err != nil{
		panic(err)
	}
	fmt.Println(res)
	err =client.Disconnect(ctx)
	if err !=nil{
		panic(err)
	}
}
