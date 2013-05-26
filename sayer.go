package sayer

/*
mkdir $GOPATH/src/pkg/sayer
go build sayer.go
cp sayer.go $GOPATH/src/pkg/sayer
go install sayer
*/

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
}

// constructor
func NewSayer(soundsFile string) *Sayer {
	s := new(Sayer)
	// alias for vlc --intf=dummy
	s.executable = "/usr/bin/cvlc"
	s.quitItem = "vlc://quit"
	s.soundsFile = soundsFile
	s.soundsMap = readMap(s.soundsFile)
	return s
}

/*	Takes variadic list of messages.
	Do not use goroutines here, they'll bump into each other.

	TODO: necessary?:
	~/.config/vlc/vlcrc:
	one-instance-when-started-from-file=0
	to avoid the sayer instance interfering with the player instance
*/
func (s *Sayer) Say(player *player.Player, messages ...string) {
	var (
		cmd *exec.Cmd
		err error
	)
	for _, sentence := range messages {
		for key, val := range s.soundsMap {
			if key == sentence {
				player.Quit()
				log.Printf("%s %s %s", s.executable, "--volume="+strconv.Itoa(player.GetVolume()), val)

				// say message and force vlc to quit right away
				cmd = exec.Command(s.executable, "--volume="+strconv.Itoa(player.GetVolume()), val, s.quitItem)
				// Run() waits for cmd to finish
				cmd.Run()
				err = player.Resume()
				if err != nil {
					log.Println(err)
				}
				// do not look further
				break
			}
		}
	}
}

/* 	sounds.json:
{	
	"30s":	"sounds/30s.mp3",
	"1m0s":	"sounds/1m0s.mp3",
	"1m30s":	"sounds/1m30s.mp3",
	"problem":	"sounds/problem.mp3"
}

	sound files acquired with

	mplayer -really-quiet -noconsolecontrols -dumpaudio -dumpfile $w.mp3 \
		"http://translate.google.com/translate_tts?tl=de&q=$w"
*/
func readMap(inFile string) map[string]string {
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

		/*	prefix with absolute path, i.e.:
			say, /path/to/map/sounds.json has entry "sounds/hey.mp3";
			then the entry gets prefixed and thus becomes
			/path/to/map/sounds/hey.mp3
		*/
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
