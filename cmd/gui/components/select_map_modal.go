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

	layers   []Map
	OnSelect func(layer Map)
}

type Map struct {
	Number int
	Icon   fyne.Resource
}

func NewSelectMapModal(layers []Map, parent fyne.Canvas, onSelect func(layer Map)) *SelectMapModal {
	modal := widget.NewModalPopUp(nil, parent)

	open := widget.NewButtonWithIcon("select map", layers[0].Icon, func() {
		modal.Show()
	})

	modal.Content = container.NewVBox(
		widget.NewLabel("Select a map"),
		container.NewGridWithColumns(4),
	)

	s := &SelectMapModal{
		layers: layers,
		Button: open,
		Modal:  modal,
	}

	handleSelect := func(layer Map) func() {
		return func() {
			onSelect(layer)
			open.SetIcon(layer.Icon)
			modal.Hide()
		}
	}

	for i := range layers {
		layerButton := widget.NewButtonWithIcon(fmt.Sprintf("Map %d", layers[i].Number+1), layers[i].Icon, handleSelect(layers[i]))
		modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Add(layerButton)
	}

	return s
}
