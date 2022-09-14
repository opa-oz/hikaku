package hikaku

import (
	"image"
	"image/color"
	"image/draw"
)

func drawCanvas(width int, height int, backgroundColor color.Color) *image.RGBA {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	return img
}

type Mask struct {
	bounds image.Rectangle
	width  int
	height int
}

func (c Mask) Draw(targetImage image.Image, img image.Image, transparency uint8) *image.RGBA {
	maskImage := drawCanvas(c.width, c.height, color.RGBA{R: transparency, G: transparency, B: transparency, A: 0xff})

	bounds := c.bounds
	mask := image.NewAlpha(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, _, _, _ := maskImage.At(x, y).RGBA()
			mask.SetAlpha(x, y, color.Alpha{A: uint8(255 - r)})
		}
	}

	m := image.NewRGBA(bounds)
	draw.Draw(m, m.Bounds(), targetImage, image.Point{}, draw.Src)
	draw.DrawMask(m, bounds, img, image.Point{}, mask, image.Point{}, draw.Over)

	return m
}
