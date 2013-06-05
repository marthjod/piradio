package alarm

/*
mkdir $GOPATH/src/pkg/alarm
go build alarm.go
cp alarm.go $GOPATH/src/pkg/alarm
go install alarm
*/

import (
	"log"
	"player"
	"sayer"
	"time"
)

type Alarm struct {
	totalDuration time.Duration
	tickBegin     time.Duration
	tickStep      time.Duration
	sayer         *sayer.Sayer
	player        *player.Player
}

func NewAlarm(totalDuration time.Duration, tickBegin time.Duration,
	tickStep time.Duration, sayer *sayer.Sayer, player *player.Player) *Alarm {

	a := new(Alarm)

	// TODO sanity checks if values are reasonable
	a.totalDuration = totalDuration
	a.tickBegin = tickBegin
	a.tickStep = tickStep
	a.sayer = sayer
	a.player = player

	// start alarm right away
	a.Start()

	return a
}

//TODO see if "a *Alarm" is (more) correct here
func (a Alarm) Ring() {
	log.Println("Ringing alarm")
	a.sayer.Say("alarm")
}

func (a *Alarm) Tick(timeLeft time.Duration) {
	log.Printf("*tick* (time left: %v)", timeLeft.String())
	// if string representation of time left
	// is found in sounds map, say it
	a.sayer.Say(timeLeft.String())
}

/* 	interval grows,
time timeLeft shrinks
*/
func (a *Alarm) Start() {
	var (
		cumulIntvl time.Duration
	)

	log.Printf("Ringing after %v total", a.totalDuration)
	time.AfterFunc(a.totalDuration, func() {
		a.Ring()
	})

	/* start ticking after (totalDuration - tickBegin)
	tick after every tickStep
	*/
	log.Printf("Ticking every %v when last %v reached (i.e. after %v)",
		a.tickStep, a.tickBegin, a.totalDuration-a.tickBegin)
	for cumulIntvl = (a.totalDuration - a.tickBegin); cumulIntvl < a.totalDuration; cumulIntvl += a.tickStep {
		// wrapped into func so param <timeLeft> 
		// is the correct one everytime
		func(timeLeft time.Duration) {
			// run func after next cumulated interval
			time.AfterFunc(cumulIntvl, func() {
				a.Tick(timeLeft)
			})
		}(a.totalDuration - cumulIntvl)
	}

}
