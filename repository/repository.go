package repository

import (
	"ExGabi/model"
	"fmt"
)

type IRepository interface{
	Add(item model.IToDoItem)
	Delete(id int) model.IToDoItem
	Update(id int,newItem model.IToDoItem) model.IToDoItem
	GetAll() []model.IToDoItem
}

type repository struct{
	dataArray []model.IToDoItem
}

func New()IRepository{
	fmt.Println("Repo initialized")
	repo :=repository{make([]model.IToDoItem,0,1)}
	return &repo
}


func (r *repository)Add(item model.IToDoItem){
	r.dataArray = append(r.dataArray,item)
}

func (r *repository)Delete(id int)model.IToDoItem{
	for i,v := range r.dataArray{
		if v.GetItemId() == id{
			item :=v
			r.dataArray = append(r.dataArray[0:i],r.dataArray[i+1:cap(r.dataArray)]...)
			return item
		}
	}
	return model.New(0,"","")
}

func (r *repository)Update(id int,newItem model.IToDoItem)model.IToDoItem{
	for i,v  := range r.dataArray{
		if v.GetItemId() == id{
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
	return model.New(0,"","")
}

func (r *repository)GetAll() []model.IToDoItem{
	return r.dataArray
}
