package main

import (
	"alarm"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"player"
	"sayer"
)

var (
	s           *sayer.Sayer
	p           *player.Player
	a           *alarm.Alarm
	reader      *bufio.Reader
	err         error
	inputKeyBuf [1]byte
	inputCmd    *exec.Cmd
)

func main() {
	reader = bufio.NewReader(os.Stdin)
	s = sayer.NewSayer("sounds.json")
	p = player.NewPlayer("streams.list")

	/*	for getchar()-like behavior cf.
		http://osdir.com/ml/go-language-discuss/2013-03/msg00081.html
	*/
	inputCmd = exec.Command("/bin/stty", "-F", "/dev/tty", "-icanon", "min", "1")
	for {
		// run (and wait for completion)
		inputCmd.Run()
		if _, err = os.Stdin.Read(inputKeyBuf[0:1]); err == nil {
			// TODO compare bytes directly
			switch string(inputKeyBuf[0]) {
			case "0":
				go p.VolumeUp(20)
			case ",":
				go p.VolumeDown(20)
			case "+":
				go p.NextStream()
			// switching through numbers one by one is presumably faster
			case "1":
				go p.NextStreamByNumber(1)
			case "2":
				go p.NextStreamByNumber(2)
			case "3":
				go p.NextStreamByNumber(3)
			case "4":
				go p.NextStreamByNumber(4)
			case "5":
				go p.NextStreamByNumber(5)
			case "6":
				go p.NextStreamByNumber(6)
			case "7":
				go p.NextStreamByNumber(7)
			case "8":
				go p.NextStreamByNumber(8)
			case "9":
				go p.NextStreamByNumber(9)
			//case "alarm\n":
			//	total, tickAfter, tickStep, err := GetAlarmParams()
			//	if err == nil {
			//		a = alarm.NewAlarm(total, tickAfter, tickStep, s, p)
			//	}
			case "-":
				p.Quit()
				os.Exit(0)
			}
		} else {
			fmt.Println(err)
		}
	}
}
