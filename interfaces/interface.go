package interfaces

import (
	"ExGabi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUseCase interface {
	AddItem(title string,description string)
	DeleteItem(id primitive.ObjectID)model.ToDoItem
	UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)model.ToDoItem
	GetItem(id primitive.ObjectID)model.ToDoItem
	GetAll()//[]model.ToDoItem
}


type IRepository interface{
	Add(item model.ToDoItem)
	Delete(id primitive.ObjectID)model.ToDoItem
	Update(id primitive.ObjectID,newItem model.ToDoItem)model.ToDoItem
	GetAll() []model.ToDoItem
}

//type IToDoItem interface {
//	GetItemId() primitive.ObjectID
//	GetTitle() string
//	GetDescription() string
//}