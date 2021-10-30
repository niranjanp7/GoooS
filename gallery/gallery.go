package main

import (
	"image/color"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

func ImageFileList(dir string) (int, []string) {
	files, err := ioutil.ReadDir(dir)
	fileList := make([]string, 0)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				splitName := strings.Split(file.Name(), ".")
				extension := splitName[len(splitName)-1]
				if extension == "jpg" || extension == "jpeg" || extension == "png" {
					fileList = append(fileList, dir+file.Name())
				}
			}
		}
	}
	if len(fileList) == 0 {
		fileList = append(fileList, "./notfound/noimage.png")
		return 1, fileList
	}
	return len(fileList), fileList
}

func FileImageTile(files []string) []fyne.CanvasObject {
	images := make([]fyne.CanvasObject, 0)
	for _, file := range files {
		img := canvas.NewImageFromFile(file)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(50, 50))
		images = append(images, img)
	}
	return images
}

func ContImgList(img []fyne.CanvasObject, target *int, imgViewer *fyne.Container, imgFileList []string, nameTile *widget.Label) []*fyne.Container {
	cont := make([]*fyne.Container, 0)
	temp := make([]int, 0)
	for i := 0; i < len(img); i++ {
		temp = append(temp, i)
	}
	for k, i := range img {
		index := k
		cont = append(cont, container.NewCenter(widget.NewButton("           ", func() {
			*target = index
			img := canvas.NewImageFromFile(imgFileList[index])
			img.SetMinSize(fyne.NewSize(400, 400))
			img.FillMode = canvas.ImageFillContain
			imgViewer.AddObject(img)
			imgViewer.Objects = imgViewer.Objects[1:]
			nameTile.SetText(imgFileList[index])
		}), i))
	}
	return cont
}

func CrateGallery(dir string, fileChoser fyne.Widget) *fyne.Container {
	targetImage := 0
	lightTheme := true
	themeOptions := map[bool]fyne.Theme{
		false: theme.DarkTheme(),
		true:  theme.LightTheme()}
	themeBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		lightTheme = !lightTheme
		fyne.CurrentApp().Settings().SetTheme(themeOptions[lightTheme])
	})
	numOfImg, imgFileList := ImageFileList(dir)
	imageList := FileImageTile(imgFileList)
	showImage := canvas.NewImageFromFile(imgFileList[targetImage])
	showImage.SetMinSize(fyne.NewSize(400, 400))
	showImage.FillMode = canvas.ImageFillContain
	imgViewer := container.NewVBox(showImage)
	nameTile := widget.NewLabel(imgFileList[targetImage])
	prevImgBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if targetImage > 0 {
			targetImage--
			img := canvas.NewImageFromFile(imgFileList[targetImage])
			img.SetMinSize(fyne.NewSize(400, 400))
			img.FillMode = canvas.ImageFillContain
			imgViewer.AddObject(img)
			imgViewer.Objects = imgViewer.Objects[1:]
			nameTile.SetText(imgFileList[targetImage])
		}
	})
	nextImgBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		if targetImage < numOfImg-1 {
			targetImage++
			img := canvas.NewImageFromFile(imgFileList[targetImage])
			img.SetMinSize(fyne.NewSize(400, 400))
			img.FillMode = canvas.ImageFillContain
			imgViewer.AddObject(img)
			imgViewer.Objects = imgViewer.Objects[1:]
			nameTile.SetText(imgFileList[targetImage])
		}
	})
	contImgList := ContImgList(imageList, &targetImage, imgViewer, imgFileList, nameTile)

	imgTileLayout := container.NewHBox()
	for _, img := range contImgList {
		imgTileLayout.AddObject(img)
	}
	imageTileSrollBar := container.NewHScroll(imgTileLayout)
	imageTileSrollBar.SetMinSize(fyne.Size{Width: 450})
	c := container.NewVBox(
		container.NewHBox(fileChoser, nameTile, layout.NewSpacer(), themeBtn),
		layout.NewSpacer(),
		container.NewHBox(
			prevImgBtn,
			layout.NewSpacer(),
			imgViewer,
			layout.NewSpacer(),
			nextImgBtn,
		),
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			imageTileSrollBar,
			layout.NewSpacer(),
		),
	)
	return c
}

func main() {
	app := app.New().NewWindow("Gallery")
	dir := "C:/msys64/home/Niranjan/goproject/gallery/imagegallery/"
	c := CrateGallery(dir, fileChoserBtn(app))
	app.SetContent(c)
	app.Resize(fyne.NewSize(650, 500))
	app.ShowAndRun()
}

func fileChoserBtn(app fyne.Window) fyne.Widget {
	fileChoser := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			save_dir := ""
			if err != nil {
				dialog.ShowError(err, app)
				return
			}
			if dir != nil {
				save_dir = dir.Path()
			}
			if save_dir != "" {
				app.SetContent(CrateGallery(save_dir+"/", fileChoserBtn(app)))
				app.Content().Refresh()
			}
		}, app)
	})
	return fileChoser
}
