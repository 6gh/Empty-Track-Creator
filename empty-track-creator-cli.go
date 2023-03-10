package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {
	// get options from cli flags
	melodyTracks := flag.Int("melody", 8, "Number of melody tracks")
	artTracks := flag.Int("art", 8, "Number of art tracks")
	midiPath := flag.String("output", "output.mid", "MIDI Output File path. Where the MIDI will be saved")
	ppq := flag.Int("ppq", 960, "Pulses per quarter note")
	bpm := flag.Int("bpm", 138, "BPM. The number of beats per minute")
	allowDrums := flag.Bool("drums", false, "If not present, does not add drum channels (ch 10)")

	flag.Parse()

	log.Println("creating midi file at: ", *midiPath)
	log.Println("melody tracks: ", *melodyTracks, "| art tracks: ", *artTracks, "| ppq: ", *ppq, "| bpm: ", *bpm, "| allow drums: ", *allowDrums)

	var (
		buffer     bytes.Buffer
		resolution = smf.MetricTicks(*ppq)
		firstTrack smf.Track
		midiData   = smf.New()
	)

	midiData.TimeFormat = resolution

	firstTrack.Add(0, smf.MetaTrackSequenceName(""))
	firstTrack.Add(0, smf.MetaTempo(float64(*bpm)))
	firstTrack.Close(0)
	midiData.Add(firstTrack)

	for i := 0; i < *melodyTracks; i++ {
		var track smf.Track

		j := i % 15
		if !*allowDrums && j == 9 {
			*melodyTracks = *melodyTracks + 1
			continue
		} else {
			log.Println("adding melody track", i, "on channel", j+1)
			if j == 9 {
				track.Add(0, smf.MetaTrackSequenceName("Rhythm"))
			}
			track.Add(0, midi.ProgramChange(uint8(j), gm.Instr_AcousticGrandPiano.Value()))
			track.Close(0)
			err := midiData.Add(track)
			handleErr(err)
		}
	}

	for i := 0; i < *artTracks; i++ {
		var track smf.Track

		log.Println("adding art track")
		track.Add(0, midi.ProgramChange(15, gm.Instr_AcousticGrandPiano.Value()))
		track.Close(0)
		err := midiData.Add(track)
		handleErr(err)
	}

	_, err := midiData.WriteTo(&buffer)
	handleErr(err)

	file, err := os.OpenFile(*midiPath, os.O_CREATE|os.O_WRONLY, 0644)
	handleErr(err)

	buffer.WriteTo(file)
	err = file.Close()
	handleErr(err)

	log.Println("finished")
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
