package ui

import (
	"log"
	"sort"
)

const (
	listLineSpacing   = 0
	subListInitialCap = 10
)

type (
	List struct {
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

		cursorVisible bool
		cursorRect    Rectangle
		selectedNode  ListNode
	}

	SubList struct {
		Name       string
		nameHeight float64
		items      []ListNode
		count      int
		origin     Point
	}

	ListItem struct {
		Name   string
		origin Point
		height float64
	}
)

type ListNode interface {
	name() string
	draw(buf *renderBuffer, f Font, size float64, clr Color) float64
	getOrigin() Point
	setOrigin(p Point)
	getHeight() float64
	setHeight(h float64)
}

func (l *List) init() {
	l.activeRect = Rectangle{
		X:      l.rect.X + l.Style.Margin[0],
		Y:      l.rect.Y + l.Style.Margin[1],
		Width:  l.rect.Width - l.Style.Margin[0]*2,
		Height: l.rect.Height - l.Style.Margin[1]*2,
	}
	l.Root = NewSubList("Root")
	l.Root.origin = Point{l.activeRect.X, l.activeRect.Y}
	// Lots of assumptions made here
	// Isn't really flexible
	l.cursorRect = Rectangle{
		X:      l.activeRect.X,
		Width:  l.activeRect.Width,
		Height: l.TextSize,
	}
	l.Root = NewSubList(l.Name)
	l.Root.origin = Point{
		l.activeRect.X,
		l.activeRect.Y,
	}
}

func (l *List) update() {
	mPos := mousePosition()
	if l.activeRect.pointInBounds(mPos) {
		l.cursorVisible = true
		l.selectedNode = l.Root.selectNode(mPos)
		if l.selectedNode != nil {
			l.cursorRect.Y = l.selectedNode.getOrigin()[1]
		} else {
			l.cursorVisible = false
		}
	} else {
		l.cursorVisible = false
	}
}

func (l *List) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)
	if l.cursorVisible {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: l.cursorRect,
			Clr:  Color{l.TextClr[0], l.TextClr[1], l.TextClr[2], 155},
		})
	}
	l.Root.draw(buf, l.Font, l.TextSize, l.TextClr)
}

func (l *List) AddItem(i ListNode) {
	l.Root.AddItem(i, l.IndentSize, l.TextSize+listLineSpacing)
}

func (l *List) SortList() {
	l.Root.sort(l.TextSize + listLineSpacing)
}

func NewSubList(name string) SubList {
	return SubList{
		Name:  name,
		items: make([]ListNode, subListInitialCap),
	}
}

func (s *SubList) AddItem(i ListNode, indentSize float64, lineSize float64) {
	name := i.name()
	var exist bool
	for i := 0; i < s.count; i += 1 {
		item := s.items[i]
		if item.name() == name {
			exist = true
			break
		}
	}
	if !exist {
		if s.count >= len(s.items) {
			newBuf := make([]ListNode, len(s.items)*2)
			copy(newBuf[:], s.items[:])
			s.items = newBuf
		}
		s.items[s.count] = i
		s.count += 1
		i.setOrigin(Point{
			s.origin[0] + indentSize,
			0,
		})
		i.setHeight(lineSize)
	} else {
		log.Printf("List %s already has a child with name %s", s.Name, name)
	}
}

func (s *SubList) name() string {
	return s.Name
}

func (s *SubList) sort(lineSize float64) {
	sortFn := func(i, j int) bool {
		return s.items[i].name() < s.items[j].name()
	}
	sort.SliceStable(s.items[:s.count], sortFn)
	yPtr := s.origin[1] + lineSize
	for i := 0; i < s.count; i += 1 {
		item := s.items[i]
		iPos := item.getOrigin()
		item.setOrigin(Point{
			iPos[0],
			yPtr,
		})
		yPtr += item.getHeight()
		switch s := item.(type) {
		case *SubList:
			s.sort(lineSize)
		}
	}
}

func (s *SubList) draw(buf *renderBuffer, f Font, size float64, clr Color) float64 {
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: Rectangle{
			X:      s.origin[0],
			Y:      s.origin[1],
			Height: size,
		},
		Clr:  clr,
		Font: f,
		Text: s.Name,
	})
	yPtr := size + listLineSpacing
	for i := 0; i < s.count; i += 1 {
		item := s.items[i]
		h := item.draw(buf, f, size, clr)
		yPtr += h
	}
	buf.addEntry(RenderEntry{
		Kind: RenderRectangle,
		Rect: Rectangle{
			X:      s.origin[0],
			Y:      s.origin[1] + size + listLineSpacing,
			Width:  1,
			Height: yPtr - size + listLineSpacing,
		},
		Clr: clr,
	})
	return yPtr
}

func (s *SubList) selectNode(at Point) ListNode {
	if at[1] >= s.origin[1] && at[1] < s.origin[1]+s.nameHeight {
		return s
	}
	var selected ListNode
	for i := 0; i < s.count; i += 1 {
		item := s.items[i]
		switch i := item.(type) {
		case *SubList:
			selected = i.selectNode(at)
		case *ListItem:
			if at[1] >= i.origin[1] && at[1] < i.origin[1]+i.height {
				selected = i
			}
		}
		if selected != nil {
			break
		}
	}
	return selected
}

func (s *SubList) getOrigin() Point {
	return s.origin
}

func (s *SubList) setOrigin(p Point) {
	s.origin = p
}

func (s *SubList) getHeight() float64 {
	height := s.nameHeight
	for i := 0; i < s.count; i += 1 {
		height += s.items[i].getHeight()
	}
	return height
}

func (s *SubList) setHeight(h float64) {
	s.nameHeight = h
}

func (l *ListItem) name() string {
	return l.Name
}

func (l *ListItem) draw(buf *renderBuffer, f Font, size float64, clr Color) float64 {
	buf.addEntry(RenderEntry{
		Kind: RenderText,
		Rect: Rectangle{
			X:      l.origin[0],
			Y:      l.origin[1],
			Height: size,
		},
		Clr:  clr,
		Font: f,
		Text: l.Name,
	})
	return size + listLineSpacing
}

func (l *ListItem) getOrigin() Point {
	return l.origin
}

func (l *ListItem) setOrigin(p Point) {
	l.origin = p
}

func (l *ListItem) getHeight() float64 {
	return l.height
}

func (l *ListItem) setHeight(h float64) {
	l.height = h
}
