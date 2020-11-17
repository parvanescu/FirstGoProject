package repository

import (
	"ExGabi/model"
	"fmt"
)

type IRepository interface{
	Add(item model.ToDoItem)
	Delete(id int) model.ToDoItem
	Update(id int,newItem model.ToDoItem) model.ToDoItem
	GetAll() []model.ToDoItem
}

type repository struct{
	dataArray []model.ToDoItem
}

func New()IRepository{
	fmt.Println("Repo initialized")
	repo :=repository{make([]model.ToDoItem,0,1)}
	return &repo
}


func (r *repository)Add(item model.ToDoItem){
	r.dataArray = append(r.dataArray,item)
}

func (r *repository)Delete(id int)model.ToDoItem{
	for i,v := range r.dataArray{
		if v.ItemId == id{
			item :=v
			r.dataArray = append(r.dataArray[0:i],r.dataArray[i+1:cap(r.dataArray)]...)
			return item
		}
	}
	return model.ToDoItem{
		ItemId:      0,
		Title:       "",
		Description: "",
	}
}

func (r *repository)Update(id int,newItem model.ToDoItem)model.ToDoItem{
	for i,v  := range r.dataArray{
		if v.ItemId == id{
			newArray := append(r.dataArray[0:i],newItem)
			item := r.dataArray[i]
			r.dataArray = append(newArray,r.dataArray[i+1:cap(r.dataArray)]...)
			return item
			// Ce metoda ?????
			//r.dataArray[i].Description = newItem.Description
			//r.dataArray[i].Title = newItem.Title
			//r.dataArray[i].ItemId=newItem.ItemId
		}
	}
	return model.ToDoItem{
		ItemId:      0,
		Title:       "",
		Description: "",
	}
}

func (r *repository)GetAll() []model.ToDoItem{
	return r.dataArray
}
