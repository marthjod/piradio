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
	player     *player.Player
	volume     int
}

// constructor
func NewSayer(soundsFile string, player *player.Player) *Sayer {
	s := new(Sayer)
	// alias for vlc --intf=dummy
	s.executable = "/usr/bin/cvlc"
	s.quitItem = "vlc://quit"
	s.soundsFile = soundsFile
	s.soundsMap = readMap(s.soundsFile)
	s.player = player
	// sayer should be louder
	s.volume = player.GetVolume() + 50
	return s
}

/*	Takes variadic list of messages.
	Do not use goroutines here, they'll bump into each other.

	TODO: necessary?:
	~/.config/vlc/vlcrc:
	one-instance-when-started-from-file=0
	to avoid the sayer instance interfering with the player instance
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
				log.Printf("%s %s %s", s.executable, "--volume="+strconv.Itoa(s.volume), val)
				// say message and force vlc to quit right away
				cmd = exec.Command(s.executable, "--volume="+strconv.Itoa(s.volume), val, s.quitItem)
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

/* 	sounds.json:
{	
	"1m0s":	"sounds/1m0s.mp3",
	"10m0s":	"sounds/10m0s.mp3",
	"problem":	"sounds/problem.mp3"
}

	sound files acquired with

	mplayer -really-quiet -noconsolecontrols -dumpaudio -dumpfile $w.mp3 \
		"http://translate.google.com/translate_tts?tl=de&q=$w"

	above json format from directory listing (use only up to "59m0s"):

	find sounds/ -type f -name '*.mp3' |
	perl -ne ' my $line = $_;
		$line =~ s#^sounds/(\d+m0s)\.mp3$#$1#;
		chomp($_);
		chomp($line);
		print "\t\"$line\": \"$_\",\n" '
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
