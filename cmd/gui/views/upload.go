package views

import (
	"embed"
	"fmt"
	"io"
	"time"

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

	ConnectModal *components.ConnectModal

	GamepadMapView *components.GamepadMapView

	SelectLayerModal *components.SelectMapModal
	ClearMapButton   *widget.Button
	UploadButton     *widget.Button
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
		uv.GamepadMapView.SelectGamepadMap(layer.Number)
	})
	selectLayerModal.Button.Disable()

	clearMapButton := widget.NewButton("Clear Map", func() {
		uv.GamepadMapView.ClearSelectedMap()
	})
	clearMapButton.Disable()

	gamepad := components.NewGamepadMap(keysIcons)
	gamepad.InfoOverlay("Please connect the device to start")
	gamepad.Disable()

	uploadButton := widget.NewButton("Upload", func() {})
	uploadButton.Disable()

	defer func() {
		uv.UploadButton.OnTapped = handleUpload(uv)
	}()

	return &UploadView{
		ConnectModal:     connectModal,
		GamepadMapView:   gamepad,
		SelectLayerModal: selectLayerModal,
		ClearMapButton:   clearMapButton,
		UploadButton:     uploadButton,
	}
}

func (uv *UploadView) Draw(window fyne.Window) {
	bottomButtonsGrid := container.New(layout.NewGridLayout(3), uv.SelectLayerModal.Button, uv.ClearMapButton, uv.UploadButton)

	window.SetContent(
		container.NewHBox(
			container.NewVBox(
				uv.ConnectModal.Button,
				layout.NewSpacer(),
				uv.GamepadMapView.Container,
				layout.NewSpacer(),
				bottomButtonsGrid,
			),
		),
	)
}

func (uv *UploadView) EnableUpload() {
	uv.UploadButton.Enable()
}

func (uv *UploadView) Reset() {
	uv.UploadButton.Disable()
	uv.SelectLayerModal.Button.Disable()
	uv.ClearMapButton.Disable()

	uv.GamepadMapView.InfoOverlay("Please connect the device to start")
	uv.GamepadMapView.Disable()

	uv.ConnectModal.Button.Enable()
	uv.ConnectModal.Button.SetText("Connect")
}

func (uv *UploadView) Upload() {
	uv.GamepadMapView.InfoOverlay(fmt.Sprintf("Uploading map %d", uv.GamepadMapView.SelectedGamepadMap()+1))
	err := uv.Controller.Upload(uint8(uv.GamepadMapView.SelectedGamepadMap()), uv.GamepadMapView.Map())
	if err != nil {
		uv.GamepadMapView.ErrorOverlay(fmt.Sprintf("Error uploading gamepad map %d: %v", uv.GamepadMapView.SelectedGamepadMap()+1, err))

		go func() {
			<-time.After(2 * time.Second)
			uv.Reset()
		}()
		return
	}

	uv.GamepadMapView.InfoOverlay(fmt.Sprintf("Map %d uploaded", uv.GamepadMapView.SelectedGamepadMap()+1))
	go func() {
		<-time.After(1 * time.Second)
		uv.GamepadMapView.HideOverlay()
	}()
}

func (uv *UploadView) Download() {
	gamepadMaps, err := uv.Controller.Download()
	if err != nil {
		uv.GamepadMapView.ErrorOverlay(fmt.Sprintf("Error downloading gamepad maps: %v", err))

		go func() {
			<-time.After(2 * time.Second)
			uv.Reset()
		}()
		return
	}

	uv.GamepadMapView.SetGamepadMaps(gamepadMaps)
}

func handleConnect(uv *UploadView, c *controller.Controller, port string) func() {
	return func() {
		var err error

		uv.GamepadMapView.InfoOverlay(fmt.Sprintf("Connecting to %s...", port))

		if c != nil {
			c.Close()
		}

		c, err := controller.NewController(port)
		if err != nil {
			uv.GamepadMapView.Disable()
			uv.GamepadMapView.ErrorOverlay(fmt.Sprintf("Error connecting to controller: %v", err))

			go func() {
				<-time.After(2 * time.Second)
				uv.Reset()
			}()
			return
		}
		uv.Controller = c

		uv.GamepadMapView.InfoOverlay("Downloading gamepad maps...")
		uv.Download()

		uv.GamepadMapView.SelectGamepadMap(uv.GamepadMapView.SelectedGamepadMap())
		uv.GamepadMapView.Enable()
		uv.GamepadMapView.HideOverlay()

		uv.SelectLayerModal.Button.Enable()
		uv.SelectLayerModal.Modal.Show()

		uv.EnableUpload()

		uv.ConnectModal.Button.SetText(fmt.Sprintf("Connected to %s", port))
	}
}

func handleUpload(uv *UploadView) func() {
	return func() {
		uv.Upload()
		uv.Download()
	}
}
