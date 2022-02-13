package editor

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/nico-ec/uwu/ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Font struct {
	faces map[int]font.Face
}

func NewFont(path string, dpi float64, sizes []int) Font {
	f := Font{
		faces: make(map[int]font.Face, len(sizes)),
	}

	fontData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}

	for _, v := range sizes {
		face, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(v),
			DPI:     dpi,
			Hinting: font.HintingNone,
		})
		if err != nil {
			panic(err)
		}
		f.faces[v] = face
	}
	return f
}

func (f *Font) GlyphAdvance(r rune, size float64) float64 {
	x, _ := f.faces[int(size)].GlyphAdvance(r)
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func (f *Font) MeasureText(t string, size float64) ui.Point {
	measure := ui.Point{}

	if v, exist := f.faces[int(size)]; !exist {
		panic("No face of size in given Font")
	} else {
		r := text.BoundString(v, t)
		measure[0] = float64(r.Dx())
		measure[1] = float64(r.Dy())
	}

	return measure
}
