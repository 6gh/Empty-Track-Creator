package main

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	sqdialog "github.com/sqweek/dialog"
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
	// window.SetIcon(theme.FyneLogo())

	helpBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			logf("Opening help dialog")

			icon := canvas.NewImageFromResource(window.Icon())
			icon.SetMinSize(fyne.NewSize(128, 128))
			icon.FillMode = canvas.ImageFillContain

			title := canvas.NewText("Empty Track Creator", color.White)
			title.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			title.Alignment = fyne.TextAlignCenter
			title.TextSize = 24

			version := canvas.NewText(fyne.CurrentApp().Metadata().Version, color.White)
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
			logf("Opening settings dialog")

			title := canvas.NewText("Additional Settings", color.White)
			title.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			title.Alignment = fyne.TextAlignCenter
			title.TextSize = 24

			drumsChk := widget.NewCheck("Allow Drums channel?", func(bool) {})

			melodyTracksRange := widget.NewEntry()
			melodyTracksRange.Validator = func(s string) error {
				if s == "" {
					return errors.New("range cannot be empty")
				}
				if !strings.Contains(s, "-") {
					return errors.New("range must be in the format of <min>-<max>")
				}

				split := strings.Split(s, "-")
				if len(split) != 2 {
					return errors.New("range must be in the format of <min>-<max>")
				}

				min, err := strconv.Atoi(split[0])
				if err != nil {
					return errors.New("min is not a number")
				}

				max, err := strconv.Atoi(split[1])
				if err != nil {
					return errors.New("max is not a number")
				}

				if max > 16 {
					return errors.New("max cannot be greater than 16")
				}
				if min < 1 {
					return errors.New("min cannot be less than 1")
				}

				if min > max {
					return errors.New("min cannot be greater than max")
				}

				return nil
			}
			artTrackRange := widget.NewEntry()
			artTrackRange.Validator = func(s string) error {
				if s == "" {
					return errors.New("range cannot be empty")
				}
				if !strings.Contains(s, "-") {
					return errors.New("range must be in the format of <min>-<max>")
				}

				split := strings.Split(s, "-")
				if len(split) != 2 {
					return errors.New("range must be in the format of <min>-<max>")
				}

				min, err := strconv.Atoi(split[0])
				if err != nil {
					return errors.New("min is not a number")
				}

				max, err := strconv.Atoi(split[1])
				if err != nil {
					return errors.New("max is not a number")
				}

				if max > 16 {
					return errors.New("max cannot be greater than 16")
				}
				if min < 1 {
					return errors.New("min cannot be less than 1")
				}

				if min > max {
					return errors.New("min cannot be greater than max")
				}

				return nil
			}

			melodyTracksRange.SetText(a.Preferences().StringWithFallback("melodyTracksRange", "1-15"))
			artTrackRange.SetText(a.Preferences().StringWithFallback("artTracksRange", "16-16"))
			drumsChk.Checked = a.Preferences().BoolWithFallback("allowDrums", false)

			dialog.ShowForm("Settings", "Save", "Cancel", []*widget.FormItem{
				{
					Text:     "Melody Tracks Range",
					Widget:   melodyTracksRange,
					HintText: "The range of channels to create melody tracks on",
				},
				{
					Text:     "Art Tracks Range",
					Widget:   artTrackRange,
					HintText: "The range of channels to create art tracks on",
				},
				{
					Text:     "CH-10",
					Widget:   drumsChk,
					HintText: "If unchecked, channel 10 will be skipped",
				},
			}, func(b bool) {
				if b {
					a.Preferences().SetString("melodyTracksRange", melodyTracksRange.Text)
					a.Preferences().SetString("artTracksRange", artTrackRange.Text)
					a.Preferences().SetBool("allowDrums", drumsChk.Checked)
					logf("Settings closed and saved")
				}
			}, window)
		}),
	)

	// create all labels first
	MelodyTrackLbl := createTxt("Melody Tracks:")
	ArtTrackLbl := createTxt("Art Tracks:")
	PPQLbl := createTxt("PPQ:")
	BPMLbl := createTxt("BPM:")

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
	MelodyTrackTXT := createNumberInput(0, 65535)
	ArtTrackTXT := createNumberInput(0, 65535)
	PPQTXT := widget.NewSelect([]string{"96", "192", "240", "480", "960", "1920", "3840", "8192"}, func(string) {})
	BPMTXT := createNumberInput(0, 65535)

	OutputBox := widget.NewMultiLineEntry()
	OutputBox.SetText("Output will go here...")

	outputButton := widget.NewButtonWithIcon("Output Path", theme.FileIcon(), func() {
		logf("Opening output path dialog")

		filePath, err := sqdialog.File().Filter("MIDI Files (.mid)", "mid").Title("Select Output Path").Save()
		if errors.Is(err, sqdialog.ErrCancelled) {
			logf("User cancelled output path dialog")
			return // user cancelled
		} else {
			handleErr(err)
		}

		// append .mid if not present
		if !strings.HasSuffix(filePath, ".mid") {
			filePath += ".mid"
		}

		logf("Output path selected: %s", filePath)
		OutputTXT.SetText(filePath)
	})

	createButton := widget.NewButton("Create", func() {
		var errs []string
		if err := OutputTXT.Validate(); err != nil {
			errs = append(errs, "output: "+err.Error())
		}
		if err := MelodyTrackTXT.Validate(); err != nil {
			errs = append(errs, "melody tracks: "+err.Error())
		}
		if err := ArtTrackTXT.Validate(); err != nil {
			errs = append(errs, "art tracks: "+err.Error())
		}
		if PPQTXT.Selected == "" {
			errs = append(errs, "ppq: cannot be empty")
		}
		if err := BPMTXT.Validate(); err != nil {
			errs = append(errs, "bpm: "+err.Error())
		}

		if len(errs) > 0 {
			dialog.ShowInformation("Invalid Options", strings.Join(errs, "\n"), window)
		} else {
			logf("starting creation | blocked ui")
			var startTime time.Time

			logf("timer enabled, starting timer now")
			startTime = time.Now()

			OutputBox.SetText("")

			melody, err := strconv.Atoi(MelodyTrackTXT.Text)
			handleErr(err)
			art, err := strconv.Atoi(ArtTrackTXT.Text)
			handleErr(err)
			pqq, err := strconv.Atoi(PPQTXT.Selected)
			handleErr(err)
			bpm, err := strconv.Atoi(BPMTXT.Text)
			handleErr(err)

			melodyTrackRangeSTR := a.Preferences().StringWithFallback("melodyTracksRange", "1-15")
			split := strings.Split(melodyTrackRangeSTR, "-")
			min, err := strconv.Atoi(split[0])
			handleErr(err)
			max, err := strconv.Atoi(split[1])
			handleErr(err)
			melodyTrackRange := []int{min, max}

			artTrackRangeSTR := a.Preferences().StringWithFallback("artTracksRange", "16-16")
			split = strings.Split(artTrackRangeSTR, "-")
			min, err = strconv.Atoi(split[0])
			handleErr(err)
			max, err = strconv.Atoi(split[1])
			handleErr(err)
			artTrackRange := []int{min, max}

			drumsEnabled := a.Preferences().BoolWithFallback("allowDrums", false)

			MelodyTrackTXT.Disable()
			ArtTrackTXT.Disable()
			OutputTXT.Disable()
			PPQTXT.Disable()
			BPMTXT.Disable()
			outputButton.Disable()

			window.SetTitle("Empty Track Creator (Running...)")

			var trackCount int
			filePath := OutputTXT.Text

			// if the user chooses a .mid file, we read the track count from that

			// check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				logf("File input does not exist, continuing with new file")

				trackCount = melody + art
			} else {
				logf("File input exists, reading track count from file")

				miditrackcount, err := ReadMIDITracks(filePath, func(format string, a ...any) {
					OutputBox.SetText(OutputBox.Text + fmt.Sprintf(format, a...) + "\n")
				})
				if err != nil {
					logf("Error reading track count from file: %s | unblocking ui", err.Error())

					dialog.ShowError(err, window)

					MelodyTrackTXT.Enable()
					ArtTrackTXT.Enable()
					OutputTXT.Enable()
					PPQTXT.Enable()
					BPMTXT.Enable()
					outputButton.Enable()
					window.SetTitle("Empty Track Creator")
					return
				} else {
					trackCount = miditrackcount + melody + art
				}
			}

			if trackCount > 65535 {
				logf("Track count would be too high (%d > 65535) | unblocking ui", trackCount)

				dialog.ShowError(fmt.Errorf("track count is too high (%d > 65535)", trackCount), window)

				MelodyTrackTXT.Enable()
				ArtTrackTXT.Enable()
				OutputTXT.Enable()
				PPQTXT.Enable()
				BPMTXT.Enable()
				outputButton.Enable()
				window.SetTitle("Empty Track Creator")

				return
			} else {
				logf("creating %v melody + %v art = %v total tracks", melody, art, melody+art)
				tracks := createTracks(melody, art, drumsEnabled, melodyTrackRange, artTrackRange, func(format string, a ...any) {
					OutputBox.SetText(OutputBox.Text + fmt.Sprintf(format, a...) + "\n")
				})
				logf("created tracks")

				logf("writing to %v", filePath)
				WriteMIDI(MIDIInfo{
					tracks:     tracks,
					trackCount: trackCount,
					midiPath:   filePath,
					ppq:        pqq,
					bpm:        bpm,
					allowDrums: drumsEnabled,
					logger: func(format string, a ...any) {
						OutputBox.SetText(OutputBox.Text + fmt.Sprintf(format, a...) + "\n")
					},
					callback: func() {
						MelodyTrackTXT.Enable()
						ArtTrackTXT.Enable()
						OutputTXT.Enable()
						PPQTXT.Enable()
						BPMTXT.Enable()
						outputButton.Enable()
						window.SetTitle("Empty Track Creator")
					},
				})
				logf("wrote to %v | unblocking ui", filePath)

				logf("timer enabled, stopping timer now | took %v", time.Since(startTime))
				OutputBox.SetText(OutputBox.Text + fmt.Sprintf("took %v", time.Since(startTime)) + "\n")
			}
		}
	})

	// set default values
	OutputTXT.SetText(a.Preferences().StringWithFallback("outputPath", "output.mid"))
	MelodyTrackTXT.SetText(a.Preferences().StringWithFallback("melodyTracks", "8"))
	ArtTrackTXT.SetText(a.Preferences().StringWithFallback("artTracks", "8"))
	PPQTXT.SetSelected(a.Preferences().StringWithFallback("ppq", "960"))
	BPMTXT.SetText(a.Preferences().StringWithFallback("bpm", "138"))

	// make rows
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
	bottomRow := container.New(
		layout.NewMaxLayout(),
		OutputBox,
	)

	// put into a column
	content := container.NewBorder(
		container.New(
			layout.NewVBoxLayout(),
			outputRow,
			tracksRow,
			midiRow,
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

		a.Preferences().SetString("outputPath", OutputTXT.Text)
		a.Preferences().SetString("melodyTracks", MelodyTrackTXT.Text)
		a.Preferences().SetString("artTracks", ArtTrackTXT.Text)
		a.Preferences().SetString("ppq", PPQTXT.Selected)
		a.Preferences().SetString("bpm", BPMTXT.Text)

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
