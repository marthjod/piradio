/*
Package sayer provides simple interrupting sound bite playing methods.

Sayers are used for playing (usually) small sound bites
corresponding to pre-defined messages.
Upon construction, they read in a file mapping message strings
to sound file paths. When prompted to "say" one or more string messages
subsequently, a Sayer will try to play the corresponding sound file
if available.

Sayers interact with a player.Player instance in order to be able to stop
this first, thus ensuring their own sound will be heard. After having
played the snippet, they will restart the respective player.Player.
(Sound snippets are played using VLC media player.)
*/
package sayer

import (
	"log"
	"os"
	"os/exec"
	"player"
	// (opt. before: export https_proxy=https://proxy.example.com:3128)
	// go get github.com/DisposaBoy/JsonConfigReader
	"encoding/json"
	"github.com/DisposaBoy/JsonConfigReader"
	"path/filepath"
	"strconv"
)

type Sayer struct {
	executable string
	quitItem   string
	soundsFile string
	soundsMap  map[string]string
	player     *player.Player
	volume     int
}

/*
NewSayer is the constructor for Sayer objects.
It needs the path to a sounds map file (JSON format) and has to know
which player.Player to stop and start.
It also increases its own volume relative to the player's volume to get
more attention.
*/
func NewSayer(soundsFile string, player *player.Player) *Sayer {
	s := new(Sayer)
	// alias for vlc --intf=dummy
	s.executable = "/usr/bin/cvlc"
	s.quitItem = "vlc://quit"
	s.soundsFile = soundsFile
	s.soundsMap = ReadMap(s.soundsFile)
	s.player = player
	// sayer should be louder than player
	s.volume = player.GetVolume() + 50
	return s
}

/*
Takes a variadic list of messages.
For each message, looks if there is an entry match. If it finds one,
the player is stopped, the appropriate sound file is played and the
player restarted using its original parameters.
*/
func (s *Sayer) Say(messages ...string) {
	var (
		cmd *exec.Cmd
		err error
	)

	for _, sentence := range messages {
		for key, val := range s.soundsMap {
			if key == sentence {
				s.player.Quit()
				log.Printf("%s %s %s", s.executable,
					"--volume="+strconv.Itoa(s.volume), val)
				// say message and force vlc to quit right away
				cmd = exec.Command(s.executable,
					"--volume="+strconv.Itoa(s.volume), val, s.quitItem)
				// Run() waits for cmd to finish
				cmd.Run()
				err = s.player.Resume()
				if err != nil {
					log.Println(err)
				}
				// do not look further
				break
			}
		}
	}
}

/*
Reads in a JSON-formatted file of strings mapped to sound file paths.
The file paths must be relative to the map file to make sure that
valid absolute paths to be played can be generated from the entries.

Ex.
	/path/to/map/sounds.json
contains
	{
		"1m0s":		"sounds/1m0s.mp3",
		"10m0s":	"sounds/10m0s.mp3",
		"problem":	"sounds/problem.mp3"
	}
then entries become prefixed into
	/path/to/map/sounds/1m0s.mp3
etc.

NB: This could and should be an unexported func, but it makes some assumptions,
so I chose it to show up in the documentation.
Whenever possible, treat as unexported, though.
*/
func ReadMap(inFile string) map[string]string {
	var (
		f           *os.File
		err         error
		outMap      map[string]string
		prefixedMap map[string]string
		pathPrefix  string
	)

	if f, err = os.Open(inFile); err == nil {
		defer f.Close()
		// wrap reader before passing it to the json decoder
		r := JsonConfigReader.New(f)
		json.NewDecoder(r).Decode(&outMap)

		pathPrefix = filepath.Dir(inFile) + string(os.PathSeparator)
		prefixedMap = make(map[string]string)
		for entry, shortPath := range outMap {
			prefixedMap[entry] = pathPrefix + shortPath
		}
	} else {
		log.Println(err)
	}

	// log.Printf("%v => %v", inFile, prefixedMap)
	return prefixedMap
}
