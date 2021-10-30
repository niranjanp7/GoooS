package main

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
)

func BoxLayout(box string) fyne.Layout {
	layouts := map[string]fyne.Layout{
		"H": layout.NewHBoxLayout(),
		"V": layout.NewVBoxLayout(),
	}
	return layouts[box]
}

func GridLayout(box string, n int) fyne.Layout {
	layouts := map[string]fyne.Layout{
		"":  layout.NewGridLayout(n),
		"C": layout.NewGridLayoutWithColumns(n),
		"R": layout.NewGridLayoutWithRows(n),
	}
	return layouts[box]
}

func Text(str string, r uint8, g uint8, b uint8, a uint8) *canvas.Text {
	return canvas.NewText(str, color.RGBA{R: r, G: g, B: b, A: a})
}

func main() {
	app := app.New().NewWindow("Calculator")

	// Calculation History
	var history []fyne.CanvasObject

	// Fields
	input := widget.NewEntry()
	input.FocusGained()
	output := widget.NewLabel("")
	historyItems := container.NewVBox(history...)
	historyScroll := container.NewVScroll(historyItems)
	historyScroll.SetMinSize(fyne.Size{Height: 100})
	historyScroll.Hide()
	historyAppear := false

	// Important Buttons
	historyBtn := widget.NewButtonWithIcon("History", theme.HistoryIcon(), func() {
		switch historyAppear {
		case true:
			historyScroll.Hide()
		case false:
			historyScroll.Show()
		}
		historyAppear = !historyAppear
	})
	backBtn := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		if len(input.Text) > 0 {
			input.SetText(input.Text[:len(input.Text)-1])
		}
	})
	evalBtn := widget.NewButton("=", func() {
		result, err := evalExp(input.Text)
		if err == nil {
			history = append(history, widget.NewLabel(input.Text+" = "+result))
			historyItems.AddObject(history[len(history)-1])
			historyItems.Refresh()
			input.SetText(result)
		} else {
			output.SetText(result)
		}
	})
	input.OnChanged = func(_ string) {
		result, err := evalExp(input.Text)
		if err == nil {
			output.SetText(result)
		}
		if input.Text == "" {
			output.SetText("")
		}
	}

	lightTheme := true
	themeOptions := map[bool]fyne.Theme{
		false: theme.DarkTheme(),
		true:  theme.LightTheme()}
	themeBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		lightTheme = !lightTheme
		fyne.CurrentApp().Settings().SetTheme(themeOptions[lightTheme])
	})

	// Set Buttons in Layout
	numrows := container.New(
		BoxLayout("V"),
		container.New(
			GridLayout("C", 2),
			historyBtn,
			backBtn,
		),
		container.New(
			GridLayout("C", 4),
			container.New(
				GridLayout("R", 3),
				widget.NewButtonWithIcon("Clear", theme.CancelIcon(), func() { input.SetText("") }),
				widget.NewButton("7", func() { input.SetText(input.Text + "7") }),
				widget.NewButton("4", func() { input.SetText(input.Text + "4") }),
			),
			container.New(
				GridLayout("R", 3),
				widget.NewButton("(", func() { input.SetText(input.Text + "(") }),
				widget.NewButton("8", func() { input.SetText(input.Text + "8") }),
				widget.NewButton("5", func() { input.SetText(input.Text + "5") }),
			),
			container.New(
				GridLayout("R", 3),
				widget.NewButton(")", func() { input.SetText(input.Text + ")") }),
				widget.NewButton("9", func() { input.SetText(input.Text + "9") }),
				widget.NewButton("6", func() { input.SetText(input.Text + "6") }),
			),
			container.New(
				GridLayout("R", 3),
				widget.NewButton("/", func() { input.SetText(input.Text + "/") }),
				widget.NewButton("*", func() { input.SetText(input.Text + "*") }),
				widget.NewButton("-", func() { input.SetText(input.Text + "-") }),
			),
		),
		container.New(
			GridLayout("C", 4),
			container.New(
				GridLayout("R", 2),
				widget.NewButton("1", func() { input.SetText(input.Text + "1") }),
				widget.NewButton("0", func() { input.SetText(input.Text + "0") }),
			),
			container.New(
				GridLayout("R", 2),
				widget.NewButton("2", func() { input.SetText(input.Text + "2") }),
				widget.NewButton(".", func() { input.SetText(input.Text + ".") }),
			),
			container.New(
				GridLayout("R", 2),
				widget.NewButton("3", func() { input.SetText(input.Text + "3") }),
				evalBtn,
			),
			container.New(
				GridLayout("C", 1),
				widget.NewButton("+", func() { input.SetText(input.Text + "+") }),
			),
		),
	)
	c := container.New(
		BoxLayout("V"),
		input,
		container.New(BoxLayout("H"), output, layout.NewSpacer(), themeBtn),
		historyScroll,
		numrows)
	app.SetContent(c) // Add Layouts to Window
	app.SetPadded(false)
	app.ShowAndRun()
}

func evalExp(exp string) (string, error) {
	ans, err := govaluate.NewEvaluableExpression(exp)
	if err == nil {
		result, e := ans.Evaluate(nil)
		if e == nil {
			return strconv.FormatFloat(result.(float64), 'f', -1, 64), e
		} else {
			return "ERROR : Expressorn Invalid", e
		}
	} else {
		return "ERROR : Expressorn Invalid", err
	}
}
