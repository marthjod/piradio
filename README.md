piradio
=======

[Go](http://golang.org) wrapping for network stream playing (think Internet radio) plus 
timer functions targeted at running on a [Raspberry Pi](http://www.raspberrypi.org/)


Description
===========

**piradio** provides convenience wrapper functions for playing and controlling network streams 
via [VLC media player](http://www.videolan.org/vlc/). Ideally, it should become usable
with acoustic feedback only and controllable via remote control, a numpad, vel sim.

Usage
=====

Commands
--------

As of now, _main.go_ recognizes the following commands from keyboard/numpad-only input:

| Command (entered without newline) | Function |
|:--------|:------------|
| `1` through `9` | If available and not already playing, play stream no. `<number>` from _streams.list_ |
| `<backspace>`    | Quit main executable / abort countdown config |
| `+`   | Increase volume |
| `-` | Decrease volume |
| `<Enter>` | Enter countdown config mode / confirm (start) countdown |
| `0` | Increase countdown time by 10 minutes |
| `000` | Increase countdown time by 1 minute |


Input
------


Because we want to be able to run _piradio_ without a terminal attached (think autostarted daemon), we need an external
helper process providing us with information about key press events ("keylogger").
[key-event.c](https://github.com/marthjod/piradio/blob/master/key-event.c) (called with the correct
input device as argument, e.g. `./key-event /dev/input/event0`) catches key events and writes them to 
a FIFO (named pipe), from which _main.go_ in turn retrieves them. The path to this FIFO must be configured in 
the config file and match the one used by the keylogger binary.


Config file
-----------

The config file serves as an external interface for setting run-time options.
Its format follows the _gitconfig_ or, more specifically, [gcfg](https://code.google.com/p/gcfg/) format
which is based on the .ini-file style.
Until now, it looks like this:

_piradio.ini_:

```ini
# Path to streams.list file
[Streams]
StreamsList = streams.list

# Path to (JSON) file containing sounds mappings
# (used by Alarms)
[Sounds]
SoundsFile = sounds.json

[Volume]
VolUpStep = 20
VolDownStep = 20

# Path to named pipe used for retrieving codes of pressed keys
[IPC]
FifoPath = /tmp/gofifo
```


Internals
=========

Internally, we use _Player_ objects to play a stream, 
_Sayer_ objects to play sound files (which act as alarm tick messages) 
and _Alarm_ objects for handling timers.

Streams
-------

Network stream URLs are read from a streams list file (containing one URL per line) 
and played starting at its first entry.
Switch to another stream by typing the new stream's list number (number keys `1`-`9`).
The streams list consequently has a 9-entry limit (advantage: only one keystroke necessary for switching streams).
Numbers not mapped to an entry are ignored.

_Player_ objects must be initialized with a valid path to a streams list.

Sounds
------

Sound files are needed for acoustic feedback (alarm ringing, alarm tick messages, error messages).
They are best placed in a _sounds/_ subdirectory with corresponding mappings in a JSON-formatted sound map file (e.g., _sounds.json_).

_Sayer_ objects must be initialized with a valid path to a _sounds.json_ and a pre-existing _Player_ object to interrupt if need be.




Alarms
------

Alarms run for a total duration and start ticking after each interval (simulating a countdown) after a certain amount of time 
has passed. For example, you can set an alarm which will ring after 12 minutes and will start ticking
every minute when 5 minutes are left (i.e. after 7 minutes). If it finds a sound file mapped 
to the current tick time,
it will stop the current _Player_, play this sound file (louder than the previous player) and resume the _Player_.



Setup
=====

Prepare Raspberry Pi
--------------------

* Download and set up Raspian as described [here](http://www.raspberrypi.org/downloads)
and [here](http://elinux.org/RPi_Easy_SD_Card_Setup#SD_card_setup) and by running
`raspi-config`
* Get the Internet connection working
* Compile Go for ARM from source as described [here](http://golang.org/doc/install/source)
* Install VLC


Prepare _piradio_
-----------------
* Install package _JsonConfigReader_ from [DisposaBoy](https://github.com/DisposaBoy/JsonConfigReader): `go get github.com/DisposaBoy/JsonConfigReader`
* Install package [_gcfg_](https://code.google.com/p/gcfg): `go get code.google.com/p/gcfg`
* In your _$GOPATH/src_, make subdirectories _player/_, _sayer/_, _alarm/_ (this may vary for different Go setups...) 
* Copy _player.go_, _sayer.go_ and _alarm.go_ to their respective directories under _$GOPATH/src_
* Install packages: 

```bash
for pkg in player sayer alarm; do go install $pkg; done
```
* Populate your _streams.list_, one URL per line
* (Alarms) Acquire some sounds for alarm ticks with, e.g.,

```bash
mplayer -really-quiet -noconsolecontrols \
  -dumpaudio -dumpfile 1m0s.mp3 \
  "http://translate.google.com/translate_tts?tl=en&q=1+minute+left"
```
* (Alarms) Add sound file names and paths to your _sounds.json_, accordingly
* Run _main.go_
	* `go run main.go` or
	* **DOES NOT WORK YET** `go run main.go --config=/path/to/piradio.ini`
* For autostart, put a [SysV Init script](https://github.com/marthjod/piradio/blob/master/piradio-sysv-init) in _/etc/init.d/_ and update runlevel configuration (`sudo update-rc.d piradio defaults`)
* Compile the "keylogger" (`gcc -o key-event key-event.c`) and add the appropriate call to the init script

> NB: the "keylogger" is only intended for a headless Raspberry Pi with no other uses, so this does not invade the user's privacy; also, the 
FIFO gets emptied whenever read, so not even the user's stream choices etc. "touch ground"...


Example hardware setup
-------------

- Raspberry Pi
- AC adapter (2000 mA)
- USB numpad
- USB WiFi dongle
- SD card
- speakers
