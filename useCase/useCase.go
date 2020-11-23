package useCase

import (
	"ExGabi/model"
	"ExGabi/repository"
	"ExGabi/response"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUseCase interface {
	AddItem(userId primitive.ObjectID,title string,description string)
	DeleteItem(userId primitive.ObjectID,itemId primitive.ObjectID)response.Item
	UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)response.Item
	GetItem(id primitive.ObjectID)response.Item
	GetAllItems() //[]model.Item

	AddUser(name string,password string)
	DeleteUser(id primitive.ObjectID)response.User
	UpdateUser(id primitive.ObjectID,newName string)response.User
	GetUser(id primitive.ObjectID)response.User
	GetAllUsers() //[]response.User
}

type UseCase struct{
	itemRepository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)AddItem(userId primitive.ObjectID,title string,description string){
	var item model.Item = model.NewItem(title,description)
	uC.itemRepository.AddItem(userId,item)
}

func (uC *UseCase)DeleteItem(userId primitive.ObjectID,itemId primitive.ObjectID)response.Item {
	return uC.itemRepository.DeleteItem(userId,itemId)
}

func (uC *UseCase)UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)response.Item {
	return uC.itemRepository.UpdateItem(id,model.NewItem(newTitle,newDescription))
}

func(uC *UseCase)GetItem(id primitive.ObjectID)response.Item {
	
	for _,v := range *uC.itemRepository.GetAllItems(){
		if v.ItemId == id{
			return v
		}
	}
	return response.NewItem("","")
}

func(uC *UseCase)GetAllItems() {
	items := *uC.itemRepository.GetAllItems()
	for _, v := range items {
		fmt.Println(v.ItemId.Hex())
		fmt.Println(v.Title)
		fmt.Println(v.Description)
	}
}

func(uC *UseCase) AddUser(name string,password string){
	newUser := model.NewUser(name,password)
	uC.itemRepository.AddUser(newUser) //returns object id of user
	}

func(uC *UseCase)DeleteUser(id primitive.ObjectID)response.User {
	return uC.itemRepository.DeleteUser(id)
	}

func(uC *UseCase)UpdateUser(id primitive.ObjectID,newName string)response.User {
	return uC.itemRepository.UpdateUser(id,newName)
	}

func(uC *UseCase)GetUser(id primitive.ObjectID)response.User {
	return uC.itemRepository.GetUserById(id)
	}

func(uC *UseCase)GetAllUsers() {
	users := *uC.itemRepository.GetAllUsers()
	for _,v := range users  {
		fmt.Println(v.UserId.Hex())
		fmt.Println(v.UserName)
		fmt.Println(v.ItemList)
	}
	}

