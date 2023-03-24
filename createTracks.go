package main

import (
	"bytes"
	"os"
	"time"

	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func createTracks(melodyTracks int, artTracks int, midiPath string, ppq int, bpm int, allowDrums bool, inputPath string, benchmark bool, logger func(format string, a ...any), callback func()) {
	var newMidi bool
	if inputPath == "" {
		newMidi = true
	} else {
		newMidi = false
	}

	logger("creating midi file at: %v", midiPath)
	if newMidi {
		logger("melody tracks: %v | art tracks: %v | ppq: %v | bpm: %v | allow drums: %v", melodyTracks, artTracks, ppq, bpm, allowDrums)
	} else {
		logger("using bpm+ppq from input file | melody tracks: %v | art tracks: %v | allow drums: %v", melodyTracks, artTracks, allowDrums)

	}

	var (
		buffer     bytes.Buffer
		resolution = smf.MetricTicks(ppq)
		firstTrack smf.Track
		midiData   *smf.SMF
		start      time.Time
	)

	if benchmark {
		logger("benchmark enabled, starting timer now")
		start = time.Now()
	}

	if inputPath != "" {
		var err error

		logger("reading from input file: %v", inputPath)
		midiData, err = smf.ReadFile(inputPath)
		handleErr(err)
	} else {
		midiData = smf.New()
	}

	if newMidi {
		midiData.TimeFormat = resolution

		firstTrack.Add(0, smf.MetaTrackSequenceName(""))
		firstTrack.Add(0, smf.MetaTempo(float64(bpm)))
		firstTrack.Close(0)
		midiData.Add(firstTrack)
	}

	for i := 0; i < melodyTracks; i++ {
		var track smf.Track

		j := i % 15
		if !allowDrums && j == 9 {
			melodyTracks = melodyTracks + 1
			continue
		} else {
			logger("adding melody track %v on channel %v", i, j+1)
			if j == 9 {
				track.Add(0, smf.MetaTrackSequenceName("Rhythm"))
			}
			track.Add(0, midi.ProgramChange(uint8(j), gm.Instr_AcousticGrandPiano.Value()))
			track.Close(0)
			err := midiData.Add(track)
			handleErr(err)
		}
	}

	for i := 0; i < artTracks; i++ {
		var track smf.Track

		logger("adding art track")
		track.Add(0, midi.ProgramChange(15, gm.Instr_AcousticGrandPiano.Value()))
		track.Close(0)
		err := midiData.Add(track)
		handleErr(err)
	}

	_, err := midiData.WriteTo(&buffer)
	handleErr(err)

	file, err := os.OpenFile(midiPath, os.O_CREATE|os.O_WRONLY, 0644)
	handleErr(err)

	buffer.WriteTo(file)
	err = file.Close()
	handleErr(err)

	if benchmark {
		elapsed := time.Since(start)
		logger("finished | took %s", elapsed)
	} else {
		logger("finished")
	}

	callback()
}
