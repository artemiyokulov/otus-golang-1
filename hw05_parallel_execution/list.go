package hw05parallelexecution

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len         int
	front, back *ListItem
}

func (instance list) Len() int {
	return instance.len
}

func (instance list) Front() *ListItem {
	return instance.front
}

func (instance list) Back() *ListItem {
	return instance.back
}

func (instance *list) pasteToFront(i *ListItem) {
	if instance.len == 0 {
		instance.front = i
		instance.back = i
	} else {
		i.Next = instance.front
		instance.front.Prev = i
		instance.front = i
	}
	instance.len++
}

func (instance *list) PushFront(v interface{}) *ListItem {
	var newItem ListItem
	newItem.Value = v

	instance.pasteToFront(&newItem)

	return instance.front
}

func (instance *list) PushBack(v interface{}) *ListItem {
	var newItem ListItem
	newItem.Value = v

	if instance.len == 0 {
		instance.front = &newItem
		instance.back = &newItem
	} else {
		newItem.Prev = instance.back
		instance.back.Next = &newItem
		instance.back = &newItem
	}

	instance.len++
	return instance.back
}

func (instance *list) exclude(i *ListItem) {
	// если исключаемый элемент front
	//nolint:gocritic
	if i.Prev == nil {
		instance.front = i.Next
		instance.front.Prev = nil
	} else if i.Next == nil { // если исключаемый элемент back
		instance.back = i.Prev
		instance.back.Next = nil
	} else {
		prev := i.Prev
		next := i.Next

		next.Prev = prev
		prev.Next = next
	}
	instance.len--
}

func (instance *list) Remove(i *ListItem) {
	if instance.len <= 1 {
		instance.len = 0
		instance.front = nil
		instance.back = nil
		return
	}

	instance.exclude(i)
}

func (instance *list) MoveToFront(i *ListItem) {
	if instance.len <= 1 {
		return
	}

	instance.exclude(i)
	instance.pasteToFront(i)
}

func NewList() List {
	return new(list)
}
