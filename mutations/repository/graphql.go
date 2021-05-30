package repository

import (
	"ExGabi/customErrors"
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



func New(client *mongo.Client) mutations.IRepository{
	fmt.Println("Repo initialized")
	repo := Repository{client: *client}
	return &repo
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

func (r *Repository) AddUserAndOrganisation(user *payload.User,organisation *payload.Organisation)(primitive.ObjectID,primitive.ObjectID,error){
	newUserId := primitive.ObjectID{}
	newOrganisationId := primitive.ObjectID{}
	err := r.client.UseSession(context.TODO(),func(sessCtx mongo.SessionContext)error{
		organisationId,err:=r.AddOrganisation(organisation)
		if err!=nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}

		userId,err:=r.AddUser(organisationId,user)
		if err!=nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}

		newUserId = userId
		newOrganisationId = organisationId
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return newUserId,newOrganisationId,err
}
func (r *Repository) UpdateUserPassword(user *payload.User,organisation *payload.Organisation) error {
	err := r.client.UseSession(context.TODO(),func(sessCtx mongo.SessionContext)error{
		err := r.setUserPassword(user.Id,user.Password,sessCtx)
		if err!=nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		err= r.setUserStatus(user.Id,true,sessCtx)
		if err!=nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		err = r.setOrganisationStatus(organisation.Id, true, sessCtx)
		if err!=nil{
			_ = sessCtx.AbortTransaction(sessCtx)
			return err
		}
		_ = sessCtx.CommitTransaction(sessCtx)
		return nil
	})
	return err
}

func (r *Repository) AddUser(organisationId primitive.ObjectID,user *payload.User) (primitive.ObjectID, error) {
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	modelUser := model.User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		OrganisationId: organisationId,
		Status:   false,
	}
	if user.Password!="" {
		modelUser.Password = user.Password
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
		OrganisationId: user.OrganisationId,
		Status: user.Status,
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

func (r *Repository) AddOrganisation(organisation *payload.Organisation)(primitive.ObjectID,error){
	organisationCollection := r.client.Database("ToDoApp").Collection("Organisations")
	modelOrganisation := model.Organisation{
		Name:   organisation.Name,
		CUI:    organisation.CUI,
		Status: false,
	}
	_,err:=r.GetOrganisationByCUI(organisation)
	if err == nil{
		return [12]byte{}, errors.New("an organisation with this CUI already exists")
	}
	res,err :=organisationCollection.InsertOne(context.TODO(),modelOrganisation)
	if err!=nil{
		return [12]byte{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id,nil
}
func (r *Repository) DeleteOrganisationById(organisation *payload.Organisation) error {
	panic("implement me")
}
func (r *Repository) UpdateOrganisation(organisation *payload.Organisation) (*response.Organisation, error) {
	organisationCollection := r.client.Database("ToDoApp").Collection("Organisations")
	modelOrganisation := &model.Organisation{
		Id:     organisation.Id,
		Name:  	organisation.Name,
		CUI:    organisation.CUI,
		Status: false,
	}
	updateValue := bson.M{
		"$set":modelOrganisation,
	}
	updatedOrganisation := &response.Organisation{}
	err := organisationCollection.FindOneAndUpdate(context.TODO(),bson.D{{"_id",organisation.Id}},updateValue).Decode(updatedOrganisation)
	if err!=nil{
		return nil, err
	}
	return updatedOrganisation,nil
}

func (r *Repository) AddPositionToOrganisation(position *payload.Position) (primitive.ObjectID, error) {
	positionCollection := r.client.Database("ToDoApp").Collection("Positions")
	modelPosition := model.Position{
		OrganisationId: position.OrganisationId,
		Name:           position.Name,
		AccessLevel:    position.AccessLevel,
	}

	id,err := positionCollection.InsertOne(context.TODO(),modelPosition)
	if err != nil{
		return [12]byte{}, err
	}
	return id.InsertedID.(primitive.ObjectID), nil
}

func (r *Repository) UpdatePosition(position *payload.Position) (*response.Position, error) {
	positionCollection := r.client.Database("ToDoApp").Collection("Positions")
	modelPosition := model.Position{
		Id:             position.Id,
		OrganisationId: position.OrganisationId,
		Name:           position.Name,
		AccessLevel:    position.AccessLevel,
	}
	updateValue := bson.M{
		"$set": modelPosition,
	}
	updatedPosition := &response.Position{}
	err := positionCollection.FindOneAndUpdate(context.TODO(),bson.D{{"_id",position.Id}},updateValue).Decode(updatedPosition)
	if err != nil{
		return nil, err
	}
	return updatedPosition, nil
}
func (r *Repository) ExchangePositions(position1 *payload.Position,position2 *payload.Position)(*response.Position,*response.Position,error){
	newPosition1 := &response.Position{}
	newPosition2 := &response.Position{}
	err := r.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		var err error
		newPosition1,err = r.setPositionAccessLevel(position1.Id,position2.AccessLevel,sessionContext)
		if err != nil{
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		newPosition2,err = r.setPositionAccessLevel(position2.Id,position1.AccessLevel,sessionContext)
		if err != nil{
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if err != nil{
		return nil, nil, err
	}
	return newPosition1, newPosition2, nil
}

func (r *Repository)setUserPassword(id primitive.ObjectID,password string,ctx context.Context)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"password":password}})
	return err
}
func (r *Repository)setUserStatus(id primitive.ObjectID, status bool,ctx context.Context)error{
	userCollection := r.client.Database("ToDoApp").Collection("Users")
	_,err :=userCollection.UpdateOne(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}
func (r *Repository)setOrganisationStatus(id primitive.ObjectID,status bool,ctx context.Context)error{
	organisationCollection := r.client.Database("ToDoApp").Collection("Organisations")
	_,err :=organisationCollection.UpdateOne(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"status":status}})
	return err
}
func (r *Repository)setPositionAccessLevel(id primitive.ObjectID,accessLevel int,ctx context.Context)(*response.Position,error){
	positionCollection := r.client.Database("ToDoApp").Collection("Positions")
	updatedPosition := & response.Position{}
	err := positionCollection.FindOneAndUpdate(ctx,bson.D{{"_id",id}},bson.M{"$set":bson.M{"accessLevel":accessLevel}}).Decode(updatedPosition)
	if err != nil{
		return nil, err
	}
	return updatedPosition, nil
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
		return nil,errors.New("not found")
	}
	return &(*users)[0],nil
}
func (r *Repository) GetUserByCredentials(user *payload.User) (*response.User, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Users")
	query := []bson.M{
		{"$match": bson.M{
			"email": user.Email,
			"password": user.Password,
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
		return nil,errors.New("no user with this credentials was found")
	}
	return &(*users)[0],nil
}
func (r *Repository) GetMatchingItems(userId primitive.ObjectID,item *payload.Item) (*[]response.Item, error) {
	itemCollection := r.client.Database("ToDoApp").Collection("Items")

	searchRegexValue := primitive.Regex{Pattern: item.Description, Options: "i"}
	search := bson.M{}
	if item.Description != "" {
		search["$or"] = []bson.M{{"description": searchRegexValue}, {"title": searchRegexValue}}
	}

	cursor,err :=itemCollection.Aggregate(context.TODO(),
		[]bson.M{
			{"$match": bson.M{"userId":userId}},
			{"$match": search},
		})

	if err != nil{
		return nil,err
	}
	defer cursor.Close(context.TODO())

	allItems := new([]response.Item)
	err = cursor.All(context.TODO(),allItems)
	if err != nil{
		return nil,err
	}
	return allItems,nil
}
func (r *Repository) GetOrganisationByCUI(organisation *payload.Organisation)(*response.Organisation,error){
	organisationsCollection := r.client.Database("ToDoApp").Collection("Organisations")
	query := []bson.M{
		{"$match": bson.M{
			"cui": organisation.CUI,
		}},
		{
			"$lookup":bson.M{
				"from" : "Users",
				"localField": "_id",
				"foreignField": "organisationId",
				"as": "users",
			}}}
	cursor,err:= organisationsCollection.Aggregate(context.TODO(),query)
	if err !=nil{
		return &response.Organisation{}, err
	}
	organisations := new([]response.Organisation)
	err = cursor.All(context.TODO(), organisations)
	if err != nil {
		return &response.Organisation{}, err
	}
	if len(*organisations) == 0{
		return nil,customErrors.NewOrganisationNotFoundError()
	}
	return &(*organisations)[0],nil
}