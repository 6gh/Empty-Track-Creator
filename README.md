# Empty Track Creator
A simple tool for a simple job

## Purpose

This small project is to quickly make a MIDI file with a specified number of "melody" tracks and "art" tracks for [Black MIDI](https://en.wikipedia.org/wiki/Black_Midi). (Primarliy to avoid spending time on creating the tracks yourself manually)

Melody tracks meaning tracks in which you would put your musical components in (melody, chords, bass, drums, etc). Art tracks meaning tracks in which you would put your Black MIDI Arts in. This is to separate these two channel-wise, so that the art does not interfere with the audio from the melody tracks

## Usage

Download the [latest release](https://github.com/6gh/Empty-Track-Creator/releases/latest) **CLI** exe. The CLI exe is necessary as the GUI exe simply runs the CLI with your options. GUI is optional.

### Using the CLI

Basic usage includes `./empty-track-creator-cli.exe -melody <number> -art <number>`, which creates an output.mid with the tracks you specified

You can also do `.\empty-track-creator-cli.exe -h` for a list of flags you can also use.

### Using the GUI

> If you want to use the GUI, you must have the CLI exe present in the same directory, and the CLI must be called `empty-track-creator-cli.exe`

To use the GUI, simply run the exe and choose your options. Then click "Create" to run the CLI and get the output in the Output box.

The GUI is very simple for the sole reason that it is just for people who might not like the CLI route and want something more user-friendly.
