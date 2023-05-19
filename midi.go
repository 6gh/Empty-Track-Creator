package main

import (
	"errors"
	"os"
)

func ReadMIDITracks(path string, logger func(format string, a ...any)) (int, error) {
	logf("reading midi path: %v", path)
	logger("reading midi path: " + path)

	// open midi file
	midiFile, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer midiFile.Close()

	// parse header track

	// parse type
	// ensure that type is of MThd
	headerType := make([]byte, 4)
	_, err = midiFile.Read(headerType)
	if err != nil {
		return -1, err
	}

	if string(headerType) != "MThd" {
		logf("invalid header track | header type: %v", string(headerType))
		return -1, errors.New("MIDI file does not contain header track")
	}

	// parse header size
	// ensure that header size is 6
	headerSize := make([]byte, 4)
	_, err = midiFile.Read(headerSize)
	if err != nil {
		return -1, err
	}

	if headerSize[0] != 0 || headerSize[1] != 0 || headerSize[2] != 0 || headerSize[3] != 6 {
		logf("invalid header size (>6) | header type: %v", string(headerType))
		return -1, errors.New("MIDI header size is not 6")
	}

	// parse format
	// ensure that format is 1
	format := make([]byte, 2)
	_, err = midiFile.Read(format)
	if err != nil {
		return -1, err
	}

	if format[0] != 0 || format[1] != 1 {
		logf("invalid midi format | header type: %v", string(headerType))
		return -1, errors.New("MIDI format is not 1")
	}

	// parse track count
	// write to trackCountInt as an integer
	trackCount := make([]byte, 2)
	_, err = midiFile.Read(trackCount)
	if err != nil {
		return -1, err
	}

	trackCountInt := int(trackCount[0])<<8 | int(trackCount[1])

	// we don't need to check the time division as it is not used

	logger("track count: %v", trackCountInt)

	// return track count
	logf("finished reading midi path: %v", path)
	logf("header: %v (%v bytes) with %v tracks", string(headerType), headerSize, trackCountInt)
	return trackCountInt, nil
}

func writePremadeMidi(inputPath string, trackCount int, trackData []byte) error {
	// only modify the track count in the header track

	// open midi file
	midiFile, err := os.OpenFile(inputPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer midiFile.Close()

	// convert track count to bytes
	trackCountBytes := NumberToBytes(trackCount, 2)

	// edit track count to be the new track count
	// write to trackCountInt as an integer
	midiFile.Seek(10, 0) // Seek to track count
	_, err = midiFile.Write(trackCountBytes)
	if err != nil {
		return err
	}

	// write track data
	midiFile.Seek(0, 2) // Seek to end of file
	_, err = midiFile.Write(trackData)
	if err != nil {
		return err
	}

	return nil
}

func writeNewMidi(info MIDIInfo) error {
	// create new midi file
	midiFile, err := os.Create(info.midiPath)
	if err != nil {
		return err
	}

	// write header track
	headerType := []byte("MThd")
	headerSize := []byte{0, 0, 0, 6}
	format := []byte{0, 1}
	trackCount := NumberToBytes(info.trackCount+1, 2) // +1 for conductor track
	timeDivision := GetDeltaTimeBytes(info.ppq)

	header := append(headerType, headerSize...)
	header = append(header, format...)
	header = append(header, trackCount...)
	header = append(header, timeDivision...)

	// write conductor track
	conductorType := []byte("MTrk")
	conductorSize := make([]byte, 4)
	var conductorData []byte

	// write tempo change
	tempoChange := []byte{0x00, 0xFF, 0x51, 0x03}  // delta time, meta event, set tempo, 3 bytes
	tempo := NumberToBytes(60_000_000/info.bpm, 3) // 60_000_000 is the number of microseconds per minute
	tempoChange = append(tempoChange, tempo...)

	// write end of track
	endOfTrack := []byte{0x00, 0xFF, 0x2F, 0x00} // delta time, meta event, end of track, 0 bytes

	// write to data
	conductorData = append(conductorData, tempoChange...)
	conductorData = append(conductorData, endOfTrack...)
	conductorSize = NumberToBytes(len(conductorData), 4)

	conductor := append(conductorType, conductorSize...)
	conductor = append(conductor, conductorData...)

	// write file
	_, err = midiFile.Write(header)
	if err != nil {
		return err
	}
	_, err = midiFile.Write(conductor)
	if err != nil {
		return err
	}
	_, err = midiFile.Write(info.tracks)
	if err != nil {
		return err
	}

	return nil
}

func WriteMIDI(info MIDIInfo) {
	// get the data from input midi file if provided
	logf("writing to midi path: %v", info.midiPath)
	info.logger("writing to midi path: " + info.midiPath)

	if _, err := os.Stat(info.midiPath); os.IsNotExist(err) {
		logf("midi file does not exist, creating new midi file")
		err := writeNewMidi(info)
		if err != nil {
			logf("could not write to midi file, error: %v", err.Error())
			info.logger("error writing to new midi file: " + err.Error())
		}
		logf("wrote new midi to path: %v", info.midiPath)
		info.logger("wrote to new midi file: " + info.midiPath)
	} else {
		// only update the trackcount in the header track
		logf("midi file exists, appending to midi file")
		err := writePremadeMidi(info.midiPath, info.trackCount, info.tracks)
		if err != nil {
			logf("could not save midi file, error: %v", err.Error())
			info.logger("error writing to premade midi file: " + err.Error())
		}
		logf("saved premade midi to path: %v", info.midiPath)
		info.logger("wrote to premade midi file: " + info.midiPath)
	}

	info.callback()
}

type MIDIInfo struct {
	tracks     []byte
	trackCount int
	midiPath   string
	ppq        int
	bpm        int
	allowDrums bool
	logger     func(format string, a ...any)
	callback   func()
}
