package views

import (
	"embed"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"

	_ "embed"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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

	PrintCheatSheetButton *widget.Button

	VersionLabel *widget.Label
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
	window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyC, Modifier: fyne.KeyModifierAlt}, func(shortcut fyne.Shortcut) {
		connectModal.RefreshPorts()
		connectModal.Modal.Show()
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
	shortcutKeys := []fyne.KeyName{
		fyne.Key1,
		fyne.Key2,
		fyne.Key3,
		fyne.Key4,
		fyne.Key5,
		fyne.Key6,
		fyne.Key7,
		fyne.Key8,
	}
	shortcutHandler := func(i int) func(fyne.Shortcut) {
		return func(shortcut fyne.Shortcut) {
			selectLayerModal.HandleSelect(maps[i])()
		}
	}
	for i := 0; i < len(shortcutKeys); i++ {
		window.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  shortcutKeys[i],
			Modifier: fyne.KeyModifierAlt,
		}, shortcutHandler(i))
	}

	clearMapButton := widget.NewButton("Clear Map", func() {
		uv.GamepadMapView.ClearSelectedMap()
	})
	clearMapButton.Disable()
	window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyX, Modifier: fyne.KeyModifierAlt}, func(shortcut fyne.Shortcut) {
		uv.GamepadMapView.ClearSelectedMap()
	})

	gamepad := components.NewGamepadMap(keysIcons)
	gamepad.InfoOverlay("Please connect the device to start")
	gamepad.Disable()

	uploadButton := widget.NewButton("Upload", func() {})
	uploadButton.Disable()
	defer func() {
		uv.UploadButton.OnTapped = handleUpload(uv)
	}()
	window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: fyne.KeyModifierAlt}, func(shortcut fyne.Shortcut) {
		handleUpload(uv)()
	})

	printCheatSheetButton := widget.NewButton("Print Cheat Sheet", func() {
		printCheatSheet(uv)
	})
	printCheatSheetButton.Disable()
	window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyP, Modifier: fyne.KeyModifierAlt}, func(shortcut fyne.Shortcut) {
		printCheatSheet(uv)
	})

	versionLabel := widget.NewLabel("")

	return &UploadView{
		ConnectModal:          connectModal,
		GamepadMapView:        gamepad,
		SelectLayerModal:      selectLayerModal,
		ClearMapButton:        clearMapButton,
		UploadButton:          uploadButton,
		PrintCheatSheetButton: printCheatSheetButton,
		VersionLabel:          versionLabel,
	}
}

func printCheatSheet(uv *UploadView) {
	var err error

	url := uv.GamepadMapView.GetCheatSheetURL()

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		uv.GamepadMapView.ErrorOverlay(fmt.Sprintf("Error opening cheat sheet: %v", err))
		go func() {
			time.Sleep(2 * time.Second)
			uv.GamepadMapView.HideOverlay()
		}()
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
				uv.PrintCheatSheetButton,
				container.NewHBox(
					layout.NewSpacer(),
					uv.VersionLabel,
				),
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

	uv.PrintCheatSheetButton.Disable()
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

	maps := uv.SelectLayerModal.Maps
	for i := range maps {
		maps[i].Empty = uv.GamepadMapView.IsEmpty(i)
	}
	uv.SelectLayerModal.SetMaps(maps)
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

		uv.ClearMapButton.Enable()

		uv.EnableUpload()

		uv.PrintCheatSheetButton.Enable()

		uv.ConnectModal.Button.SetText(fmt.Sprintf("Connected to %s", port))

		firmwareVersion, err := uv.Controller.GetFirmwareVersion()
		if err != nil {
			uv.GamepadMapView.ErrorOverlay(fmt.Sprintf("Error getting firmware version: %v", err))

			go func() {
				<-time.After(2 * time.Second)
				uv.Reset()
			}()

			return
		}

		uv.VersionLabel.SetText(strings.ReplaceAll(firmwareVersion, "\n", " "))
	}
}

func handleUpload(uv *UploadView) func() {
	return func() {
		uv.Upload()
		uv.Download()
	}
}
