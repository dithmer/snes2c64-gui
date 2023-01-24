package main

import (
	"image/color"
	"snes2c64gui/cmd/gui/views"

	fyne "fyne.io/fyne/v2"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {

	myApp := app.New()

	myWindow := myApp.NewWindow("Snes2C64")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(true)

	// create a new theme based on the default theme
	customTheme := &CustomTheme{}
	// set the new theme
	myApp.Settings().SetTheme(customTheme)

	uploadView := views.NewUploadView(myWindow)
	uploadView.Draw(myWindow)

	myWindow.Show()
	myApp.Run()
}

type CustomTheme struct{}

var _ fyne.Theme = (*CustomTheme)(nil)

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameSeparator {
		// light gray
		return color.RGBA{0x80, 0x80, 0x80, 0xff}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
