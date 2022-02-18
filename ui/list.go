package ui

import (
	"log"
	"sort"
)

const (
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
		Receiver   ListReceiver

		cursorVisible bool
		cursorRect    Rectangle
		selectedNode  ListNode
	}

	SubList struct {
		ItemName   string
		Collapsed  bool
		nameHeight float64
		items      []ListNode
		count      int
		origin     Point
	}

	ListItem struct {
		ItemName string
		origin   Point
		height   float64
	}

	ListReceiver interface {
		OnItemSelected(item ListNode)
	}
)

type ListNode interface {
	Name() string
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
		l.selectedNode = l.Root.selectNode(mPos)
		if l.selectedNode != nil {
			l.cursorVisible = true
			l.cursorRect.Y = l.selectedNode.getOrigin()[1]
			if isMouseJustPressed() {
				switch s := l.selectedNode.(type) {
				case *SubList:
					s.Collapsed = !s.Collapsed
					l.Root.orderItems(l.TextSize)
				case *ListItem:
					if l.Receiver != nil {
						l.Receiver.OnItemSelected(l.selectedNode)
					} else {
						log.SetPrefix("[UI Debug]: ")
						log.Println("No receiver attached to this list")
					}
				}
			}
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
	l.Root.AddItem(i, l.IndentSize, l.TextSize)
}

func (l *List) SortList() {
	l.Root.sort(l.TextSize)
}

func NewSubList(name string) SubList {
	return SubList{
		ItemName:  name,
		Collapsed: false,
		items:     make([]ListNode, subListInitialCap),
	}
}

func (s *SubList) AddItem(i ListNode, indentSize float64, lineSize float64) {
	name := i.Name()
	var exist bool
	for i := 0; i < s.count; i += 1 {
		item := s.items[i]
		if item.Name() == name {
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
		log.Printf("List %s already has a child with name %s", s.Name(), name)
	}
}

func (s *SubList) Name() string {
	return s.ItemName
}

func (s *SubList) sort(lineSize float64) {
	sortFolderFn := func(i, j int) bool {
		_, iFolder := s.items[i].(*SubList)
		_, jFolder := s.items[j].(*SubList)
		return iFolder && !jFolder
	}
	sortFn := func(i, j int) bool {
		return s.items[i].Name() < s.items[j].Name()
	}
	sort.SliceStable(s.items[:s.count], sortFolderFn)
	maxFolderIndex := 0
	for i := 0; i < s.count; i += 1 {
		_, isFolder := s.items[i].(*SubList)
		if !isFolder {
			break
		}
		maxFolderIndex += 1
	}

	sort.SliceStable(s.items[:maxFolderIndex], sortFn)
	sort.SliceStable(s.items[maxFolderIndex:s.count], sortFn)
	s.orderItems(lineSize)
}

func (s *SubList) orderItems(lineSize float64) {
	if s.Collapsed {
		return
	}
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
		Text: s.ItemName,
	})
	yPtr := size
	if !s.Collapsed {
		for i := 0; i < s.count; i += 1 {
			item := s.items[i]
			h := item.draw(buf, f, size, clr)
			yPtr += h
		}
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X:      s.origin[0],
				Y:      s.origin[1] + size,
				Width:  1,
				Height: yPtr - size,
			},
			Clr: clr,
		})
	}
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
	if !s.Collapsed {
		for i := 0; i < s.count; i += 1 {
			height += s.items[i].getHeight()
		}
	}
	return height
}

func (s *SubList) setHeight(h float64) {
	s.nameHeight = h
}

func (l *ListItem) Name() string {
	return l.ItemName
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
		Text: l.ItemName,
	})
	return size
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
