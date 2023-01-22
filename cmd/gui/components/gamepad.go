package components

import (
	"embed"
	"fmt"
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

type GamepadMap struct {
	*fyne.Container

	cols []*fyne.Container
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

func NewGamepadMap(keys []*canvas.Image) *GamepadMap {
	keyCount := 10

	gamepadMap := &GamepadMap{
		container.NewHBox(),
		make([]*fyne.Container, keyCount),
	}

	for i := 0; i < keyCount; i++ {
		buttonsContainer := container.NewVBox()

		// create icon from resource and add it to buttonsContainer
		buttonsContainer.Add(keys[i])
		buttonsContainer.Add(widget.NewSeparator())

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
			iconPressSwitch := widgets.NewIconPressSwitch(staticResource)

			buttonsContainer.Add(iconPressSwitch)
		}

		gamepadMap.cols[i] = buttonsContainer

		gamepadMap.Container.Add(layout.NewSpacer())
		gamepadMap.Container.Add(buttonsContainer)
		gamepadMap.Container.Add(layout.NewSpacer())

		if i < keyCount-1 {
			gamepadMap.Container.Add(widget.NewSeparator())
		}
	}

	return gamepadMap
}

func (m *GamepadMap) Enable() {
	for _, number := range m.cols {
		for _, button := range number.Objects {
			if _, ok := button.(*widgets.IconPressSwitch); ok {
				button.(*widgets.IconPressSwitch).Enable()
			}
		}
	}
}

func (m *GamepadMap) Disable() {
	for _, number := range m.cols {
		for _, button := range number.Objects {
			if _, ok := button.(*widgets.IconPressSwitch); ok {
				button.(*widgets.IconPressSwitch).Disable()
			}
		}
	}
}

func (m *Gamepad) Enable() {
	m.gamepadMapView.Enable()
}

func (m *Gamepad) Disable() {
	m.gamepadMapView.Disable()

	radioGroup := m.Container.Objects[0].(*widget.RadioGroup)
	radioGroup.Disable()
}

func (m *GamepadMap) SetMap(gamepadMap controller.GamepadMap) {
	for i, number := range gamepadMap {
		var j int
		for _, button := range m.cols[i].Objects {
			log.Printf("i: %d, j: %d, number: %d", i, j, number)
			if _, ok := button.(*widgets.IconPressSwitch); ok {
				button.(*widgets.IconPressSwitch).SetActive(int(number)&pow2(j) != 0)
				j++
			}
		}
	}
}

func (m *GamepadMap) Map() controller.GamepadMap {
	var gamepadMap controller.GamepadMap

	for i, number := range m.cols {
		var j int
		for _, button := range number.Objects {
			if _, ok := button.(*widgets.IconPressSwitch); ok {
				if button.(*widgets.IconPressSwitch).Active() {
					gamepadMap[i] |= (uint8(pow2(j)))
				}
				j++
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

type Gamepad struct {
	*fyne.Container

	gamepadMapView *GamepadMap

	selectedMap int
	GamepadMaps []controller.GamepadMap
}

func NewGamepad(keys []*canvas.Image) *Gamepad {
	mapCount := 8

	gamepadMap := NewGamepadMap(keys)

	gamepad := &Gamepad{
		container.NewVBox(
			gamepadMap.Container,
		),
		gamepadMap,
		0,
		make([]controller.GamepadMap, mapCount),
	}

	return gamepad
}

func (m *Gamepad) SetMaps(gamepadMaps []controller.GamepadMap) {
	m.GamepadMaps = gamepadMaps
}

func (m *Gamepad) SetSelectedMap(mapIndex int) {
	m.selectedMap = mapIndex

	m.gamepadMapView.SetMap(m.GamepadMaps[mapIndex])
}

func (m *Gamepad) SelectedMap() int {
	return m.selectedMap
}

func (m *Gamepad) Map() *GamepadMap {
	return m.gamepadMapView
}
