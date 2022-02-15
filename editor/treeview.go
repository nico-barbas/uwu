package editor

import (
	"github.com/nico-ec/uwu/ui"
)

const (
	treeviewWidth = 140
)

type treeview struct {
	list *ui.List
}

func newTreeview(edParent ui.Handle, sepImg *Image, font *Font) treeview {
	treeview := treeview{}
	treeview.list = &ui.List{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{232, 152, 168, 255},
			Img:     sepImg,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Style: ui.Style{
			Padding: 3,
			Margin:  ui.Point{5, 0},
		},

		Name:       "Root",
		Font:       font,
		TextSize:   12,
		TextClr:    uwuTextClr,
		IndentSize: 10,
	}
	ui.AddWidget(edParent, treeview.list, treeviewWidth)

	return treeview
}

func (t *treeview) loadProject(p *project) {
	t.list.Receiver = t
	for k, v := range p.root.nodes {
		switch f := v.(type) {
		case *folder:
			subList := ui.NewSubList(k)
			t.list.AddItem(&subList)
			populateSubList(&subList, f)
		case file:
			t.list.AddItem(&ui.ListItem{ItemName: k})
		}
	}
	t.list.SortList()
}

func (t *treeview) OnItemSelected(item ui.ListNode) {
	openProjectFile(item.Name())
}

func populateSubList(l *ui.SubList, f *folder) {
	for k, v := range f.nodes {
		switch f := v.(type) {
		case *folder:
			subList := ui.NewSubList(k)
			populateSubList(&subList, f)
			l.AddItem(&subList, 10, 12)
		case file:
			l.AddItem(&ui.ListItem{ItemName: k}, 10, 12)
		}
	}
}
