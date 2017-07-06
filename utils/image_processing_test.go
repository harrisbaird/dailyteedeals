package utils_test

import (
	"image/color"
	"testing"

	. "github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestTransformImage(t *testing.T) {
	defaultTransform := Transform{Size: 100}

	testCases := []struct {
		name      string
		input     string
		output    string
		transform Transform
	}{
		{"Horizontal", "1.png", "1.png", defaultTransform},
		{"Vertical", "2.png", "2.png", defaultTransform},
		{"JPEG", "3.jpg", "3.png", defaultTransform},
		{"Transparent PNG", "4.png", "4.png", defaultTransform},
		{"Transparent GIF", "5.gif", "5.png", defaultTransform},
		{"Non-animated GIF", "6.gif", "6.png", defaultTransform},
		{"Animated GIF", "7.gif", "7.png", defaultTransform},
		{"No size", "1.png", "8.png", Transform{}},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			input := LoadImageFixture("transform/input/" + tt.input)
			output := TransformImage(input, tt.transform, color.Black)
			fixture := LoadImageFixture("transform/output/" + tt.output)
			st.Expect(t, output, fixture)
		})
	}
}

func TestFindBackgroundColor(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"1", "1.png", "#243cc4"},
		{"2", "2.png", "#131522"},
		{"3", "3.png", "#000000"},
		{"4", "4.png", "#dcdddf"},
		{"5", "5.png", "#a19ba5"},
		{"6", "6.png", "#000000"},
		{"7", "7.png", "#ffffff"},
		{"8", "8.gif", "#65d8ff"},
		{"9", "9.gif", "#47b6e0"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			img := LoadImageFixture("background/" + tt.input)
			color := FindBackgroundColor(img)
			st.Expect(t, ColorToHex(color), tt.want)
		})
	}
}

func TestColorToHex(t *testing.T) {
	testCases := []struct {
		name  string
		input color.Color
		want  string
	}{
		{"1", color.Black, "#000000"},
		{"2", color.White, "#ffffff"},
		{"3", color.RGBA{255, 0, 0, 0}, "#ff0000"},
		{"4", color.RGBA{0, 255, 0, 0}, "#00ff00"},
		{"5", color.RGBA{0, 0, 255, 0}, "#0000ff"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			st.Expect(t, ColorToHex(tt.input), tt.want)
		})
	}
}
