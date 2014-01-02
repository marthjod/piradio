piradio
=======

[Go](http://golang.org) wrapping for network stream playing (think Internet radio) plus
timer functions targeted at running on a [Raspberry Pi](http://www.raspberrypi.org/)

![](https://github.com/marthjod/piradio/blob/master/piradio.gif?raw=true)


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

| Command (key) | Function |
|:--------|:------------|
| `1` through `9` | If available and not already playing, play stream no. `<number>` from _streams.list_ |
| `<Backspace>`    | Quit main executable / abort countdown config and reset countdown time^1 |
| `+`   | Increase volume |
| `-` | Decrease volume |
| `<Enter>` | Enter countdown config mode / confirm (start) countdown^1 |
| `0` | Increase countdown time by 10 minutes |
| `,` | Increase countdown time by 1 minute (`0`and `000` are the same on some numpads, unfortunately) |

^1 when in countdown config mode

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
Its format follows the [gcfg](https://code.google.com/p/gcfg/) format
based on .ini-file style.
Until now, it must look like this:

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
They are best placed in a _sounds/_ subdirectory with corresponding mappings in a JSON-formatted sound map file.

Example _sounds.json_:

```json
{
	"1m0s":		"sounds/1m0s.mp3",
	"10m0s":	"sounds/10m0s.mp3",
	"problem":	"sounds/problem.mp3"
}
```

NB: All parts concerned with time durations (_Alarms_ and _Sayers_ called by _Alarm_ functions)
only recognize full-minute times up to _59m0s_ (format is that of Go's
[time.Duration string representation](http://golang.org/pkg/time/#Duration.String)) as of now.

Generate the map file's chunk (given you have already sound clips in a _sounds/_ subdirectory, see below)
with something like

```bash
	find sounds/ -type f -name '*.mp3' |
	perl -ne ' my $line = $_;
		$line =~ s#^sounds/(\d+m0s)\.mp3$#$1#;
		chomp($_);
		chomp($line);
		print "\t\"$line\": \"$_\",\n" '
```

_Sayer_ objects must be initialized with a valid path to a sound map file and a pre-existing _Player_ object to interrupt.


Alarms
------

Alarms run for a total duration and start ticking after each interval (simulating a countdown) after a certain amount of time
has passed. For example, you can set an alarm which will ring after 12 minutes and will start ticking
every minute when 5 minutes are left (i.e. after 7 minutes). If it finds a sound file mapped
to the current tick time,
it will stop the current _Player_, play this sound file (louder than the previous player) and resume the _Player_.

TODOs/Known bugs
------

- Setting a countdown for intervals < 5 min causes freezing (negative intervals etc.). Resort to intervals > 10 min for now.

Setup
=====

Prepare Raspberry Pi
--------------------

* Download and set up Raspbian as described [here](http://www.raspberrypi.org/downloads)
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
* Build the packages to check if they compile

```bash
for pkg in player sayer alarm; do go build $pkg; done
```

* Copy _player.go_, _sayer.go_ and _alarm.go_ to their respective directories under _$GOPATH/src_
* Install packages

```bash
for pkg in player sayer alarm; do go install $pkg; done
```

* Populate your _streams.list_, one URL per line (no limit, but only first 9 will be
accessible via keys)
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
FIFO gets emptied whenever read, so not even the user's stream choices etc. really "touch ground" (unless you log to a log file).



Example hardware setup
-------------

- Raspberry Pi
- AC adapter (2000 mA)
- USB numpad
- USB WiFi dongle
- SD card
- speakers


Code
-----

| Source file | Godoc'd? |
|:--------|:------------|
| sayer.go | [x] |
| player.go    | [ ] |
| alarm.go  | [ ] |
| main.go  |  [ ]  |

- Godoc: `godoc -http=:6000`
