package editor

import "github.com/nico-ec/uwu/ui"

type treeview struct {
}

func populateSubList(l *ui.SubList, f *folder) {
	for k, v := range f.nodes {
		switch f := v.(type) {
		case *folder:
			subList := ui.NewSubList(k)
			populateSubList(&subList, f)
			l.AddItem(&subList)
		case file:
			l.AddItem(&ui.ListItem{Name: k})
		}
	}
}
