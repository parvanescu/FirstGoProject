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


func (r *Repository)AddItem(item payload.Item)(primitive.ObjectID,error){
	modelItem := model.NewItem(item.Title,item.Description,item.UserId)

	callback := func(sessCtx mongo.SessionContext) (interface{},error){
			var id primitive.ObjectID
			var status bool
			var err error
			if id,err  = r.insertItem(modelItem);err != nil{
				return nil,err
			}

			if status,err = r.getUserStatus(item.UserId);err != nil{
				return nil,err
			}
			if status == false{
				if err := r.setUserStatus(item.UserId,true);err !=nil{
					return nil,err
				}
			}
			return id,nil
	}

	session, err:= r.client.StartSession()
	if err !=nil{
		return [12]byte{}, err
	}
	defer session.EndSession(context.TODO())

	result,err := session.WithTransaction(context.TODO(),callback)
	if err != nil{
		return [12]byte{}, err
	}
	return result.(primitive.ObjectID),nil

}

func (r *Repository)DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)error {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	callback := func(sessCtx mongo.SessionContext) (interface{},error){
		var deletedItem response.Item
		err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id},{"userId",userId}}).Decode(&deletedItem)
		if err !=nil{
			return nil,err
		}
		user,err := r.GetUserById(deletedItem.UserId)
		if err != nil{
			return nil, err
		}
		if len(user.Items)==0{
			err:= r.setUserStatus(userId,false)
			if err != nil{
				return nil, err
			}
		}
		return nil, nil
	}
	session, err:= r.client.StartSession()
	if err !=nil{
		return err
	}
	defer session.EndSession(context.TODO())

	_,err = session.WithTransaction(context.TODO(),callback)
	if err != nil{
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

func (r *Repository)GetAllItems() (*[]response.Item,error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	cursor,err :=itemCollection.Find(context.TODO(),bson.D{})
	if err != nil{
		return nil,err
	}
	defer cursor.Close(context.TODO())

	var allItems = new([]response.Item)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		return nil,err
	}
	return allItems,nil
}

func (r *Repository)GetItemById(id primitive.ObjectID) (*response.Item,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	var item *response.Item
	err :=itemCollection.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(item)
	if err !=nil {
		return nil,err
	}
	return item,nil
}

func (r *Repository)AddUser(user payload.User)(primitive.ObjectID,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := model.NewUser(user.UserName,user.Password,user.Status)
	res,err :=userCollection.InsertOne(context.TODO(),modelUser)
	if err!=nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,nil
}
func (r *Repository)DeleteUser(id primitive.ObjectID)error{
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	callback := func(sessCtx mongo.SessionContext) (interface{},error){
		var deletedUser response.User
		err :=itemCollection.FindOneAndDelete(context.TODO(),bson.D{{"_id",id}}).Decode(&deletedUser)
		if err != nil{
			return nil, err
		}
		for _,v := range deletedUser.Items{
			err=r.DeleteItem(id,v.ItemId)
			if err != nil{
				return nil, err
			}
		}
		return nil, nil
	}
	session, err:= r.client.StartSession()
	if err !=nil{
		return err
	}
	defer session.EndSession(context.TODO())

	_,err = session.WithTransaction(context.TODO(),callback)
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
func (r *Repository)GetAllUsers() (*[]response.User,error){
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

	allItems := new([]response.User)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		return nil, err
	}

	return allItems,nil
}
func (r *Repository)GetUserById(id primitive.ObjectID) (response.User,error){
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
		return response.User{}, err
	}
	user := new([]response.User)
	err = cursor.All(context.TODO(),user)
	if err != nil{
		return response.User{}, err
	}
	return (*user)[0],nil

	//TODO: new querry
}

func (r *Repository)getUserStatus(id primitive.ObjectID)(bool,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	projection := bson.D{{"status",1}}
	var status struct{Status bool `bson:"status" json:"status"`}
	err :=userCollection.FindOne(context.TODO(),bson.D{{"_id",id}},options2.FindOne().SetProjection(projection)).Decode(&status)
	return status.Status,err
}

func (r *Repository)setUserStatus(id primitive.ObjectID, status bool)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(context.TODO(),bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}

func (r *Repository)insertItem(item model.Item)(primitive.ObjectID,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	res,err :=itemCollection.InsertOne(context.TODO(),item)
	if err != nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,err
}

//func (r *Repository)addItemTransaction(sessionContext mongo.SessionContext)(interface{},error){
//	var id primitive.ObjectID
//	var status bool
//	var err error
//	if id,err  = r.insertItem(modelItem);err != nil{
//		return nil,err
//	}
//
//	if status,err = r.GetUserStatus(item.UserId);err != nil{
//		return nil,err
//	}
//	if status == false{
//		if err := r.SetUserStatus(item.UserId,true);err !=nil{
//			return nil,err
//		}
//	}
//	return id,nil
//}