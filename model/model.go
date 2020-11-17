package model
import "ExGabi/interfaces"
//type IToDoItem interface {
//	GetItemId() int
//	GetTitle() string
//	GetDescription() string
//}

type toDoItem struct{
	ItemId      int
	Title       string
	Description string
}

func New(itemId int,title string,description string) interfaces.IToDoItem {
	return toDoItem{
		ItemId:      itemId,
		Title:       title,
		Description: description,
	}
}

func (item toDoItem)GetItemId()int{
	return item.ItemId
}

func (item toDoItem)GetTitle()string{
	return item.Title
}

func (item toDoItem)GetDescription()string{
	return item.Description
}
