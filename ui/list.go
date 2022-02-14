package ui

import (
	"log"
	"sort"
)

const (
	listLineSpacing   = 0
	subListInitialCap = 10
)

type List struct {
	widgetRoot

	Background Background
	Style      Style

	activeRect Rectangle
	Name       string
	Font       Font
	TextSize   float64
	TextClr    Color
	Root       SubList
	IndentSize float64
}

func (l *List) init() {
	l.activeRect = Rectangle{
		X:      l.rect.X + l.Style.Margin[0],
		Y:      l.rect.Y + l.Style.Margin[1],
		Width:  l.rect.Width - l.Style.Margin[0]*2,
		Height: l.rect.Height - l.Style.Margin[1]*2,
	}
	l.Root = NewSubList(l.Name)
}

func (l *List) update() {
	mPos := mousePosition()
	if l.activeRect.pointInBounds(mPos) {

	}
}

func (l *List) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)
	RootRect := Rectangle{
		X:      l.activeRect.X,
		Y:      l.activeRect.Y,
		Height: l.TextSize,
	}
	l.Root.draw(buf, l.Font, RootRect, l.TextClr, l.IndentSize)
}

func (l *List) AddItem(i ListNode) {
	l.Root.AddItem(i)
}

func (l *List) SortList() {
	l.Root.sort()
}

type ListNode interface {
	name() string
	draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64) float64
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

func NewSubList(name string) SubList {
	return SubList{
		Name:  name,
		items: make([]ListNode, 0, subListInitialCap),
	}
}

func (s *SubList) AddItem(i ListNode) {
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

func (s *SubList) sort() {
	sortFn := func(i, j int) bool {
		return s.items[i].name() < s.items[j].name()
	}
	sort.SliceStable(s.items, sortFn)
	for _, item := range s.items {
		switch s := item.(type) {
		case *SubList:
			s.sort()
		}
	}
}

func (s *SubList) draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64) float64 {
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: r,
		Clr:  clr,
		Font: f,
		Text: s.Name,
	})
	yPtr := r.Height + listLineSpacing
	for _, item := range s.items {
		childRect := Rectangle{
			X:      r.X + indent,
			Y:      r.Y + yPtr,
			Height: r.Height,
		}
		h := item.draw(buf, f, childRect, clr, indent)
		yPtr += h
	}
	buf.addEntry(RenderEntry{
		Kind: RenderRectangle,
		Rect: Rectangle{
			X:      r.X,
			Y:      r.Y + r.Height + listLineSpacing,
			Width:  1,
			Height: yPtr - r.Height + listLineSpacing,
		},
		Clr: clr,
	})
	return yPtr
}

func (l *ListItem) name() string {
	return l.Name
}

func (l *ListItem) draw(buf *renderBuffer, f Font, r Rectangle, clr Color, indent float64) float64 {
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: r,
		Clr:  clr,
		Font: f,
		Text: l.Name,
	})
	return r.Height + listLineSpacing
}
