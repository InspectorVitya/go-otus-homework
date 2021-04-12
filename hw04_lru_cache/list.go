package hw04lrucache

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый элемент списка
	Back() *listItem                   // последний элемент списка
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(item *listItem)             // удалить элемент
	MoveToFront(item *listItem)        // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	length int
	front  *listItem
	tail   *listItem
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *listItem {
	item := &listItem{
		Value: v,
		Next:  l.front,
	}
	return l.pushFront(item)
}

func (l *list) PushBack(v interface{}) *listItem {
	item := &listItem{
		Value: v,
		Prev:  l.tail,
	}
	if l.Len() != 0 {
		l.tail.Next = item
	} else {
		l.front = item
	}
	l.tail = item
	l.length++
	return l.tail
}

func (l *list) Remove(item *listItem) {
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
	if l.tail == item {
		l.tail = item.Prev
	}
	l.length--
}

func (l *list) MoveToFront(item *listItem) {
	if l.front != item {
		l.Remove(item)
		l.pushFront(item)
	}
}

func (l *list) pushFront(item *listItem) *listItem {
	item.Prev = nil
	item.Next = l.front
	if l.Len() != 0 {
		l.front.Prev = item
	} else {
		l.tail = item
	}
	l.front = item
	l.length++
	return l.front
}
