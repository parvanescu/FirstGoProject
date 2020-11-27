package repository

import (
	"ExGabi/model"
	"ExGabi/payload"
	"ExGabi/response"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
)

//mongodb+srv://UserToDoList:<password>@firstcluster.cwp9s.mongodb.net/<dbname>?retryWrites=true&w=majority

type Repository struct{
	client mongo.Client
}

func New(client *mongo.Client)IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{*client}
	return &repo
}


func (r *Repository)AddItem(item payload.Item)primitive.ObjectID{
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	modelItem := model.NewItem(item.Title,item.Description,item.UserId)
	res,err :=itemCollection.InsertOne(context.TODO(),modelItem)
	if err!=nil{
		panic(err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id
}

func (r *Repository)DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)error {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	var deletedItem response.Item
	err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id},{"userId",userId}}).Decode(&deletedItem)
	if err !=nil{
		return err
	}
	return nil
}

func (r *Repository)UpdateItem(id primitive.ObjectID,newItem payload.Item)(response.Item,error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	modelItem := model.Item{ItemId: id,Title: newItem.Title,Description: newItem.Description,UserId: newItem.UserId}
	var updatedItem response.Item
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",id}},modelItem).Decode(&updatedItem)
	return updatedItem,err
}

func (r *Repository)GetAllItems() *[]response.Item {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	cursor,err :=itemCollection.Find(context.TODO(),bson.D{})
	if err != nil{
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var allItems = new([]response.Item)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		panic(err)
	}
	return allItems
}

func (r *Repository)GetItemById(id primitive.ObjectID) *response.Item{
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	var item *response.Item
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(item)
	if err !=nil{
		panic(err)
	}
	return item
}

func (r *Repository)AddUser(user payload.User)primitive.ObjectID{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := model.NewUser(user.UserName,user.Password,user.Status)
	res,err :=userCollection.InsertOne(context.TODO(),modelUser)
	if err!=nil{
		panic(err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id
}
func (r *Repository)DeleteUser(id primitive.ObjectID)error{
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	var deletedUser model.User
	err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id}}).Decode(&deletedUser)
	if err != nil{
		return err
	}
	return nil
}
func (r *Repository)UpdateUser(id primitive.ObjectID,user payload.User)(response.User,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	var updatedItem response.User
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",id}},user).Decode(&updatedItem)
	if err !=nil{
		return updatedItem,err
	}
	return updatedItem,nil
}
func (r *Repository)GetAllUsers() *[]response.User{
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
		panic(err)
	}
	defer cursor.Close(context.TODO())

	allItems := new([]response.User)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		panic(err)
	}

	return allItems
}
func (r *Repository)GetUserById(id primitive.ObjectID) response.User{
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	//var user model.User
	//err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(&user)
	query := []bson.M{{
		"$lookup":bson.M{
			"from" : "Items",
			"localField": "_id",
			"foreignField": "userId",
			"as": "items",
		}},
		{"$match": bson.M{
		"_id": id,
		}}}
	cursor,err:= itemCollection.Aggregate(context.TODO(),query)
	if err !=nil{
		panic(err)
	}
	user := new([]response.User)
	err = cursor.All(context.TODO(),user)
	return (*user)[0]

	//TODO: new querry
}

func (r *Repository)GetUserStatus(id primitive.ObjectID)(bool,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	projection := bson.D{{"status",1}}
	var status struct{Status bool `bson:"status" json:"status"`}
	err :=userCollection.FindOne(context.TODO(),bson.D{{"_id",id}},options2.FindOne().SetProjection(projection)).Decode(&status)
	return status.Status,err
}

func (r *Repository)SetUserStatus(id primitive.ObjectID, status bool)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(context.TODO(),bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}