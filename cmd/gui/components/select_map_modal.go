package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SelectMapModal struct {
	Button *widget.Button

	Modal *widget.PopUp

	maps     []Map
	OnSelect func(layer Map)
}

type Map struct {
	Number int
	Icon   fyne.Resource
}

func NewSelectMapModal(maps []Map, parent fyne.Canvas, onSelect func(layer Map)) *SelectMapModal {
	modal := widget.NewModalPopUp(nil, parent)

	open := widget.NewButtonWithIcon(fmt.Sprintf("(current map %d) select map", maps[0].Number+1), maps[0].Icon, func() {
		modal.Show()
	})

	modal.Content = container.NewVBox(
		widget.NewLabel("Select a map"),
		container.NewGridWithColumns(4),
	)

	s := &SelectMapModal{
		maps:   maps,
		Button: open,
		Modal:  modal,
	}

	handleSelect := func(layer Map) func() {
		return func() {
			onSelect(layer)
			open.SetIcon(layer.Icon)
			open.SetText(fmt.Sprintf("(current map %d) select map", layer.Number+1))
			modal.Hide()
		}
	}

	for i := range maps {
		layerButton := widget.NewButtonWithIcon(fmt.Sprintf("Map %d", maps[i].Number+1), maps[i].Icon, handleSelect(maps[i]))
		modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Add(layerButton)
	}

	return s
}
