package repository

import (
	"ExGabi/model"
	"ExGabi/response"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//mongodb+srv://UserToDoList:<password>@firstcluster.cwp9s.mongodb.net/<dbname>?retryWrites=true&w=majority
type IRepository interface{
	AddItem(userId primitive.ObjectID,item model.Item)primitive.ObjectID
	DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)response.Item
	UpdateItem(id primitive.ObjectID,newItem model.Item)response.Item
	GetAllItems() *[]response.Item
	GetItemById(id primitive.ObjectID) *response.Item

	AddUser(user model.User)primitive.ObjectID
	DeleteUser(id primitive.ObjectID)response.User
	UpdateUser(id primitive.ObjectID,newUserName string)response.User
	GetAllUsers() *[]response.User
	GetUserById(id primitive.ObjectID) response.User
}

type Repository struct{
	client mongo.Client
}

func New(client *mongo.Client)IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{*client}
	return &repo
}


func (r *Repository)AddItem(userId primitive.ObjectID,item model.Item)primitive.ObjectID{
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	res,err :=itemCollection.InsertOne(context.TODO(),item)
	if err!=nil{
		panic(err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	r.appendUserItem(model.NewUserItem(userId,id))
	return id
}

func (r *Repository)DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)response.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	var deletedItem response.Item
	err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id}}).Decode(&deletedItem)
	if err !=nil{
		panic(err)
	}
	r.deleteUserItem(model.NewUserItem(userId,id))
	return deletedItem
}

func (r *Repository)UpdateItem(id primitive.ObjectID,newItem model.Item)response.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	var updatedItem response.Item
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",id}},newItem).Decode(&updatedItem)
	if err !=nil{
		panic(err)
	}
	return updatedItem
}

func (r *Repository)GetAllItems() *[]response.Item {
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
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
	itemCollection := r.client.Database("DB_ToDoItem").Collection("ToDo_Collection")
	var item *response.Item
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(item)
	if err !=nil{
		panic(err)
	}
	return item
}

func (r *Repository)AddUser(user model.User)primitive.ObjectID{
	userCollection := r.client.Database("DB_ToDoItem").Collection("User_Collection")
	res,err :=userCollection.InsertOne(context.TODO(),user)
	if err!=nil{
		panic(err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id
}
func (r *Repository)DeleteUser(id primitive.ObjectID)response.User{
	itemCollection := r.client.Database("DB_ToDoItem").Collection("User_Collection")
	var deletedUser model.User
	err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id}}).Decode(&deletedUser)
	if err != nil{
		panic(err)
	}

	userItems := r.getUserItems(deletedUser.UserId)
	var itemsIds []primitive.ObjectID
	for _,v := range *userItems{
		itemsIds = append(itemsIds,v.ItemId)
	}
	for _,v := range *userItems{
		r.DeleteItem(deletedUser.UserId,v.ItemId)
	}
	//for _,v := range *r.getUserItems(deletedUser.UserId){
	//	r.deleteUserItem(v)
	//}
	return response.NewUser(deletedUser,itemsIds)
}
func (r *Repository)UpdateUser(id primitive.ObjectID,newUserName string)response.User{
	itemCollection := r.client.Database("DB_ToDoItem").Collection("User_Collection")
	var updatedItem model.User
	err :=itemCollection.FindOneAndUpdate(context.TODO(),bson.D{{"_id",id}},
	bson.D{{"$set", bson.D{{"name", newUserName}}}}).Decode(&updatedItem)
	if err !=nil{
		panic(err)
	}
	userItems := r.getUserItems(updatedItem.UserId)
	var itemsIds []primitive.ObjectID
	for _,v := range *userItems{
		itemsIds = append(itemsIds,v.ItemId)
	}
	return response.NewUser(updatedItem,itemsIds)
}
func (r *Repository)GetAllUsers() *[]response.User{
	itemCollection := r.client.Database("DB_ToDoItem").Collection("User_Collection")
	cursor,err :=itemCollection.Find(context.TODO(),bson.D{})
	if err != nil{
		panic(err)
	}
	defer cursor.Close(context.TODO())

	allItems := new([]model.User)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		panic(err)
	}
	if allItems != nil {
		users := new([]response.User)
		for _, user := range *allItems {
			userItems := r.getUserItems(user.UserId)
			var itemsIds []primitive.ObjectID
			for _, v := range *userItems {
				itemsIds = append(itemsIds, v.ItemId)
			}
			*users = append(*users,response.NewUser(user,itemsIds))
		}
		return users
	}
	return nil
}
func (r *Repository)GetUserById(id primitive.ObjectID) response.User{
	itemCollection := r.client.Database("DB_ToDoItem").Collection("User_Collection")
	var user model.User
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(&user)
	if err !=nil{
		panic(err)
	}
	userItems := r.getUserItems(user.UserId)
	var itemsIds []primitive.ObjectID
	for _,v := range *userItems{
		itemsIds = append(itemsIds,v.ItemId)
	}
	return response.NewUser(user,itemsIds)
}


func (r *Repository)appendUserItem(userItem model.UserItem){
	userItemCollection := r.client.Database("DB_ToDoItem").Collection("UserItemRelation")
	_,err :=userItemCollection.InsertOne(context.TODO(),userItem)
	if err != nil{
		panic(err)
	}
}

func (r *Repository)deleteUserItem(userItem model.UserItem){
	userItemCollection := r.client.Database("DB_ToDoItem").Collection("UserItemRelation")
	var deletedItem model.UserItem
	err :=userItemCollection.FindOneAndDelete(context.TODO(),bson.D{{"userId",userItem.UserId},{"itemId",userItem.ItemId}}).Decode(&deletedItem)
	if err != nil{
		panic(err)
	}

}

func (r *Repository)getUserItems(userId primitive.ObjectID)*[]model.UserItem{
	userItemRelation := r.client.Database("DB_ToDoItem").Collection("UserItemRelation")
	cursor,err := userItemRelation.Find(context.TODO(),bson.D{{"userId",userId}})
	if err!=nil{
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var allItems = new([]model.UserItem)
	err = cursor.All(context.TODO(),allItems)
	if err!=nil{
		panic(err)
	}
	return allItems
}
