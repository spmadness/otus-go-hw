package hw04lrucache

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
	First  *ListItem
	Last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	if l.First == nil {
		l.First = li
		l.Last = li
	} else {
		li.Next = l.First
		li.Next.Prev = li
		l.First = li
	}
	l.length++
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	if l.Last == nil {
		l.PushFront(v)
	} else {
		li.Prev = l.Last
		li.Prev.Next = li
		l.Last = li
		l.length++
	}
	return li
}

func (l *list) Remove(i *ListItem) {
	if i != nil {
		if i.Prev == nil {
			l.First = i.Next
		} else {
			i.Prev.Next = i.Next
		}
		if i.Next == nil {
			l.Last = i.Prev
		} else {
			i.Next.Prev = i.Prev
		}
		i = nil
		l.length--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i != nil {
		l.Remove(i)
		l.PushFront(i.Value)
	}
}

func NewList() List {
	return new(list)
}
