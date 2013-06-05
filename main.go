package main

import (
	// "alarm"
	// go get code.google.com/p/gcfg
	"code.google.com/p/gcfg"
	"fmt"
	"os"
	"os/exec"
	"player"
	// "sayer"
	"flag"
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

func main() {

	var (
		// a           *alarm.Alarm
		// s           *sayer.Sayer
		p           *player.Player
		err         error
		inputKeyBuf [1]byte
		inputCmd    *exec.Cmd
		conf        Config
		confFile    string
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

	p = player.NewPlayer(conf.Streams.StreamsList)
	/*	Sayer and Alarm DISABLED FOR NOW
		s = sayer.NewSayer(conf.Sounds.SoundsFile, p)
		s will be used fo Alarm a (see below)
	*/

	/*	simulate getchar()-like behavior, cf.
		http://osdir.com/ml/go-language-discuss/2013-03/msg00081.html
	*/
	inputCmd = exec.Command("/bin/stty", "-F", "/dev/tty", "-icanon", "min", "1")

	for {
		// run (and wait for completion)
		inputCmd.Run()
		if _, err = os.Stdin.Read(inputKeyBuf[0:1]); err == nil {
			switch inputKeyBuf[0] {
			case '0':
				go p.VolumeUp(conf.Volume.VolUpStep)
			case ',':
				go p.VolumeDown(conf.Volume.VolDownStep)
			case '+':
				go p.NextStream()
			// switching through numbers one by one is presumably fastest,
			// compiler- and performance-wise
			case '1':
				go p.NextStreamByNumber(1)
			case '2':
				go p.NextStreamByNumber(2)
			case '3':
				go p.NextStreamByNumber(3)
			case '4':
				go p.NextStreamByNumber(4)
			case '5':
				go p.NextStreamByNumber(5)
			case '6':
				go p.NextStreamByNumber(6)
			case '7':
				go p.NextStreamByNumber(7)
			case '8':
				go p.NextStreamByNumber(8)
			case '9':
				go p.NextStreamByNumber(9)
			//case "alarm\n":
			//	total, tickAfter, tickStep, err := GetAlarmParams()
			//	if err == nil {
			//		a = alarm.NewAlarm(total, tickAfter, tickStep, s, p)
			//	}
			case '-':
				p.Quit()
				os.Exit(0)
			}
		} else {
			fmt.Println(err)
		}
	}
}
