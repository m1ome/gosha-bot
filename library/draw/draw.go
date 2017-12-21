package draw

import (
	"fmt"
	"strings"
	"errors"

	"github.com/fogleman/gg"
)

const (
	StringSize = 15
	StringHeight = 28
)

var (
	ErrTooLongWord = errors.New("Too long word in sentence")
)

// Defaults
// "/Library/Fonts/Arial.ttf"
func Text(img, font, text, output string) error{
	im, err := gg.LoadImage(img)
	if err != nil {
		return fmt.Errorf("Error reading image: %v", err)
	}

	bounds := im.Bounds()
	size := bounds.Size()

	image := gg.NewContext(size.X, size.Y)
	image.DrawImage(im, 0, 0)

	txt := gg.NewContext(300, 300)
	txt.Rotate(gg.Radians(6))
	txt.Clear()
	txt.SetRGB(0,0,0,)
	if err := txt.LoadFontFace(font, 28); err != nil {
		return fmt.Errorf("Error looking font: %v", err)
	}

	words := strings.Split(text, " ")
	var tmp []string
	var symbols int
	var idx int

	for _, word := range words{
		if symbols + len(word) > StringSize {
			// Adding text
			txt.DrawString(strings.Join(tmp, " "), 50, 50 + float64(StringHeight*idx))

			// Making new array of words
			tmp = []string{word}

			// Resting counters
			symbols = 0
			idx++
		} else {
			tmp = append(tmp, word)
			symbols += len(word) + 1
		}
	}
	txt.DrawString(strings.Join(tmp, " "), 50, 50+float64(StringHeight*idx))

	txt.Clip()

	image.DrawImage(txt.Image(), 240, 515)
	image.SavePNG(output)

	return nil
}