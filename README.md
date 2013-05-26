piradio
=======

Go wrappings for network stream (e.g., Internet radio) playing plus timer functions targeted at running on a Raspberry Pi


Description
===========

**piradio** provides convenience wrapper functions for playing and controlling network streams via [VLC media player](www.videolan.org/vlc/).


Internals
=========

Internally, we use _Player_ objects to play a stream, 
_Sayer_ objects to play sound files (which act as alarm tick messages) 
and _Alarm_ objects for handling timers.

Preliminaries/Setup/Howto
=========================



