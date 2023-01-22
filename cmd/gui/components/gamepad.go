package components

import (
	"embed"
	"fmt"
	"io"
	"log"
	"snes2c64gui/cmd/gui/widgets"
	"snes2c64gui/pkg/controller"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GamepadMap struct {
	*fyne.Container

	cols []*fyne.Container
}

var c64buttons = []string{
	"btn_1",
	"btn_2",
	"btn_3",
	"btn_a",
	"joy_up",
	"joy_down",
	"joy_left",
	"joy_right",
}

//go:embed assets/*
var assets embed.FS

func NewGamepadMap() *GamepadMap {
	keyCount := 10

	gamepadMap := &GamepadMap{
		container.NewHBox(),
		make([]*fyne.Container, keyCount),
	}

	for i := 0; i < keyCount; i++ {
		buttonsContainer := container.NewVBox()
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
			button.(*widgets.IconPressSwitch).Enable()
		}
	}
}

func (m *GamepadMap) Disable() {
	for _, number := range m.cols {
		for _, button := range number.Objects {
			button.(*widgets.IconPressSwitch).Disable()
		}
	}
}

func (m *Gamepad) Enable() {
	m.gamepadMapView.Enable()

	radioGroup := m.Container.Objects[0].(*widget.RadioGroup)
	radioGroup.Enable()
}

func (m *Gamepad) Disable() {
	m.gamepadMapView.Disable()

	radioGroup := m.Container.Objects[0].(*widget.RadioGroup)
	radioGroup.Disable()
}

func (m *GamepadMap) SetMap(gamepadMap controller.GamepadMap) {
	for i, number := range gamepadMap {
		for j := 0; j < 8; j++ {
			m.cols[i].Objects[j].(*widgets.IconPressSwitch).SetActive(int(number)&pow2(j) != 0)
		}
	}
}

func (m *GamepadMap) Map() controller.GamepadMap {
	var gamepadMap controller.GamepadMap

	for i, number := range m.cols {
		for j, button := range number.Objects {
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

type Gamepad struct {
	*fyne.Container

	gamepadMapView *GamepadMap

	selectedMap int
	GamepadMaps []controller.GamepadMap
}

func NewGamepad() *Gamepad {
	mapCount := 8

	gamepadMap := NewGamepadMap()

	gamepad := &Gamepad{
		container.NewVBox(
			widget.NewRadioGroup([]string{}, func(string) {}),
			gamepadMap.Container,
		),
		gamepadMap,
		0,
		make([]controller.GamepadMap, mapCount),
	}

	var options []string
	for i := 0; i < mapCount; i++ {
		options = append(options, fmt.Sprintf("Map %d", i))
	}

	radioGroup := gamepad.Container.Objects[0].(*widget.RadioGroup)
	radioGroup.Options = options
	radioGroup.Horizontal = true
	radioGroup.SetSelected("Map 0")
	radioGroup.Disable()
	radioGroup.OnChanged = gamepad.handleMapSelect()
	return gamepad
}

func (m *Gamepad) SetMaps(gamepadMaps []controller.GamepadMap) {
	m.GamepadMaps = gamepadMaps
}

func (m *Gamepad) SetSelectedMap(mapIndex int) {
	m.selectedMap = mapIndex

	radioGroup := m.Container.Objects[0].(*widget.RadioGroup)
	radioGroup.SetSelected(fmt.Sprintf("Map %d", mapIndex))

	m.gamepadMapView.SetMap(m.GamepadMaps[mapIndex])
}

func (m *Gamepad) SelectedMap() int {
	return m.selectedMap
}

func (m *Gamepad) Map() *GamepadMap {
	return m.gamepadMapView
}

func (m *Gamepad) handleMapSelect() func(value string) {
	return func(value string) {
		m.SetSelectedMap(int(value[4]) - 48)
	}
}
