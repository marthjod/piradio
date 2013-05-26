piradio
=======

Go wrappings for network stream playing (think Internet radio) plus 
timer functions targeted at running on a [Raspberry Pi](http://www.raspberrypi.org/)


Description
===========

**piradio** provides convenience wrapper functions for playing and controlling network streams 
via [VLC media player](http://www.videolan.org/vlc/).

As of now, it recognizes the following commands from keyboard input:

| Command (enter with newline) | Function |
|:--------|:------------|
| `next`    | Quit currently playing stream and start next one         |
| `quit`    | Quit main executable |
| `volup`   | Increase volume by 100 |
| `voldown` | Decrease volume by 100 |
| `alarm`   | Set up alarm |


Setting up an alarm (user input marked with `#`) and alarm output:

```bash
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

What's missing
==============

In the future, it should be possible to acquire input by more comfortable means than typing commands on  
a keyboard; for example, via remote control. 


Internals
=========

Internally, we use _Player_ objects to play a stream, 
_Sayer_ objects to play sound files (which act as alarm tick messages) 
and _Alarm_ objects for handling timers.

Streams
-------

Streams are read from _streams.list_ (containing one network stream URL per line) 
and played in the file's order. If the last stream is reached, the next stream
will be the first one in the list again.


Alarms
------

Alarms run for a total duration and start ticking after each interval after a certain amount of time 
has passed. For example, you can set an alarm which will ring after 2 minutes 30 seconds and will start ticking
every 10 seconds when 1 minute is left (i.e. after 1 minute 30 seconds). If it finds a sound file mapped 
to the current tick time,
it will stop the current Player, play this sound file and resume the Player.


Preliminaries/Setup/Howto
=========================

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

