# Empty Track Creator
A simple tool for a simple job.

## Purpose

This small project is to quickly make a MIDI file with a specified number of "melody" tracks and "art" tracks for [Black MIDI](https://en.wikipedia.org/wiki/Black_Midi). (Primarliy to avoid spending time on creating the tracks yourself manually)

Melody tracks meaning tracks in which you would put your musical components in (melody, chords, bass, drums, etc). Art tracks meaning tracks in which you would put your Black MIDI Arts in. This is to separate these two channel-wise, so that the art does not interfere with the audio from the melody tracks

## Usage

Download the [latest release](https://github.com/6gh/Empty-Track-Creator/releases/latest). Currently, the only built release is for windows. This is due to me not having a Linux or Mac machine, so I am not able to verify that it works on these OSes.

## Building 

You will need to install the packages required using Go and also follow [Fyne getting started guide](https://developer.fyne.io/started/) to install and use fyne (gui framework). After that just use `fyne package` and you will get your executable.

## License

[MIT](https://github.com/6gh/Empty-Track-Creator/blob/master/LICENSE)
