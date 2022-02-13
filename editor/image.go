package editor

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	data *ebiten.Image
}

func (i *Image) GetWidth() float64 {
	return float64(i.data.Bounds().Dx())
}

func (i *Image) GetHeight() float64 {
	return float64(i.data.Bounds().Dy())
}
