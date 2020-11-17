package useCase

import (
	"ExGabi/model"
	"ExGabi/repository"
	"fmt"
)

type IUseCase interface {
	AddItem(id int,title string,description string)
	DeleteItem(id int)model.ToDoItem
	UpdateItem(id int,newTitle string,newDescription string)model.ToDoItem
	GetItem(id int)model.ToDoItem
	GetAll()//[]model.ToDoItem
}

type UseCase struct{
	repository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)Init(repository repository.IRepository){
	uC.repository = repository
}

func (uC *UseCase)AddItem(id int,title string,description string){
	var item model.ToDoItem = model.ToDoItem{
		ItemId:      id,
		Title:       title,
		Description: description,
	}
	uC.repository.Add(item)
}

func (uC *UseCase)DeleteItem(id int)model.ToDoItem{
	return model.ToDoItem{
		ItemId:      0,
		Title:       "",
		Description: "",
	}
}

func (uC *UseCase)UpdateItem(id int,newTitle string,newDescription string)model.ToDoItem{
	return model.ToDoItem{
		ItemId:      0,
		Title:       "",
		Description: "",
	}
}

func(uC *UseCase)GetItem(id int)model.ToDoItem{
	return model.ToDoItem{
		ItemId:      0,
		Title:       "",
		Description: "",
	}
}

func(uC *UseCase)GetAll(){
	items := uC.repository.GetAll()
	for _,v := range items{
		fmt.Println(v)
	}
}
