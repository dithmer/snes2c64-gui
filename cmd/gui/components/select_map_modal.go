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

	Maps     []Map
	OnSelect func(layer Map)
}

type Map struct {
	Number int
	Icon   fyne.Resource
	Empty  bool
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
		Maps:   maps,
		Button: open,
		Modal:  modal,
	}

	s.OnSelect = onSelect

	s.Refresh()

	return s
}

func (s *SelectMapModal) Refresh() {
	handleSelect := func(layer Map) func() {
		return func() {
			s.OnSelect(layer)
			s.Button.SetIcon(layer.Icon)
			s.Button.SetText(fmt.Sprintf("(current map %d) select map", layer.Number+1))
			s.Modal.Hide()
		}
	}

	s.Modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Objects = nil

	for i := range s.Maps {
		var layerButton *widget.Button

		if s.Maps[i].Empty {
			layerButton = widget.NewButtonWithIcon(fmt.Sprintf("Map %d", s.Maps[i].Number+1), s.Maps[i].Icon, handleSelect(s.Maps[i]))
		} else {
			layerButton = widget.NewButtonWithIcon(fmt.Sprintf("Map %d (e)", s.Maps[i].Number+1), s.Maps[i].Icon, handleSelect(s.Maps[i]))
		}
		s.Modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Add(layerButton)
	}
}

func (s *SelectMapModal) SetMaps(maps []Map) {
	s.Maps = maps
	s.Refresh()
}
