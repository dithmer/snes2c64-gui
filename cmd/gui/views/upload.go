package views

import (
	"embed"
	"fmt"
	"io"
	"log"

	_ "embed"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"snes2c64gui/cmd/gui/components"
	"snes2c64gui/pkg/controller"
)

type UploadView struct {
	Controller *controller.Controller

	LogsView *components.Logs

	ConnectModal     *components.ConnectModal
	SelectLayerModal *components.SelectMapModal

	Gamepad *components.Gamepad

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
	connectModal := components.NewConnectModal(window.Canvas(), func(port string) {
		handleConnect(uv, uv.Controller, port)()
	})

	maps := make([]components.Map, len(selectMapModalMapIcons))

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
		uv.selectedMap = layer.Number
	})
	selectLayerModal.Button.Disable()

	gamepad := components.NewGamepad(keysIcons)

	uploadButton := widget.NewButton("Upload", func() {})
	uploadButton.Disable()

	logsView := components.NewLogs()

	defer func() {
		uv.UploadButton.OnTapped = handleUpload(uv)
	}()

	return &UploadView{
		ConnectModal:     connectModal,
		SelectLayerModal: selectLayerModal,
		Gamepad:          gamepad,
		LogsView:         logsView,
		UploadButton:     uploadButton,
	}
}

func (uv *UploadView) Draw(window fyne.Window) {
	bottomButtonsGrid := container.New(layout.NewGridLayout(2), uv.SelectLayerModal.Button, uv.UploadButton)

	window.SetContent(
		container.NewHBox(
			container.NewVBox(
				uv.ConnectModal.Button,
				layout.NewSpacer(),
				uv.Gamepad.Container,
				layout.NewSpacer(),
				bottomButtonsGrid,
			),
		),
	)
}

func (uv *UploadView) EnableUpload() {
	uv.UploadButton.Enable()
}

func (uv *UploadView) Upload() {
	err := uv.Controller.Upload(uint8(uv.Gamepad.SelectedMap()), uv.Gamepad.Map().Map())
	if err != nil {
		log.Fatalf("Error uploading: %v", err)
		return
	}
}

func (uv *UploadView) Download() {
	gamepadMaps, err := uv.Controller.Download()
	if err != nil {
		log.Fatalf("Error downloading: %v", err)
		return
	}

	uv.Gamepad.SetMaps(gamepadMaps)
}

func handleConnect(uv *UploadView, c *controller.Controller, port string) func() {
	return func() {
		var err error

		c, err := controller.NewController(port)
		if err != nil {
			log.Fatalf("Error connecting to %s: %v", port, err)
			return
		}
		uv.Controller = c

		uv.Download()

		uv.Gamepad.SetSelectedMap(uv.selectedMap)
		uv.Gamepad.Enable()
		uv.SelectLayerModal.Button.Enable()
		uv.SelectLayerModal.Modal.Show()
		uv.EnableUpload()
	}
}

func handleUpload(uv *UploadView) func() {
	return func() {
		uv.Upload()
		uv.Download()
	}
}
