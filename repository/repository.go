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
	_ "go.mongodb.org/mongo-driver/mongo/options"
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


func (r *Repository)AddItem(item *payload.Item)(primitive.ObjectID,error){
	newItemId := primitive.ObjectID{}
	err := r.client.UseSession(context.TODO(), func(sessCtx mongo.SessionContext)error{
		if err := sessCtx.StartTransaction();err!=nil {
			return err
		}
		itemId,err := r.insertItem(item,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		item.Id=itemId
		user,err := r.GetUserById(item.UserId)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}

		if  user.Status== false{
			if err := r.setUserStatus(item.UserId,true,sessCtx);err !=nil{
					return err
				}
		}
		newItemId = itemId
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return newItemId,err

}

func (r *Repository)DeleteItem(item *payload.Item)error {
	err := r.client.UseSession(context.TODO(), func(sessCtx mongo.SessionContext)error{
		deletedItem,err := r.removeItem(item,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		user,err := r.GetUserById(deletedItem.UserId)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		if len(user.Items)==0{
			err:= r.setUserStatus(item.UserId,false,sessCtx)
			if err != nil{
				_ = sessCtx.AbortTransaction(sessCtx)
				return err
			}
		}
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return err
}

func (r *Repository)UpdateItem(item *payload.Item)(*response.Item,error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	modelItem := model.Item{Id: item.Id,Title: item.Title,Description: item.Description,UserId: item.UserId}
	var updatedItem *response.Item
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",item.Id}},modelItem).Decode(updatedItem)
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

func (r *Repository)AddUser(user *payload.User)(primitive.ObjectID,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := model.User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Password: user.Password,
		Status:   false,
	}
	res,err :=userCollection.InsertOne(context.TODO(),modelUser)
	if err!=nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,nil
}
func (r *Repository)DeleteUser(user *payload.User)error{
	err := r.client.UseSession(context.TODO(),func(sessCtx mongo.SessionContext)error{
		deletedUser,err := r.removeUser(user,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		for _,v := range deletedUser.Items{
			_,err = r.removeItem(&payload.Item{
				Id:     v.Id,
				UserId: v.UserId,
			},sessCtx)
			if err != nil{
				_ = sessCtx.AbortTransaction(sessCtx)
				return err
			}
		}
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return err
}
func (r *Repository)UpdateUser(user *payload.User)(*response.User,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	var updatedItem *response.User
	err :=itemCollection.FindOneAndReplace(context.TODO(),bson.D{{"_id",user.Id}},user).Decode(updatedItem)
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

	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil{
		return nil, err
	}

	return users,nil
}
func (r *Repository)GetUserById(id primitive.ObjectID) (*response.User,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
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
		return &response.User{}, err
	}
	users := new([]response.User)
	err = cursor.All(context.TODO(),users)
	if err != nil {
		return &response.User{}, err
	}
	return &(*users)[0],nil
}

func (r *Repository)GetUserByCredentials(user *payload.User)(*response.User,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{{
		"$lookup":bson.M{
			"from" : "Items",
			"localField": "_id",
			"foreignField": "userId",
			"as": "items",
		}},
		{"$match": bson.M{
			"email": user.Email,
			"password": user.Password,
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
	return &(*users)[0],nil
}

func (r *Repository)setUserStatus(id primitive.ObjectID, status bool,ctx context.Context)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}

func (r *Repository)insertItem(item *payload.Item,ctx context.Context)(primitive.ObjectID,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	res,err :=itemCollection.InsertOne(ctx,item)
	if err != nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,err
}

func (r *Repository)removeItem(item *payload.Item,ctx context.Context)(*response.Item,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	var deletedItem *response.Item
	err :=itemCollection.FindOneAndDelete(ctx,bson.D{{"_id",item.Id},{"userId",item.UserId}}).Decode(deletedItem)
	if err !=nil{
		return nil, err
	}
	return deletedItem,nil
}

func (r * Repository)removeUser(user *payload.User,ctx context.Context)(*response.User,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	var deletedUser *response.User
	err :=userCollection.FindOneAndDelete(ctx,bson.D{{"_id",user.Id}}).Decode(deletedUser)
	if err != nil{
		return nil, err
	}
	return deletedUser,nil
}