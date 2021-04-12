package hw04lrucache


type List interface {
	Len() int // длина списка
	Front() *ListItem // первый элемент списка
	Back() *ListItem // последний элемент списка
	PushFront(v interface{}) *ListItem  // добавить значение в начало
	PushBack(v interface{}) *ListItem // добавить значение в конец
	Remove(i *ListItem) // удалить элемент
	MoveToFront(i *ListItem) // переместить элемент в начало
}

type ListItem struct {
	Value interface{}  // значение
	Next  *ListItem  // следующий элемент
	Prev  *ListItem // предыдущий элемент
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}
func (l *list) Len() int  {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.front,
	}
	return l.pushFront(item)
}
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.Len() != 0 {
		l.back.Next = item
	} else {
		l.front = item
	}
	l.back = item
	l.length++
	return l.back
}
func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	if l.front == item {
		l.front = item.Next
	}
	if l.back == item {
		l.back = item.Prev
	}
	l.length--
}
func (l *list) MoveToFront(item *ListItem) {
	if l.front != item {
		l.Remove(item)
		l.pushFront(item)
	}
}


func (l *list) pushFront(item *ListItem) *ListItem {
	item.Prev = nil
	item.Next = l.front
	if l.Len() != 0 {
		l.front.Prev = item
	} else {
		l.back = item
	}
	l.front = item
	l.length++
	return l.front
}