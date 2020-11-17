package interfaces

type IUseCase interface {
	AddItem(id int,title string,description string)
	DeleteItem(id int)IToDoItem
	UpdateItem(id int,newTitle string,newDescription string)IToDoItem
	GetItem(id int)IToDoItem
	GetAll()//[]model.ToDoItem
}


type IRepository interface{
	Add(item IToDoItem)
	Delete(id int)IToDoItem
	Update(id int,newItem IToDoItem)IToDoItem
	GetAll() []IToDoItem
}

type IToDoItem interface {
	GetItemId() int
	GetTitle() string
	GetDescription() string
}