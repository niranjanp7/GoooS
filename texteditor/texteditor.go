package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	w := app.NewWindow("Text Editor")
	fileStatus := FileStatus{false, false, nil}
	c, input := Content(&fileStatus)
	w.SetContent(c)
	w.SetMainMenu(makeMenu(app, w, input, &fileStatus))
	w.Resize(fyne.NewSize(540, 400))
	w.SetOnClosed(func() {
		if fileStatus.uri != nil {
			fileStatus.uri.Close()
		}
	})
	w.SetPadded(false)
	w.ShowAndRun()
}

type FileStatus struct {
	saved  bool
	edited bool
	uri    fyne.URIWriteCloser
}

func Content(fileStatus *FileStatus) (*fyne.Container, *widget.Entry) {
	input := widget.NewMultiLineEntry()
	input.SetText("")
	input.FocusGained()
	input.OnChanged = func(_ string) {
		*&fileStatus.edited = true
	}
	return container.New(layout.NewGridLayoutWithColumns(1), input), input
}

func NewWindowOpen(a fyne.App, w fyne.Window, fileStatus *FileStatus, material string) {
	newFileStatus := FileStatus{false, false, nil}
	neww := a.NewWindow("New File")
	c, newInput := Content(&newFileStatus)
	newInput.SetText(material)
	neww.SetMainMenu(makeMenu(a, w, newInput, &newFileStatus))
	neww.Resize(fyne.NewSize(540, 400))
	neww.SetContent(c)
	neww.SetOnClosed(func() {
		if fileStatus.uri != nil {
			fileStatus.uri.Close()
		}
	})
	neww.SetPadded(false)
	neww.Show()
}

func makeMenu(a fyne.App, w fyne.Window, input *widget.Entry, fileStatus *FileStatus) *fyne.MainMenu {
	newFile := fyne.NewMenuItem("New", func() {
		NewWindowOpen(a, w, fileStatus, "")
	})
	openFile := fyne.NewMenuItem("Open", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})
	saveFile := fyne.NewMenuItem("Save", func() {
		if fileStatus.uri != nil {
			fileSaved(fileStatus.uri, input.Text, w)
		} else {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if writer == nil {
					log.Println("Cancelled")
					return
				}
				if writer != nil {
					fileSaved(writer, input.Text, w)
					*&fileStatus.uri = writer
					*&fileStatus.saved = true
					*&fileStatus.edited = false
				}
			}, w)
		}
	})
	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	})
	file := fyne.NewMenu("File", newFile, openFile, saveFile)
	if !fyne.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	return fyne.NewMainMenu(
		file,
	)
}

func fileSaved(f fyne.URIWriteCloser, file string, w fyne.Window) {
	//defer f.Close()
	_, err := f.Write([]byte(file))
	if err != nil {
		dialog.ShowError(err, w)
	}
	//err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
}
