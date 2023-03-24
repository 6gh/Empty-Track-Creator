package main

import (
	"os"
	"time"

	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func createTracks(options Info) {
	log := options.logger
	var newMidi bool
	if options.inputPath == "" {
		newMidi = true
	} else {
		newMidi = false
	}

	log("Creating midi file at: %v", options.midiPath)
	if newMidi {
		log("melody tracks: %v | art tracks: %v | ppq: %v | bpm: %v | allow drums: %v", options.melodyTracks, options.artTracks, options.ppq, options.bpm, options.allowDrums)
	} else {
		log("using bpm+ppq from input file | melody tracks: %v | art tracks: %v | allow drums: %v", options.melodyTracks, options.artTracks, options.allowDrums)
	}

	var (
		resolution = smf.MetricTicks(options.ppq)
		firstTrack smf.Track
		midiData   *smf.SMF
		start      time.Time
	)

	if options.benchmark {
		log("benchmark enabled, starting timer now")
		start = time.Now()
	}

	if options.inputPath != "" {
		var err error

		log("reading from input file: %v", options.inputPath)
		midiData, err = smf.ReadFile(options.inputPath)
		handleErr(err)
	} else {
		midiData = smf.New()
	}

	if newMidi {
		midiData.TimeFormat = resolution

		firstTrack.Add(0, smf.MetaTrackSequenceName(""))
		firstTrack.Add(0, smf.MetaTempo(float64(options.bpm)))
		firstTrack.Close(0)
		midiData.Add(firstTrack)
	}

	for i := 0; i < options.melodyTracks; i++ {
		var track smf.Track

		j := i % 15
		if !options.allowDrums && j == 9 {
			options.melodyTracks = options.melodyTracks + 1
			continue
		} else {
			log("[M-%v] adding melody track on channel %v", i+1, j+1)
			if j == 9 {
				track.Add(0, smf.MetaTrackSequenceName("Rhythm"))
			}
			track.Add(0, midi.ProgramChange(uint8(j), gm.Instr_AcousticGrandPiano.Value()))
			track.Close(0)
			err := midiData.Add(track)
			handleErr(err)
		}
	}

	for i := 0; i < options.artTracks; i++ {
		var track smf.Track

		log("[A-%v] adding art track", i+1)
		track.Add(0, midi.ProgramChange(15, gm.Instr_AcousticGrandPiano.Value()))
		track.Close(0)
		err := midiData.Add(track)
		handleErr(err)
	}

	file, err := os.OpenFile(options.midiPath, os.O_CREATE|os.O_WRONLY, 0644)
	handleErr(err)

	_, err = midiData.WriteTo(file)
	handleErr(err)
	err = file.Close()
	handleErr(err)

	if options.benchmark {
		elapsed := time.Since(start)
		log("finished | took %s", elapsed)
	} else {
		log("finished")
	}

	callback := options.callback
	callback()
}

type Info struct {
	melodyTracks int
	artTracks    int
	midiPath     string
	ppq          int
	bpm          int
	allowDrums   bool
	inputPath    string
	benchmark    bool
	logger       func(format string, a ...any)
	callback     func()
}
