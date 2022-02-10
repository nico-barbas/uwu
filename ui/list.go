package ui

import (
	"fmt"
	"log"
)

const (
	listLineSpacing   = 0
	subListInitialCap = 10
)

type List struct {
	widgetRoot

	Background Background

	Name       string
	Font       Font
	TextSize   float64
	TextClr    Color
	root       SubList
	IndentSize float64
}

func (l *List) init() {
	l.root = newSubList(l.Name)
}

func (l *List) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)
	rootRect := Rectangle{
		X:      l.rect.X,
		Y:      l.rect.Y,
		Height: l.TextSize,
	}
	l.root.draw(buf, l.Font, rootRect, l.TextClr, l.IndentSize)
}

func (l *List) AddItem(i ListNode) {
	l.root.addItem(i)
}

type ListNode interface {
	name() string
	draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64)
}

type (
	SubList struct {
		Name  string
		items []ListNode
	}

	ListItem struct {
		Name string
	}
)

func newSubList(name string) SubList {
	return SubList{
		Name:  name,
		items: make([]ListNode, 0, subListInitialCap),
	}
}

func (s *SubList) addItem(i ListNode) {
	name := i.name()
	var exist bool
	for _, item := range s.items {
		if item.name() == name {
			exist = true
			break
		}
	}
	if !exist {
		s.items = append(s.items, i)
	} else {
		log.Printf("List %s already has a child with name %s", s.Name, name)
	}
}

func (s *SubList) name() string {
	return s.Name
}

func (s *SubList) draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64) {
	fmt.Println("Start List")
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: r,
		Clr:  clr,
		Font: f,
		Text: s.Name,
	})
	i := 1
	for _, item := range s.items {
		childRect := Rectangle{
			X:      r.X + indent,
			Y:      r.Y + (r.Height+listLineSpacing)*float64(i),
			Height: r.Height,
		}
		item.draw(buf, f, childRect, clr, indent)
		i += 1
	}
}

func (l *ListItem) name() string {
	return l.Name
}

func (l *ListItem) draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64) {
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: r,
		Clr:  clr,
		Font: f,
		Text: l.Name,
	})
}
