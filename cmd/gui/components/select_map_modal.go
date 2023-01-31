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

func (s *SelectMapModal) HandleSelect(m Map) func() {

	return func() {
		s.OnSelect(m)
		s.Button.SetIcon(m.Icon)
		s.Button.SetText(fmt.Sprintf("(current map %d) select map", m.Number+1))
		s.Modal.Hide()
	}
}

func (s *SelectMapModal) Refresh() {
	s.Modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Objects = nil

	for i := range s.Maps {
		var layerButton *widget.Button

		if s.Maps[i].Empty {
			layerButton = widget.NewButtonWithIcon(fmt.Sprintf("Map %d", s.Maps[i].Number+1), s.Maps[i].Icon, s.HandleSelect(s.Maps[i]))
		} else {
			layerButton = widget.NewButtonWithIcon(fmt.Sprintf("Map %d (e)", s.Maps[i].Number+1), s.Maps[i].Icon, s.HandleSelect(s.Maps[i]))
		}
		s.Modal.Content.(*fyne.Container).Objects[1].(*fyne.Container).Add(layerButton)
	}
}

func (s *SelectMapModal) SetMaps(maps []Map) {
	s.Maps = maps
	s.Refresh()
}
