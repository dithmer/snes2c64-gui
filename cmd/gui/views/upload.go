package views

import (
	"embed"
	"fmt"
	"io"

	_ "embed"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"snes2c64gui/cmd/gui/components"
	"snes2c64gui/pkg/controller"
)

type UploadView struct {
	Controller *controller.Controller

	StatusBar *widget.Label
	LogsView  *components.Logs

	ConnectModalButton     fyne.CanvasObject
	SelectLayerModalButton fyne.CanvasObject
	CurrentMapLabel        *widget.Label
	CurrentMapIcon         *widget.Icon
	Gamepad                *components.Gamepad

	selectedMap int

	UploadButton *widget.Button
}

//go:embed assets/*
var assets embed.FS

var selectMapModalMapIcons = []string{
	"snes_button_b-full",
	"snes_button_y-full",
	"dpad_up",
	"dpad_down",
	"dpad_left",
	"dpad_right",
	"snes_button_a-full",
	"snes_button_x-full",
}

var keyIconNames = []string{
	"snes_button_b-full",
	"snes_button_y-full",
	"dpad_up",
	"dpad_down",
	"dpad_left",
	"dpad_right",
	"snes_button_a-full",
	"snes_button_x-full",
	"snes_shoulder_l",
	"snes_shoulder_R",
}

func NewUploadView(window fyne.Window) (uv *UploadView) {
	statusBar := widget.NewLabel("")

	connectModal := components.NewConnectModal(window.Canvas(), func(port string) {
		handleConnect(uv, uv.Controller, port)()
	})

	maps := make([]components.Map, len(selectMapModalMapIcons))

	currentMapLabel := widget.NewLabel("Current map: 1")
	currentMapIcon := widget.NewIcon(nil)

	keysIcons := make([]*canvas.Image, 10)

	for i, icon := range selectMapModalMapIcons {
		resource, err := assets.Open(fmt.Sprintf("assets/%s.png", icon))
		if err != nil {
			panic(fmt.Sprintf("Error opening asset: %v", err))
		}

		b, err := io.ReadAll(resource)
		if err != nil {
			panic(fmt.Sprintf("Error reading asset: %v", err))
		}

		if err := resource.Close(); err != nil {
			panic(fmt.Sprintf("Error closing asset: %v", err))
		}

		staticResource := fyne.NewStaticResource(icon, b)

		if i == 0 {
			currentMapIcon.SetResource(staticResource)
		}

		maps[i] = components.Map{
			Icon:   staticResource,
			Number: i,
		}
	}

	for i, icon := range keyIconNames {
		resource, err := assets.Open(fmt.Sprintf("assets/%s.svg.png", icon))
		if err != nil {
			panic(fmt.Sprintf("Error opening asset: %v", err))
		}

		b, err := io.ReadAll(resource)
		if err != nil {
			panic(fmt.Sprintf("Error reading asset: %v", err))
		}

		if err := resource.Close(); err != nil {
			panic(fmt.Sprintf("Error closing asset: %v", err))
		}

		staticResource := fyne.NewStaticResource(icon, b)

		keysIcons[i] = canvas.NewImageFromResource(staticResource)
		keysIcons[i].SetMinSize(fyne.NewSize(80, 80))
		keysIcons[i].FillMode = canvas.ImageFillContain
	}

	selectLayerModal := components.NewSelectMapModal(maps, window.Canvas(), func(layer components.Map) {
		uv.Gamepad.SetSelectedMap(layer.Number)
		uv.CurrentMapLabel.SetText(fmt.Sprintf("Current map: %d", layer.Number+1))
		uv.CurrentMapIcon.SetResource(layer.Icon)
		uv.selectedMap = layer.Number
	})
	selectLayerModal.(*components.SelectMapModal).Disable()

	gamepad := components.NewGamepad(keysIcons)

	uploadButton := widget.NewButton("Upload", func() {})
	uploadButton.Disable()

	logsView := components.NewLogs()

	defer func() {
		uv.SetStatus("Starting up")

		uv.UploadButton.OnTapped = handleUpload(uv)
	}()

	return &UploadView{
		StatusBar:              statusBar,
		ConnectModalButton:     connectModal,
		SelectLayerModalButton: selectLayerModal,
		Gamepad:                gamepad,
		LogsView:               logsView,
		UploadButton:           uploadButton,
		CurrentMapLabel:        currentMapLabel,
		CurrentMapIcon:         currentMapIcon,
	}
}

func (uv *UploadView) Draw(window fyne.Window) {

	window.SetContent(
		container.NewHBox(
			container.NewVBox(
				container.NewHBox(
					uv.ConnectModalButton,
					uv.SelectLayerModalButton,
					uv.CurrentMapLabel,
					uv.CurrentMapIcon,
				),
				uv.Gamepad.Container,
				uv.UploadButton,
			),
			uv.StatusBar,
		),
	)
}

func (uv *UploadView) SetStatus(status string) {
	uv.StatusBar.SetText("Status: " + status)
	uv.LogsView.Add(status)
}

func (uv *UploadView) EnableUpload() {
	uv.UploadButton.Enable()
}

func (uv *UploadView) Upload() {
	uv.SetStatus("Uploading")

	err := uv.Controller.Upload(uint8(uv.Gamepad.SelectedMap()), uv.Gamepad.Map().Map())
	if err != nil {
		uv.SetStatus(err.Error())
		return
	}

	uv.SetStatus("Upload complete")
}

func (uv *UploadView) Download() {
	uv.SetStatus("Downloading")

	gamepadMaps, err := uv.Controller.Download()
	if err != nil {
		uv.SetStatus(err.Error())
		return
	}

	uv.Gamepad.SetMaps(gamepadMaps)

	uv.SetStatus("Download complete")

}

func handleConnect(uv *UploadView, c *controller.Controller, port string) func() {
	return func() {
		var err error

		uv.SetStatus(fmt.Sprintf("Connecting to %s", port))

		c, err := controller.NewController(port)
		if err != nil {
			uv.SetStatus(fmt.Sprintf("Error connecting to %s: %v", port, err))
			return
		}
		uv.Controller = c

		uv.SetStatus("Connected")

		uv.Download()

		uv.Gamepad.SetSelectedMap(uv.selectedMap)
		uv.Gamepad.Enable()
		uv.SelectLayerModalButton.(*components.SelectMapModal).Enable()
		uv.SelectLayerModalButton.(*components.SelectMapModal).Modal.Show()
		uv.EnableUpload()
	}
}

func handleUpload(uv *UploadView) func() {
	return func() {
		uv.Upload()
		uv.Download()
	}
}
