package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
}

func NewIconPressSwitch(icon fyne.Resource) *IconPressSwitch {
	i := &IconPressSwitch{Icon: icon, enabled: true, active: false}

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

	var tapAnimation *fyne.Animation
	if !i.active {
		tapAnimation = SetActiveAnimation(i.activeBackground, i)
	} else {
		tapAnimation = SetPassiveAnimation(i.activeBackground, i)
	}

	tapAnimation.Curve = fyne.AnimationEaseOut
	tapAnimation.Stop()
	tapAnimation.Start()

	i.active = !i.active

	if i.OnToggled != nil {
		i.OnToggled(i.active)
	}
}

func SetActiveAnimation(bg *canvas.Rectangle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationStandard, func(p float32) {
		mid := w.Size().Width
		size := mid * p
		bg.Move(fyne.NewPos(0, w.Size().Height-w.Size().Height/5))
		bg.Resize(fyne.NewSize(size, w.Size().Height/5))
		canvas.Refresh(bg)
	})
}

func SetPassiveAnimation(bg *canvas.Rectangle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationShort, func(p float32) {
		mid := w.Size().Width
		size := mid * (1 - p)
		bg.Move(fyne.NewPos(0, w.Size().Height-w.Size().Height/5))
		bg.Resize(fyne.NewSize(size, w.Size().Height/5))
		canvas.Refresh(bg)
	})
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
	i.activeBackground = canvas.NewRectangle(theme.PrimaryColor())

	r := &iconPressSwitchRenderer{
		background:       i.background,
		activeBackground: i.activeBackground,
		iconPressSwitch:  i,
	}

	i.activeBackground.Move(fyne.NewPos(0, i.Size().Height-i.Size().Height/5))
	i.activeBackground.Resize(fyne.NewSize(0, i.Size().Height/5))
	i.activeBackground.Refresh()
	i.activeBackground.Show()

	r.icon = canvas.NewImageFromResource(i.Icon)
	r.icon.FillMode = canvas.ImageFillContain
	r.icon.Resize(fyne.NewSize(50, 50))
	r.icon.Refresh()
	r.icon.Show()

	return r
}

type iconPressSwitchRenderer struct {
	background       *canvas.Rectangle
	activeBackground *canvas.Rectangle
	icon             *canvas.Image

	iconPressSwitch *IconPressSwitch
}

func (r *iconPressSwitchRenderer) Destroy() {
}

func (r *iconPressSwitchRenderer) Layout(size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	r.icon.Move(pos)

	r.icon.Resize(fyne.NewSize(size.Width, size.Height))
	r.background.Resize(fyne.NewSize(size.Width, size.Height))
}

func (r *iconPressSwitchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 40)
}

func (r *iconPressSwitchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.icon, r.activeBackground}
}

func (r *iconPressSwitchRenderer) Refresh() {
	r.activeBackground.Move(fyne.NewPos(0, r.iconPressSwitch.Size().Height-r.iconPressSwitch.Size().Height/5))

	var size float32
	if r.iconPressSwitch.active {
		size = r.iconPressSwitch.Size().Width
	} else {
		size = 0
	}
	r.activeBackground.Resize(fyne.NewSize(size, r.iconPressSwitch.Size().Height/5))
	r.activeBackground.Refresh()
	r.activeBackground.Show()

	r.icon = canvas.NewImageFromResource(r.iconPressSwitch.Icon)
	r.icon.FillMode = canvas.ImageFillContain
	r.icon.Resize(fyne.NewSize(50, 50))
	r.icon.Refresh()
	r.icon.Show()
}
