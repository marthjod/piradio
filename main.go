package main

import (
	// "alarm"
	// go get code.google.com/p/gcfg
	"code.google.com/p/gcfg"
	"fmt"
	"os"
	// "os/exec"
	"player"
	// "sayer"
	"flag"
	"syscall"
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

	[IPC]
	FifoPath = /tmp/gofifo
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
	IPC struct {
		FifoPath string
	}
}

func main() {

	var (
		// a           *alarm.Alarm
		// s           *sayer.Sayer
		p         *player.Player
		err       error
		input     string
		conf      Config
		confFile  string
		fifo      *os.File
		bytesRead int
		inputKey  []byte
		// keyEventListener *exec.Cmd
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

	// create named pipe (fifo)
	err = syscall.Mkfifo(conf.IPC.FifoPath, syscall.S_IFIFO|0666)
	if err != nil {
		fmt.Println(conf.IPC.FifoPath, err)
	}

	fifo, err = os.Open(conf.IPC.FifoPath)
	if err != nil {
		fmt.Printf("Could not acquire control input from %s, aborting (%s).",
			conf.IPC.FifoPath, err)
		os.Exit(1)
	}

	/*
		// gets started externally as of now
		// because we need root permissions
		keyEventListener = exec.Command("sudo", "./key-event", "/dev/input/event0")
		err = keyEventListener.Start()
		if err != nil {
			fmt.Printf("Could not start key event listener, aborting.")
			os.Exit(1)
		}
	*/

	inputKey = make([]byte, 2)

	for {
		bytesRead, err = fifo.Read(inputKey)

		if err == nil && bytesRead == 2 {
			// ignore null bytes
			// (checking this earlier panicked [?])
			if inputKey[0] != 0 {
				input = string(inputKey)
				fmt.Printf("Read from fifo: bytes %v = string (key) <%v>\n", inputKey, input)

				switch input {
				case "78":
					go p.VolumeUp(conf.Volume.VolUpStep)
				case "74":
					go p.VolumeDown(conf.Volume.VolDownStep)
				case "79":
					go p.NextStreamByNumber(1)
				case "80":
					go p.NextStreamByNumber(2)
				case "81":
					go p.NextStreamByNumber(3)
				case "75":
					go p.NextStreamByNumber(4)
				case "76":
					go p.NextStreamByNumber(5)
				case "77":
					go p.NextStreamByNumber(6)
				case "71":
					go p.NextStreamByNumber(7)
				case "72":
					go p.NextStreamByNumber(8)
				case "73":
					go p.NextStreamByNumber(9)
				case "14":
					p.Quit()
					os.Exit(0)
				}
			}
		} else if err.Error() != "EOF" {
			// "EOF" is expected if no data waiting
			fmt.Println(err)
		}
	}
}
