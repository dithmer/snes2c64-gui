package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type IconPressSwitch struct {
	widget.BaseWidget

	Icon fyne.Resource

	OnToggled func(bool)

	background       *canvas.Rectangle
	activeBackground *canvas.Rectangle
	active, enabled  bool

	Width, Height int
}

func NewIconPressSwitch(icon fyne.Resource, width int, height int) *IconPressSwitch {
	i := &IconPressSwitch{Icon: icon, enabled: true, active: false, Width: width, Height: height}
	i.ExtendBaseWidget(i)
	return i
}

func (i *IconPressSwitch) Active() bool {
	return i.active
}

func (i *IconPressSwitch) Tapped(p *fyne.PointEvent) {
	if !i.enabled {
		return
	}

	// press must be within the background
	if p.Position.X < 0 || p.Position.Y < 0 || p.Position.X > i.background.Size().Width || p.Position.Y > i.background.Size().Height {
		return
	}

	i.active = !i.active

	if i.OnToggled != nil {
		i.OnToggled(i.active)
	}
	i.Refresh()
}

func (i *IconPressSwitch) Enable() {
	i.enabled = true
	i.Refresh()
}

func (i *IconPressSwitch) Disable() {
	i.enabled = false
	i.Refresh()
}

func (i *IconPressSwitch) SetActive(active bool) {
	i.active = active
	i.Refresh()
}

func (i *IconPressSwitch) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)

	i.background = canvas.NewRectangle(theme.BackgroundColor())

	i.activeBackground = canvas.NewRectangle(color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 200})

	r := &iconPressSwitchRenderer{
		iconPressSwitch:  i,
		background:       i.background,
		activeBackground: i.activeBackground,
		layout:           layout.NewVBoxLayout(),
	}

	r.icon = canvas.NewImageFromResource(i.Icon)
	r.icon.FillMode = canvas.ImageFillContain
	r.icon.Refresh()
	r.icon.Show()

	return r
}

type iconPressSwitchRenderer struct {
	background       *canvas.Rectangle
	activeBackground *canvas.Rectangle
	icon             *canvas.Image

	iconPressSwitch *IconPressSwitch

	layout fyne.Layout
}

func (r *iconPressSwitchRenderer) Destroy() {
}

func (r *iconPressSwitchRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.SetMinSize(size)

	r.activeBackground.SetMinSize(size)
	r.activeBackground.Resize(size)

	r.icon.SetMinSize(size)
	r.icon.Resize(size)

}

func (r *iconPressSwitchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(r.iconPressSwitch.Width), float32(r.iconPressSwitch.Height))
}

func (r *iconPressSwitchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.icon, r.activeBackground}
}

func (r *iconPressSwitchRenderer) Refresh() {
	r.background.Refresh()

	r.activeBackground.Refresh()
	if r.iconPressSwitch.active {
		r.activeBackground.Hide()
	} else {
		r.activeBackground.Show()
	}

	r.icon.Refresh()
	r.icon.Show()

	canvas.Refresh(r.iconPressSwitch)
}
