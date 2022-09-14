package hikaku

import (
	"image"
	"image/color"
)

func Grayscale(img image.Image) (result *image.Gray) {
	bounds := img.Bounds()

	result = image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			pixel = color.GrayModel.Convert(pixel)
			result.Set(x, y, pixel)
		}
	}

	return
}
