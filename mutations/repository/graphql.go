package repository

import (
	"ExGabi/model"
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct{
	client mongo.Client
}

func (r *Repository) AddItem(userId primitive.ObjectID, item *payload.Item) (primitive.ObjectID, error) {
	newItemId := primitive.ObjectID{}
	err := r.client.UseSession(context.TODO(), func(sessCtx mongo.SessionContext)error{
		if err := sessCtx.StartTransaction();err!=nil {
			return err
		}
		dbItem := &model.Item{Title: item.Title,Description: item.Description,UserId: userId}
		itemId,err := r.insertItem(dbItem,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		user,err := r.GetUserById(userId)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}

		if  user.Status== false{
			if err := r.setUserStatus(userId,true,sessCtx);err !=nil{
				return err
			}
		}
		newItemId = itemId
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return newItemId,err

}
func (r *Repository) DeleteItem(userId primitive.ObjectID, item *payload.Item) error {
	err := r.client.UseSession(context.TODO(), func(sessCtx mongo.SessionContext)error{
		_,err := r.removeItem(userId,item,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		user,err := r.GetUserById(userId)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		if len(user.Items)==0{
			err:= r.setUserStatus(userId,false,sessCtx)
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
func (r *Repository) UpdateItem(userId primitive.ObjectID, item *payload.Item) (*response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	modelItem := model.Item{Id: item.Id,Title: item.Title,Description: item.Description,UserId: userId}
	updateValue := bson.M{
		"$set": modelItem,
	}
	updatedItem :=&response.Item{}
	err :=itemCollection.FindOneAndUpdate(context.TODO(),bson.D{{"_id",item.Id}},updateValue).Decode(updatedItem)
	return updatedItem,err
}



func (r *Repository) AddUser(user *payload.User) (primitive.ObjectID, error) {
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := model.User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Password: user.Password,
		Status:   false,
	}
	_,err := r.GetUserByEmail(user)
	if err==nil{
		return [12]byte{},errors.New("email already used")
	}
	res,err :=userCollection.InsertOne(context.TODO(),modelUser)
	if err!=nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,nil
}
func (r *Repository) DeleteUser(user *payload.User) error {
	err := r.client.UseSession(context.TODO(),func(sessCtx mongo.SessionContext)error{
		deletedUser,err := r.removeUser(user,sessCtx)
		if err != nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		for _,v := range deletedUser.Items{
			_,err = r.removeItem(user.Id,&payload.Item{
				Id:     v.Id,
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
func (r *Repository) UpdateUser(user *payload.User) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := &model.User{
		Id:        user.Id,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		Email:     user.Email,
		Password:  user.Password,
		Status:    false,
	}
	updateValue := bson.M{
		"$set": modelUser,
	}
	updatedItem := &response.User{}
	err :=itemCollection.FindOneAndUpdate(context.TODO(),bson.D{{"_id",user.Id}},updateValue).Decode(updatedItem)
	if err !=nil{
		return updatedItem,err
	}
	return updatedItem,nil
}


func New(client *mongo.Client) mutations.IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{client: *client}
	return &repo
}

func (r *Repository)setUserStatus(id primitive.ObjectID, status bool,ctx context.Context)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}
func (r *Repository)insertItem(item *model.Item,ctx context.Context)(primitive.ObjectID,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	res,err :=itemCollection.InsertOne(ctx,item)
	if err != nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,err
}
func (r *Repository)removeItem(userId primitive.ObjectID,item *payload.Item,ctx context.Context)(*response.Item,error){
	itemCollection := r.client.Database("ToDoApp").Collection("Items")
	deletedItem := &response.Item{}
	err :=itemCollection.FindOneAndDelete(ctx,bson.D{{"_id",item.Id},{"userId",userId}}).Decode(deletedItem)
	if err !=nil{
		return nil, err
	}
	return deletedItem,nil
}
func (r * Repository)removeUser(user *payload.User,ctx context.Context)(*response.User,error){
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	deletedUser,err := r.GetUserById(user.Id)
	if err != nil{
		return nil, err
	}
	userCollection.FindOneAndDelete(ctx,bson.D{{"_id",user.Id}})
	if err != nil{
		return nil, err
	}
	return deletedUser,nil
}

func (r *Repository) GetUserById(id primitive.ObjectID) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{
			"$lookup":bson.M{
				"from" : "Items",
				"localField": "_id",
				"foreignField": "userId",
				"as": "items",
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
	if len(*users) == 0{
		return nil,errors.New("no user with this id was found")
	}
	return &(*users)[0],nil
}
func (r *Repository) GetUserByEmail(user *payload.User) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"email": user.Email,
		}},
		{
			"$lookup":bson.M{
				"from" : "Items",
				"localField": "_id",
				"foreignField": "userId",
				"as": "items",
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
	if len(*users) == 0{
		return nil,errors.New("no user with this id was found")
	}
	return &(*users)[0],nil
}
