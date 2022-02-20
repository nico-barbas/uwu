package editor

import "github.com/nico-ec/uwu/ui"

var (
	lightTheme = theme{
		backgroundClr1: ui.Color{247, 231, 230, 255},
		backgroundClr2: ui.Color{247, 231, 230, 255},
		backgroundClr3: ui.Color{255, 95, 131, 255},
		dividerClr:     ui.Color{255, 95, 131, 255},
		normalTextClr:  ui.Color{255, 95, 131, 255},
		normalTextClr2: ui.Color{247, 231, 230, 255},

		syntaxKeywordClr: ui.Color{200, 106, 255, 255},
		syntaxNormalClr:  ui.Color{255, 95, 131, 255},
		syntaxNumberClr:  ui.Color{213, 133, 128, 255},
	}
)

type theme struct {
	backgroundClr1 ui.Color // The textbox area
	backgroundClr2 ui.Color // The treeview, window header and tabview header
	backgroundClr3 ui.Color // The statusbar and the selected tab
	dividerClr     ui.Color
	rulerClr       ui.Color
	normalTextClr  ui.Color
	normalTextClr2 ui.Color

	syntaxNormalClr   ui.Color
	syntaxSymbolClr   ui.Color
	syntaxCommentClr  ui.Color
	syntaxKeywordClr  ui.Color
	syntaxFunctionClr ui.Color
	syntaxNumberClr   ui.Color
	syntaxStringClr   ui.Color
}
