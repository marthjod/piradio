package main

import (
	"alarm"
	"bufio"
	"fmt"
	"log"
	"os"
	"player"
	"sayer"
	"strings"
	"time"
)

var (
	s      *sayer.Sayer
	p      *player.Player
	a      *alarm.Alarm
	reader *bufio.Reader
	err    error
	line   []byte
)

func ReadStdin(msg string) (string, error) {
	var (
		stdin *bufio.Reader
	)

	stdin = bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	line, err = stdin.ReadBytes('\n')

	return strings.TrimRight(string(line), "\n"), err

}

func GetTimeParam(msg string) (time.Duration, error) {
	var (
		duration time.Duration
		buf      string
		err      error
	)

	err = nil
	duration = time.Since(time.Now())

	if buf, err = ReadStdin(msg); err == nil {
		duration, err = time.ParseDuration(buf)
	}

	return duration, err
}

func GetAlarmParams() (time.Duration, time.Duration, time.Duration, error) {
	var (
		total     time.Duration
		tickBegin time.Duration
		tickStep  time.Duration
		err       error
	)

	total, err = GetTimeParam("Total: ")
	if err != nil {
		log.Println(err)
	} else {

		tickBegin, err = GetTimeParam("Start ticking at ... left: ")
		if err != nil {
			log.Println(err)
		} else {

			tickStep, err = GetTimeParam("Tick every ...: ")
			if err != nil {
				log.Println(err)
			}
		}
	}

	if err == nil {
		log.Printf("Alarm params: total = %v, begin ticking at %v left, tick every %v", total.String(), tickBegin.String(), tickStep.String())
	}

	return total, tickBegin, tickStep, err
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	s = sayer.NewSayer("sounds.json")
	p = player.NewPlayer("streams.list")
	//a = alarm.NewAlarm(2*time.Minute, 1*time.Minute+40*time.Second, 5*time.Second, s, p)
	for {
		if line, err = reader.ReadBytes('\n'); err == nil {
			// line is []byte
			switch string(line) {
			case "volup\n":
				go p.VolumeUp(100)
			case "voldown\n":
				go p.VolumeDown(100)
			case "next\n":
				go p.NextStream()
			// switching through numbers one by one is presumably faster
			case "1\n":
				go p.NextStreamByNumber(1)
			case "2\n":
				go p.NextStreamByNumber(2)
			case "3\n":
				go p.NextStreamByNumber(3)
			case "4\n":
				go p.NextStreamByNumber(4)
			case "5\n":
				go p.NextStreamByNumber(5)
			case "6\n":
				go p.NextStreamByNumber(6)
			case "7\n":
				go p.NextStreamByNumber(7)
			case "8\n":
				go p.NextStreamByNumber(8)
			case "9\n":
				go p.NextStreamByNumber(9)
			case "alarm\n":
				total, tickAfter, tickStep, err := GetAlarmParams()
				if err == nil {
					a = alarm.NewAlarm(total, tickAfter, tickStep, s, p)
				}
			case "quit\n":
				p.Quit()
				os.Exit(0)
			}
		} else {
			fmt.Println(err)
		}
	}
}
