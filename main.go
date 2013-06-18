package main

import (
	// "alarm"
	// go get code.google.com/p/gcfg
	"code.google.com/p/gcfg"
	"fmt"
	"os"
	"player"
	// "sayer"
	"flag"
	"strconv"
	"strings"
)

/*	the expected config file key-value structure;
	if not matching, main() will panic
	https://code.google.com/p/gcfg/

	Ex. piradio.ini:

	---

	# Path to streams.list file
	[Streams]
	StreamsList = streams.list

	# Path to (JSON) file containing sounds mappings
	[Sounds]
	SoundsFile = sounds.json

	[Volume]
	VolUpStep = 20
	VolDownStep = 20

	...
	---
*/
type Config struct {
	Streams struct {
		StreamsList string
	}
	Sounds struct {
		SoundsFile string
	}
	Volume struct {
		VolUpStep   int
		VolDownStep int
	}
}

/*	uses an underlying Bash script for acquiring keyboard input
	TODO this is presumably slow
*/
func GetKey(externalCmd string) int64 {
	var (
		proc      *os.Process
		procAttr  os.ProcAttr
		procState *os.ProcessState
		inputKey  int64
		err       error
	)

	inputKey = 0

	procAttr.Files = []*os.File{nil, nil, nil}
	proc, err = os.StartProcess(externalCmd, []string{""}, &procAttr)
	procState, err = proc.Wait()

	/*	we expect information via exit statuses (i.e., error codes)
		from the system script
	 */
	inputKey, err = strconv.ParseInt(strings.TrimPrefix(procState.String(), "exit status "), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	return inputKey
}

func main() {

	var (
		// a           *alarm.Alarm
		// s           *sayer.Sayer
		p        *player.Player
		err      error
		inputKey int64
		conf     Config
		confFile string
		getKeyCmd string
	)

	flag.StringVar(&confFile, "config", "piradio.ini",
		"Configuration file to parse for mandatory and default values")
	flag.Parse()

	// read in config file into struct
	err = gcfg.ReadFileInto(&conf, confFile)
	// if config not as expected, bail
	if err != nil {
		// TODO user feedback first...,
		// then
		panic(err)
	}

	// TODO add to config
	getKeyCmd = "./getkey.sh"

	p = player.NewPlayer(conf.Streams.StreamsList)
	/*	Sayer and Alarm DISABLED FOR NOW
		s = sayer.NewSayer(conf.Sounds.SoundsFile, p)
		s will be used fo Alarm a (see below)
	*/

	for {
		/*	run and wait for completion
			it is ok that this blocks
			because what else should we do in the meantime?
			TODO unfortunately, this seems to block noticeably =(
		*/
		inputKey = GetKey(getKeyCmd)

		switch inputKey {
		case 11:
			go p.VolumeUp(conf.Volume.VolUpStep)
		case 12:
			go p.VolumeDown(conf.Volume.VolDownStep)
		case 1:
			go p.NextStreamByNumber(1)
		case 2:
			go p.NextStreamByNumber(2)
		case 3:
			go p.NextStreamByNumber(3)
		case 4:
			go p.NextStreamByNumber(4)
		case 5:
			go p.NextStreamByNumber(5)
		case 6:
			go p.NextStreamByNumber(6)
		case 7:
			go p.NextStreamByNumber(7)
		case 8:
			go p.NextStreamByNumber(8)
		case 9:
			go p.NextStreamByNumber(9)
		case 13:
			p.Quit()
			os.Exit(0)
		}
	}

}
