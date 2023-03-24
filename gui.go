package main

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"path"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createGUI() {
	logf("Opening GUI")

	a := app.NewWithID("xyz.6gh.emptytrackcreator")
	window := a.NewWindow("Empty Track Creator")
	window.SetMaster()
	window.Resize(fyne.NewSize(750, 800))
	window.CenterOnScreen()

	// for development purposes
	// uncomment when building
	window.SetIcon(theme.FyneLogo())

	helpBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			icon := canvas.NewImageFromResource(window.Icon())
			icon.SetMinSize(fyne.NewSize(128, 128))
			icon.FillMode = canvas.ImageFillContain

			title := canvas.NewText("Empty Track Creator", color.White)
			title.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			title.Alignment = fyne.TextAlignCenter
			title.TextSize = 24

			version := canvas.NewText("v1.2.0", color.White)
			version.Alignment = fyne.TextAlignCenter
			version.TextSize = 16

			creator := canvas.NewText("Created by 6gh", color.White)
			creator.Alignment = fyne.TextAlignCenter

			fyneUrl, err := url.Parse("https://fyne.io/")
			handleErr(err)
			fyneLbl := widget.NewHyperlink("Made with Fyne", fyneUrl)
			fyneLbl.Alignment = fyne.TextAlignCenter

			repoUrl, err := url.Parse("https://github.com/6gh/Empty-Track-Creator")
			handleErr(err)
			repoLbl := widget.NewHyperlink("Check on GitHub", repoUrl)
			repoLbl.Alignment = fyne.TextAlignCenter

			hBox := container.New(layout.NewGridLayout(2), fyneLbl, repoLbl)
			vBox := container.New(layout.NewVBoxLayout(), icon, title, version, creator, hBox)

			dialog.ShowCustom("About", "Close", vBox, window)
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			title := canvas.NewText("Additional Settings", color.White)
			title.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			title.Alignment = fyne.TextAlignCenter
			title.TextSize = 24

			timerChkBox := widget.NewCheck("Enable Timer (shows how much time the program took)", func(checked bool) {
				a.Preferences().SetBool("timer", checked)
			})
			timerChkBox.SetChecked(a.Preferences().BoolWithFallback("timer", false))

			vBox := container.New(layout.NewVBoxLayout(), title, timerChkBox)

			dialog.ShowCustom("Settings", "Close", vBox, window)
		}),
	)

	// create all labels first
	MelodyTrackLbl := createTxt("Melody Tracks:")
	ArtTrackLbl := createTxt("Art Tracks:")
	PPQLbl := createTxt("PPQ:")
	BPMLbl := createTxt("BPM:")

	// create all inputs
	InputTXT := widget.NewEntry()
	InputTXT.Validator = func(s string) error {
		if s != "" && path.Ext(s) != ".mid" {
			return errors.New("file must be a .mid file")
		}
		return nil
	}
	OutputTXT := widget.NewEntry()
	OutputTXT.Validator = func(s string) error {
		if s == "" {
			return errors.New("path cannot be empty")
		}
		if path.Ext(s) != ".mid" {
			return errors.New("file must be a .mid file")
		}
		return nil
	}
	MelodyTrackTXT := createNumberInput(0, 500)
	ArtTrackTXT := createNumberInput(0, 500)
	PPQTXT := widget.NewSelect([]string{"96", "192", "240", "480", "960", "1920", "3840", "8192"}, func(string) {})
	BPMTXT := createNumberInput(0, 600)
	DrumsChk := widget.NewCheck("Allow Drums channel?", func(bool) {})

	OutputBox := widget.NewMultiLineEntry()
	OutputBox.SetText("Output will go here...")

	inputButton := widget.NewButtonWithIcon("Input Path", theme.FileIcon(), func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, _ error) {
			if reader != nil {
				InputTXT.SetText(reader.URI().Path())
			}
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mid"}))
		fileDialog.Show()
	})
	outputButton := widget.NewButtonWithIcon("Output Path", theme.FileIcon(), func() {
		fileDialog := dialog.NewFileSave(func(reader fyne.URIWriteCloser, _ error) {
			if reader != nil {
				p := reader.URI().Path()

				if path.Ext(p) != ".mid" {
					p += ".mid"
				}
				OutputTXT.SetText(p)
			}
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mid"}))
		fileDialog.Show()
	})

	createButton := widget.NewButton("Create", func() {
		var errors []string
		if err := InputTXT.Validate(); err != nil {
			errors = append(errors, "input: "+err.Error())
		}
		if err := OutputTXT.Validate(); err != nil {
			errors = append(errors, "output: "+err.Error())
		}
		if err := MelodyTrackTXT.Validate(); err != nil {
			errors = append(errors, "melody tracks: "+err.Error())
		}
		if err := ArtTrackTXT.Validate(); err != nil {
			errors = append(errors, "art tracks: "+err.Error())
		}
		if PPQTXT.Selected == "" {
			errors = append(errors, "ppq: cannot be empty")
		}
		if err := BPMTXT.Validate(); err != nil {
			errors = append(errors, "bpm: "+err.Error())
		}

		if len(errors) > 0 {
			dialog.ShowInformation("Invalid Options", strings.Join(errors, "\n"), window)
		} else {
			OutputBox.SetText("")

			melody, err := strconv.Atoi(MelodyTrackTXT.Text)
			handleErr(err)
			art, err := strconv.Atoi(ArtTrackTXT.Text)
			handleErr(err)
			pqq, err := strconv.Atoi(PPQTXT.Selected)
			handleErr(err)
			bpm, err := strconv.Atoi(MelodyTrackTXT.Text)
			handleErr(err)

			MelodyTrackTXT.Disable()

			ArtTrackTXT.Disable()
			OutputTXT.Disable()
			PPQTXT.Disable()
			BPMTXT.Disable()
			DrumsChk.Disable()
			InputTXT.Disable()
			inputButton.Disable()
			outputButton.Disable()
			window.SetTitle("Empty Track Creator (Running...)")

			createTracks(
				melody,
				art,
				OutputTXT.Text,
				pqq,
				bpm,
				DrumsChk.Checked,
				InputTXT.Text,
				a.Preferences().BoolWithFallback("timer", false),
				func(format string, a ...any) {
					OutputBox.SetText(OutputBox.Text + fmt.Sprintf(format, a...) + "\n")
				},
				func() {
					MelodyTrackTXT.Enable()
					ArtTrackTXT.Enable()
					OutputTXT.Enable()
					PPQTXT.Enable()
					BPMTXT.Enable()
					DrumsChk.Enable()
					InputTXT.Enable()
					inputButton.Enable()
					outputButton.Enable()
					window.SetTitle("Empty Track Creator")
				},
			)
		}
	})

	// set default values
	InputTXT.SetText(a.Preferences().StringWithFallback("inputPath", ""))
	OutputTXT.SetText(a.Preferences().StringWithFallback("outputPath", "output.mid"))
	MelodyTrackTXT.SetText(a.Preferences().StringWithFallback("melodyTracks", "8"))
	ArtTrackTXT.SetText(a.Preferences().StringWithFallback("artTracks", "8"))
	PPQTXT.SetSelected(a.Preferences().StringWithFallback("ppq", "960"))
	BPMTXT.SetText(a.Preferences().StringWithFallback("bpm", "138"))
	DrumsChk.SetChecked(a.Preferences().BoolWithFallback("allowDrums", false))

	// make rows
	inputRow := container.New(
		layout.NewFormLayout(),
		container.New(
			layout.NewHBoxLayout(),
			inputButton,
		),
		InputTXT,
	)
	outputRow := container.New(
		layout.NewFormLayout(),
		container.New(
			layout.NewHBoxLayout(),
			outputButton,
		), OutputTXT,
	)
	tracksRow := container.New(layout.NewGridLayout(2),
		container.New(layout.NewFormLayout(), MelodyTrackLbl, MelodyTrackTXT),
		container.New(layout.NewFormLayout(), ArtTrackLbl, ArtTrackTXT),
	)
	midiRow := container.New(layout.NewGridLayout(2),
		container.New(layout.NewFormLayout(), PPQLbl, PPQTXT),
		container.New(layout.NewFormLayout(), BPMLbl, BPMTXT),
	)
	drumRow := container.New(layout.NewHBoxLayout(), DrumsChk)
	bottomRow := container.New(
		layout.NewMaxLayout(),
		OutputBox,
	)

	// put into a column
	content := container.NewBorder(
		container.New(
			layout.NewVBoxLayout(),
			inputRow,
			outputRow,
			tracksRow,
			midiRow,
			drumRow,
			createButton,
		),
		helpBar,
		nil,
		nil,
		container.New(
			layout.NewMaxLayout(),
			bottomRow,
		),
	)

	window.SetCloseIntercept(func() {
		logf("GUI closed, saving settings")
		if OutputTXT.Text == "" {
			OutputTXT.SetText("output.mid")
		}
		if MelodyTrackTXT.Text == "" {
			MelodyTrackTXT.SetText("output.mid")
		}
		if ArtTrackTXT.Text == "" {
			ArtTrackTXT.SetText("output.mid")
		}
		if PPQTXT.Selected == "" {
			PPQTXT.SetSelected("960")
		}
		if BPMTXT.Text == "" {
			BPMTXT.SetText("output.mid")
		}

		a.Preferences().SetString("inputPath", InputTXT.Text)
		a.Preferences().SetString("outputPath", OutputTXT.Text)
		a.Preferences().SetString("melodyTracks", MelodyTrackTXT.Text)
		a.Preferences().SetString("artTracks", ArtTrackTXT.Text)
		a.Preferences().SetString("ppq", PPQTXT.Selected)
		a.Preferences().SetString("bpm", BPMTXT.Text)
		a.Preferences().SetBool("allowDrums", DrumsChk.Checked)

		window.Close()
	})

	window.SetContent(content)
	window.ShowAndRun()

}

// helper functions to create objects with the same settings
func createTxt(text string) *canvas.Text {
	txt := canvas.NewText(text, color.White)
	txt.Alignment = fyne.TextAlignLeading
	txt.TextSize = 16

	return txt
}

func createNumberInput(min int, max int) *widget.Entry {
	entry := widget.NewEntry()
	entry.Validator = func(input string) error {
		if input == "" {
			return errors.New("cannot be empty")
		}

		num, err := strconv.Atoi(input)

		if err != nil {
			return errors.New("not a number")
		}

		if num < min || num > max {
			return errors.New("number out of range")
		}

		return nil
	}
	return entry
}
