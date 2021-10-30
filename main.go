package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()            // Crate New App
	w := a.NewWindow("VarOS") // Create New Window

	//Test Object
	var text [4]*canvas.Text
	for i := 0; i < 4; i++ {
		text[i] = canvas.NewText(fmt.Sprint(i), color.Black)
	}

	// Box Layout
	hbox := layout.NewHBoxLayout() // Horizontal Box Layout
	vbox := layout.NewVBoxLayout() // Vertical Box Layout

	// Left Layout
	leftBox := container.New(
		vbox,
		widget.NewIcon(theme.FileApplicationIcon()),
		widget.NewIcon(theme.FileApplicationIcon()),
		widget.NewIcon(theme.FileApplicationIcon()))

	// Right Layout
	rightLayout := container.New(
		vbox,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.HomeIcon(), func() {}),
		widget.NewButtonWithIcon("", theme.ComputerIcon(), func() {}),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.LogoutIcon(), func() { w.Close() }))

	// Combine all Layouts
	c := container.New(
		hbox,
		leftBox,
		layout.NewSpacer(),
		rightLayout)
	w.SetContent(c)       // Add Layouts to Window
	w.SetFullScreen(true) // Set Window to Full Screen
	w.SetPadded(false)
	w.ShowAndRun() // Run Window
}
