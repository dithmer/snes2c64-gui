package components

import (
	"embed"
	"fmt"
	"image/color"
	"io"
	"log"
	"snes2c64gui/cmd/gui/widgets"
	"snes2c64gui/pkg/controller"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GamepadMapView struct {
	Container *fyne.Container

	selectedGamepadMap int
	gamepadMaps        []controller.GamepadMap
}

var c64buttons = []string{
	"joy_up",
	"joy_down",
	"joy_left",
	"joy_right",
	"btn_1",
	"btn_2",
	"btn_3",
	"btn_a",
}

//go:embed assets/*
var assets embed.FS

func NewGamepadMap(snesKeyImages []*canvas.Image) *GamepadMapView {
	// snesKeyCount is the number of keys existing on the snes controller
	snesKeyCount := len(snesKeyImages)

	// mainContainer is the max container to display the gamepad map alongside with an overlay to show the state of the gamepad

	gamepadMapContainer := container.NewHBox()

	for i := 0; i < snesKeyCount; i++ {
		gamepadMapColContainer := container.NewVBox()

		gamepadMapColContainer.Add(snesKeyImages[i])
		gamepadMapColContainer.Add(widget.NewSeparator())

		c64ButtonsContainer := container.NewVBox()
		for _, button := range c64buttons {
			resource, _ := assets.Open(fmt.Sprintf("assets/c64_%s.svg.png", button))
			b, err := io.ReadAll(resource)
			if err != nil {
				log.Fatalf("failed to read resource: %v", err)
			}
			if err := resource.Close(); err != nil {
				log.Fatalf("failed to close resource: %v", err)
			}

			staticResource := fyne.NewStaticResource(fmt.Sprintf("c64_%s", button), b)
			iconPressSwitch := widgets.NewIconPressSwitch(staticResource, 50, 50)

			c64ButtonsContainer.Add(iconPressSwitch)
		}

		gamepadMapColContainer.Add(c64ButtonsContainer)

		gamepadMapContainer.Add(layout.NewSpacer())
		gamepadMapContainer.Add(gamepadMapColContainer)
		gamepadMapContainer.Add(layout.NewSpacer())

		if i < snesKeyCount-1 {
			gamepadMapContainer.Add(widget.NewSeparator())
		}
	}

	mainContainer := container.NewMax()
	mainContainer.Add(gamepadMapContainer)

	// overlay background
	overlayBackground := canvas.NewRectangle(color.RGBA{0, 0, 0, 220})
	overlayBackground.Hide()
	mainContainer.Add(overlayBackground)

	overlayText := canvas.NewText("Gamepad Map", color.White)
	overlayText.Alignment = fyne.TextAlignCenter
	overlayText.TextSize = 20
	overlayText.Hide()
	mainContainer.Add(overlayText)

	return &GamepadMapView{
		Container: mainContainer,
	}
}

func (m *GamepadMapView) InfoOverlay(text string) {
	overlayRect := m.Container.Objects[1].(*canvas.Rectangle)
	overlayRect.Show()

	overlayText := m.Container.Objects[2].(*canvas.Text)
	overlayText.Text = text
	overlayText.Color = color.White
	overlayText.Show()

	m.Container.Refresh()
}

func (m *GamepadMapView) ErrorOverlay(text string) {
	overlayRect := m.Container.Objects[1].(*canvas.Rectangle)
	overlayRect.Show()

	overlayText := m.Container.Objects[2].(*canvas.Text)
	overlayText.Text = text
	overlayText.Color = color.RGBA{255, 0, 0, 255}
	overlayText.Show()

	m.Container.Refresh()
}

func (m *GamepadMapView) HideOverlay() {
	overlayRect := m.Container.Objects[1].(*canvas.Rectangle)
	overlayRect.Hide()

	overlayText := m.Container.Objects[2].(*canvas.Text)
	overlayText.Hide()

	m.Container.Refresh()
}

func (m *GamepadMapView) Enable() {
	for _, gamepadMapColContainer := range m.getGamepadMapColContainers() {
		c64ButtonsContainer := gamepadMapColContainer.Objects[2].(*fyne.Container)

		for _, button := range c64ButtonsContainer.Objects {
			button.(*widgets.IconPressSwitch).Enable()
		}
	}
}

func (m *GamepadMapView) getGamepadMapColContainers() []*fyne.Container {
	gamepadMapContainer := m.Container.Objects[0].(*fyne.Container)

	gamepadMapColContainers := make([]*fyne.Container, 0)

	for _, container := range gamepadMapContainer.Objects {
		if _, ok := container.(*fyne.Container); ok {
			gamepadMapColContainers = append(gamepadMapColContainers, container.(*fyne.Container))
		} else {
			continue
		}
	}

	return gamepadMapColContainers
}

func (m *GamepadMapView) Disable() {
	for _, number := range m.getGamepadMapColContainers() {
		c64ButtonsContainer := number.Objects[2].(*fyne.Container)

		for _, button := range c64ButtonsContainer.Objects {
			button.(*widgets.IconPressSwitch).Disable()
		}
	}
}

func (m *GamepadMapView) SelectGamepadMap(index int) {
	m.selectedGamepadMap = index

	gamepadMap := m.gamepadMaps[index]

	for i, number := range gamepadMap {
		c64ButtonsContainer := m.getGamepadMapColContainers()[i].Objects[2].(*fyne.Container)

		for j, button := range c64ButtonsContainer.Objects {
			button.(*widgets.IconPressSwitch).SetActive(int(number)&pow2(j) != 0)
		}
	}
}

func (m *GamepadMapView) Map() controller.GamepadMap {
	var gamepadMap controller.GamepadMap

	for i, number := range m.getGamepadMapColContainers() {
		c64ButtonsContainer := number.Objects[2].(*fyne.Container)

		for j, button := range c64ButtonsContainer.Objects {
			if button.(*widgets.IconPressSwitch).Active() {
				gamepadMap[i] |= (uint8(pow2(j)))
			}
		}
	}

	return gamepadMap
}

func pow2(n int) int {
	if n == 0 {
		return 1
	}

	return 2 * pow2(n-1)
}

func (m *GamepadMapView) SetGamepadMaps(gamepadMaps []controller.GamepadMap) {
	m.gamepadMaps = gamepadMaps

}

func (m *GamepadMapView) SelectedGamepadMap() int {
	return m.selectedGamepadMap
}
