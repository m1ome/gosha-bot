package draw

import (
	"testing"
	"os"
)

func TestText(t *testing.T) {
	var tests = []struct{
		Text string
		Filename string
	}{
		{
			"Hello",
			"basic.png",
		},
		{
			`This is pretty long text and should be definetely scaled properly`,
			"long.png",
		},
		{
			"Тест русского языка с переносами и всякими ништяками",
			"long_rus.png",
		},
	}

	// Running all cases
	for _, test := range tests {
		os.Remove("tmp/" + test.Filename)
		Text("../../assets/image.png", "/Library/Fonts/Arial.ttf", test.Text, "tmp/" + test.Filename)
	}
}
