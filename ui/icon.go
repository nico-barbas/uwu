package ui

type Icon struct {
	widgetRoot

	Img    Image
	ImgClr Color
}

func (i *Icon) draw(buf *renderBuffer) {
	if i.Img == nil {
		return
	}
	iconRect := Rectangle{
		X: i.rect.X + (i.rect.Width/2 - i.Img.GetWidth()/2),
		Y: i.rect.Y + (i.rect.Height/2 - i.Img.GetHeight()/2),
	}
	buf.addEntry(RenderEntry{
		Kind: RenderImage,
		Rect: iconRect,
		Img:  i.Img,
		Clr:  i.ImgClr,
	})
}
