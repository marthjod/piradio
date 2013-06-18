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
| `1` through `9` | If available and not already playing, play stream no. `<number>` in _streams.list_ |
| `<backspace>`    | Quit main executable |
| `+`   | Increase volume |
| `-` | Decrease volume |

However, keyboard input depends on an underlying [Bash script](https://github.com/marthjod/piradio/blob/master/getkey.sh)
which can always monitor a specific input device.
Because we want to be able to run _piradio_ automatically, without a terminal attached, we need an external
helper process returning information about a key pressed to _main.go_.



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

```

What's missing
==============

Feedback to the user should only be acoustic, thus making it possible to control _piradio_ without having to look at it.

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




Alarms (disabled fttb)
------

Alarms run for a total duration and start ticking after each interval (simulating a countdown) after a certain amount of time 
has passed. For example, you can set an alarm which will ring after 2 minutes 30 seconds and will start ticking
every 10 seconds when 1 minute is left (i.e. after 1 minute 30 seconds). If it finds a sound file mapped 
to the current tick time,
it will stop the current _Player_, play this sound file and resume the _Player_.

### Set up alarm (timer)

**NB: Currently not featured in _main.go_**

Setting up an alarm (user input marked with `#`) and example alarm output:

```
alarm #
Total: 2m30s #
Start ticking at ... left: 1m #
Tick every ...: 10s #
Ringing after 2m30s total
Ticking every 10s when last 1m0s reached (i.e. after 1m30s)
*tick* (time left: 1m0s)
/usr/bin/cvlc --volume=<current player volume> <sounds path>/1m0s.mp3
Resuming player
*tick* (time left: 50s)
*tick* (time left: 40s)
*tick* (time left: 30s)
/usr/bin/cvlc --volume=<current player volume> <sounds path>/30s.mp3
Resuming player
*tick* (time left: 20s)
*tick* (time left: 10s)
Ringing alarm
/usr/bin/cvlc --volume=<current player volume> <sounds path>/alarm.mp3
Resuming player
```


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
go install player
go install sayer
go install alarm
```

* Populate your _streams.list_, one URL per line
* (Alarms) Acquire some sounds for alarm ticks with, e.g.,

```bash
mplayer -really-quiet -noconsolecontrols \
  -dumpaudio -dumpfile 1m30s.mp3 \
  "http://translate.google.com/translate_tts?tl=en&q=1+minute+30+seconds+left"
```

* (Alarms) Add sound file names and paths to your _sounds.json_, accordingly
* Run _main.go_
	* `go run main.go` or
	* `go run main.go --config=/path/to/piradio.ini`

* Find out the device name for the attached input device and map key codes, if needed (see [Bash script](https://github.com/marthjod/piradio/blob/master/getkey.sh))
* For autostart, put a [SysV Init script](https://github.com/marthjod/piradio/blob/master/piradio-sysv-init) in _/etc/init.d/_ and update runlevel configuration (`sudo update-rc.d piradio defaults`)


Example hardware setup
-------------

- Raspberry Pi
- AC adapter (2000 mA)
- USB numpad
- USB WiFi dongle
- SD card
