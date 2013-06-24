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
	sayer  *sayer.Sayer
	player *player.Player
}

func NewAlarm(sayer *sayer.Sayer, player *player.Player) *Alarm {

	a := new(Alarm)

	a.sayer = sayer
	a.player = player

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
func (a *Alarm) Start(totalDuration time.Duration, tickBegin time.Duration,
	tickStep time.Duration) {
	var (
		cumulIntvl time.Duration
	)

	log.Printf("Ringing after %v total", totalDuration)
	time.AfterFunc(totalDuration, func() {
		a.Ring()
	})

	/* start ticking after (totalDuration - tickBegin)
	tick after every tickStep
	*/
	log.Printf("Ticking every %v when last %v reached (i.e. after %v)",
		tickStep, tickBegin, totalDuration-tickBegin)
	for cumulIntvl = (totalDuration - tickBegin); cumulIntvl < totalDuration; cumulIntvl += tickStep {
		// wrapped into func so param <timeLeft> 
		// is the correct one everytime
		func(timeLeft time.Duration) {
			// run func after next cumulated interval
			time.AfterFunc(cumulIntvl, func() {
				a.Tick(timeLeft)
			})
		}(totalDuration - cumulIntvl)
	}

}
