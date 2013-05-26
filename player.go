package player

/*
mkdir $GOPATH/src/pkg/player
go build player.go
cp player.go $GOPATH/src/pkg/player
go install player
*/

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	executable      string
	cmd             *exec.Cmd
	socketFile      string
	pid             int
	volume          int
	conn            net.Conn
	streamsFile     string
	streamsList     []string
	currentStreamNo int
}

/*	Player constructor
	Builds the player command and starts it.
*/
func NewPlayer(streamsFile string) *Player {
	p := new(Player)

	p.volume = 200
	// do not use cvlc b/c it would use only the dummy interface
	p.executable = "/usr/bin/vlc"
	p.streamsFile = streamsFile
	p.currentStreamNo = 0

	p.streamsList = readList(p.streamsFile)
	if len(p.streamsList) < 1 {
		log.Fatal("No streams loaded!")
	}

	// start playing right away
	p.Start()

	return p
}

func (p *Player) Start() {
	var (
		err error
	)

	// each socket file should have a different name
	rand.Seed(time.Now().UTC().UnixNano())
	p.socketFile = "/tmp/vlc" + strconv.Itoa(rand.Int()) + ".sock"

	// build command
	p.cmd = exec.Command(p.executable, "--intf=oldrc",
		"--rc-unix="+p.socketFile,
		"--rc-fake-tty",
		"--volume="+strconv.Itoa(p.volume),
		//"--http-proxy=proxy.example.com:3128",
		p.streamsList[p.currentStreamNo])

	if err = p.cmd.Start(); err == nil {
		p.pid = p.cmd.Process.Pid
		log.Printf("Started %v (PID %v) %v", p.cmd.Path, p.pid, p.cmd.Args)
		p.connectToSocket()
	} else {
		log.Println(err)
	}
}

func (p *Player) NextStream() {
	p.Quit()
	// reset to first stream if stream index out of range
	if p.currentStreamNo++; p.currentStreamNo > len(p.streamsList)-1 {
		p.currentStreamNo = 0
	}
	p.Start()
}

func (p *Player) GetVolume() int {
	return p.volume
}

// lowercase: private/unexported
func (p *Player) sendToSocket(msg string) error {
	var (
		err error
	)

	if p.conn != nil {
		_, err = p.conn.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		} else {
			// msg = strings.Replace(msg, "\n", "\\n", -1)
			// log.Printf("'%v' -> %v", msg, p.socketFile)
		}
	} else {
		log.Println("No socket connection!")
		err = errors.New("No socket connection!")
	}

	return err
}

// lowercase: private/unexported
func (p *Player) kill() error {
	var (
		proc *os.Process
		err  error
	)

	if proc, err = os.FindProcess(p.pid); err == nil {
		if err = proc.Kill(); err == nil {
			// log.Printf("Killed process %v", p.pid)
		} else {
			log.Println("Could not kill process", err)
		}
	} else {
		log.Println("Could not find process", err)
	}

	return err
}

func (p *Player) SetVolume(vol int) error {
	var (
		err error
	)

	if err = p.sendToSocket("volume " + strconv.Itoa(vol) + "\n"); err == nil {
		p.volume = vol
	} else {
		log.Println(err)
	}

	return err
}

func (p *Player) Quit() error {
	var (
		err  error
		proc *os.Process
	)

	err = p.sendToSocket("quit\n")
	if err != nil {
		log.Println(err)
	} else {
		// log.Println("Sent quit request to VLC socket")
	}
	proc, err = os.FindProcess(p.pid)
	if err != nil || proc == nil {
		// kill if sending "quit\n" did not suffice
		if err = p.kill(); err != nil {
			log.Println(err)
		} else {
			// log.Println("Killed")
		}
	}
	return err
}

// VLC's "volup" and "voldown" do not work as expected
func (p *Player) VolumeUp(step int) error {
	var (
		err error
	)

	p.volume += step
	err = p.SetVolume(p.volume)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (p *Player) VolumeDown(step int) error {
	var (
		err error
	)

	p.volume -= step
	err = p.SetVolume(p.volume)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (p *Player) Resume() error {
	var (
		err error
	)

	log.Printf("Resuming player")
	p.Start()
	// p.volume should still hold the previous value
	err = p.SetVolume(p.volume)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (p *Player) connectToSocket() error {
	var (
		i            int
		max_attempts int
		err          error
	)

	max_attempts = 20

	for i = 1; i <= max_attempts; i++ {
		// log.Printf("Trying to connect to socket %v (attempt #%v)", p.socketFile, i)
		if p.conn, err = net.Dial("unix", p.socketFile); err == nil {
			log.Printf("Connected to %v", p.socketFile)
			break
		}
		// check again after short while
		time.Sleep(300 * time.Millisecond)
	}

	if err != nil && i == max_attempts {
		log.Println("No socket connection established", err)
	}

	return err
}

func readList(inFile string) []string {
	var (
		err     error
		strbuf  string
		fdata   []byte
		outList []string
	)

	if fdata, err = ioutil.ReadFile(inFile); err == nil {
		strbuf = string(fdata)
		outList = strings.Split(strbuf, "\n")
		// kill possible last empty element 
		// caused by file's trailing newline
		if outList[len(outList)-1] == "" {
			outList = outList[:len(outList)-1]
		}
	} else {
		log.Println(err)
	}

	// log.Printf("%v => %v", inFile, outList)
	return outList
}
