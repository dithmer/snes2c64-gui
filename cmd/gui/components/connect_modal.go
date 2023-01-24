package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

type ConnectModal struct {
	Button *widget.Button

	Modal *widget.PopUp

	OnConnect func(port string)
}

func NewConnectModal(parent fyne.Canvas, onConnect func(port string)) *ConnectModal {
	serialPorts, err := serial.GetPortsList()
	if err != nil {
		// TODO: Handle error in ui
		panic(fmt.Sprintf("Error getting serial ports: %v", err))
	}

	modal := widget.NewModalPopUp(nil, parent)

	portsGrid := container.NewGridWithColumns(3)
	for i := range serialPorts {
		portButton := widget.NewButton(serialPorts[i], func() {
			modal.Hide()
			onConnect(serialPorts[i])
		})
		portsGrid.Add(portButton)
	}

	open := widget.NewButton("Connect", func() {
		modal.Show()
	})

	modal.Content = container.NewVBox(
		widget.NewLabel("Select a serial port"),
		portsGrid,
	)

	return &ConnectModal{
		Button: open,
		Modal:  modal,
	}

}
