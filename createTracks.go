package main

func createTracks(melody int, art int, allowDrums bool, logger func(format string, a ...any)) []byte {
	var tracksData []byte

	if melody == 0 && art == 0 {
		logf("no tracks to create | dont know how this happened since there are checks in place to prevent this. please report")
		logger("no tracks to create")
		return tracksData
	}

	for i := 0; i < melody; i++ {
		var track []byte

		j := i % 15
		if !allowDrums && j == 9 {
			logf("[m-%v] skipping drum channel as j is %v", i+1, j)
			logger("[M] skipping drum channel")
			melody = melody + 1
			continue
		} else {
			logf("[m-%v] adding track on channel %v", i+1, j)
			logger("[M-%v] adding melody track on channel %v", i+1, j+1)
			createTrack(j, &track)

			tracksData = append(tracksData, track...)
		}
	}

	for i := 0; i < art; i++ {
		var track []byte

		logf("[a-%v] adding art track on channel 16", i+1)
		logger("[A-%v] adding art track", i+1)
		createTrack(15, &track)

		tracksData = append(tracksData, track...)
	}

	return tracksData
}

func createTrack(j int, bytes *[]byte) {
	trackType := []byte{0x4d, 0x54, 0x72, 0x6b} // MTrk
	trackLength := make([]byte, 4)              // size of track
	var trackEvents []byte                      // events in track

	// 0 ticks, ff, 03, 00
	// ff 03 is track name event
	// this sets the track name to nothing
	trackEvents = append(trackEvents, []byte{0x00, 0xff, 0x03, 0x00}...)

	// 0 ticks, cn, pp
	// cn pp is program change event
	// n is channel number, pp is program number
	// this sets the instrument to piano
	// it also sets the channel for the track
	trackEvents = append(trackEvents, []byte{0x00, byte(192 + j), 0x00}...)

	// 0 ticks, ff, 2f, 00
	// ff 2f is end of track event
	trackEvents = append(trackEvents, []byte{0x00, 0xff, 0x2f, 0x00}...)

	// update track length
	trackLength[0] = byte(len(trackEvents) >> 24)
	trackLength[1] = byte(len(trackEvents) >> 16)
	trackLength[2] = byte(len(trackEvents) >> 8)
	trackLength[3] = byte(len(trackEvents))

	trackBytes := append(trackType, trackLength...)
	trackBytes = append(trackBytes, trackEvents...)

	*bytes = append(*bytes, trackBytes...)
}
