package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

type ConnectModal struct {
	Button    *widget.Button
	Modal     *widget.PopUp
	OnConnect func(port string)

	serialPortButtonGrid *fyne.Container
}

func NewConnectModal(parent fyne.Canvas, onConnect func(port string)) *ConnectModal {
	c := &ConnectModal{}

	modal := widget.NewModalPopUp(nil, parent)
	portsGrid := container.NewGridWithColumns(3)
	open := widget.NewButton("Connect", func() {
		c.RefreshPorts()
		modal.Show()
	})

	modal.Content = container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Select a serial port"),
			widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
				c.RefreshPorts()
			}),
		),
		portsGrid,
	)

	c.Button = open
	c.Modal = modal
	c.serialPortButtonGrid = portsGrid
	c.OnConnect = onConnect

	return c
}

func (c *ConnectModal) RefreshPorts() {
	serialPorts, err := serial.GetPortsList()
	if err != nil {
		// TODO: Handle error in ui
		panic(fmt.Sprintf("Error getting serial ports: %v", err))
	}

	c.serialPortButtonGrid.Objects = nil

	for i := range serialPorts {
		portButton := widget.NewButton(serialPorts[i], func() {
			c.Modal.Hide()
			c.OnConnect(serialPorts[i])
		})
		c.serialPortButtonGrid.Add(portButton)
	}
}
