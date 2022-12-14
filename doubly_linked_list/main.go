package main

import (
	"fmt"
	"strconv"
)

type entry struct {
	data string
	prev *entry
	next *entry
}

type list struct {
	entries int
	first   *entry
	last    *entry
}

func CreateList() *list {
	list := &list{
		entries: 0,
		first:   nil,
		last:    nil,
	}
	return list
}

func Len(x *list) int {
	return x.entries
}

func First(x *list) *entry {
	return x.first
}

func Last(x *list) *entry {
	return x.last
}

func (y *list) ItemPushFront(x string) {
	newitem := &entry{
		data: x,
		prev: nil,
	}
	if y.entries == 0 {
		newitem.next = nil
		y.last = newitem
	} else {
		y.first.prev = newitem
		newitem.next = y.first
	}
	y.first = newitem
	y.entries++
	return
}

func (y *list) ItemPushEnd(x string) {
	newitem := &entry{
		data: x,
		next: nil,
	}
	if y.entries == 0 {
		newitem.prev = nil
		y.first = newitem
	} else {
		y.last.next = newitem
		newitem.prev = y.last
	}
	y.last = newitem
	y.entries++
}

func (y *list) Range() {
	fmt.Println("Linked list")
	fmt.Println("Number of Entries:", y.entries)
	fmt.Println("First entry pointer:", y.first)
	fmt.Println("Last entry pointer:", y.last)
	current := y.first
	for current != nil {
		fmt.Println("=================")
		fmt.Println("Previous entry:", current.prev)
		fmt.Println("Data:", current.data)
		fmt.Println("Next entry:", current.next)
		current = current.next
	}
	fmt.Println("=================")
}

func (y *list) RemoveEntry(entry *entry) {
	if y.entries <= 0 {
		fmt.Println("Empty list!")
		return
	}
	current := y.first
	for current != nil {
		if current != entry {
			current = current.next
		} else {
			prev := current.prev
			next := current.next
			if current.next == nil {
				prev.next = nil
			}
			if current.prev == nil {
				next.prev = nil
			}
			if current.prev != nil && current.next != nil {
				next.prev = prev
				prev.next = next
			}
			current.Clear()
			y.entries -= 1
			current = current.next
		}
	}
	if entry == y.first {
		y.first = entry.next
	}
	if entry == y.last {
		y.last = entry.prev
	}
	return
}

func (y *entry) Clear() {
	y = &entry{}
	return
}

func (y *list) NextEntry(entry *entry) *entry {
	if y.entries <= 0 {
		fmt.Println("Empty list!")
		return nil
	}
	current := y.first
	for current != nil {
		if current != entry {
			current = current.next
		} else {
			return current.next
		}
	}
	return nil
}

func (y *list) PrevEntry(entry *entry) *entry {
	if y.entries <= 0 {
		fmt.Println("Empty list!")
		return nil
	}
	current := y.first
	for current != nil {
		if current != entry {
			current = current.next
		} else {
			return current.prev
		}
	}
	return nil
}

func (y *list) GetItems(x string) []*entry {
	var g []*entry
	current := y.first
	for current != nil {
		if current.data == x {
			g = append(g, current)
		}
		current = current.next
	}
	return g
}

func main() {
	fmt.Println("Creating list")
	y := CreateList()
	fmt.Println("Len x function, number of entries:", Len(y))
	fmt.Println("First entry function", First(y))
	fmt.Println("Last entry function", Last(y))
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Adding 5 items to end \n")
	for i := 1; i < 5; i++ {
		g := strconv.Itoa(i)
		y.ItemPushEnd(g)
	}
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Adding 5 items to front \n")
	for i := 6; i < 11; i++ {
		g := strconv.Itoa(i)
		y.ItemPushFront(g)
	}
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Adding item test \n")
	y.ItemPushEnd("test")
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Get items:test:", y.GetItems("test"))
	t := y.GetItems("test")
	fmt.Println("Next item:test")
	for _, v := range t {
		fmt.Println(y.NextEntry(v))
	}
	fmt.Println("Previous item:test")
	for _, v := range t {
		fmt.Println(y.PrevEntry(v))
	}
	fmt.Println("Removing items:test")
	for _, v := range t {
		y.RemoveEntry(v)
	}
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Removing items:3")
	b := y.GetItems("3")
	for _, v := range b {
		y.RemoveEntry(v)
	}
	fmt.Println("Ranging list")
	y.Range()
	fmt.Println("Len x function", Len(y))
}
