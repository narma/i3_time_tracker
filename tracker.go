package main

import (
	set "github.com/deckarep/golang-set"
	"github.com/hoisie/redis"
	"log"
	//"strconv"
	"time"
)

type (
	Tracker struct {
		options        *mainOptions
		db             redis.Client
		is_counting    bool
		is_going_relax bool
		started        time.Time
	}

	mainOptions struct {
		WorkWorkspaces set.Set
	}
)

func NewMainOptions() *mainOptions {
	return &mainOptions{
		// Default almost all, except 9, 10
		WorkWorkspaces:      set.NewSetFromSlice([]interface{}{"1", "2", "3", "4", "5", "6", "7", "8"}),
		AllowedRelaxMinutes: 3,
	}
}

func NewTracker(options *mainOptions) *Tracker {
	t := &Tracker{options: options}
	log.Printf("I'm working on %s\n", t.options.WorkWorkspaces)
	c := time.Tick(1 * time.Minute)
	go func() {
		for now := range c {
			t.db.Incrby("tracker:elapsed", 60)
			t.started.Add(60 * time.Second)
		}
	}()
	return t
}

func (t *Tracker) IsHardWork(wspace string) bool {
	return t.options.WorkWorkspaces.Contains(wspace)
}

// exclusive or
func ex_or(a bool, b bool) bool {
	return (a || b) && !(a && b)
}

func (t *Tracker) OnWorkspaceChange(wspace string) {
	state_changed := ex_or(t.IsHardWork(wspace), t.is_counting || !t.is_going_relax)
	if state_changed == false {
		return
	}

	if t.is_counting {
		t.Stop()
	} else {
		t.Start()
	}
}

func (t *Tracker) Stop() {
	t.is_going_relax = time.AfterFunc(t.options.AllowedRelaxMinutes*time.Minute, func() {
		log.Printf("Stop work\n")
		t.is_counting = false
		elapsed := time.Since(t.started)

		t.db.Rpush("tracker:stopped_at", []byte(time.Now().String()))
		t.db.Incrby("tracker:elapsed", int64(elapsed.Seconds()-t.options.AllowedRelaxMinutes*60))
	})

}

func (t *Tracker) Start() {
	log.Printf("Start work\n")
	t.is_counting = true
	t.started = time.Now()
	t.db.Rpush("tracker:started_at", []byte(t.started.String()))
}

func (t *Tracker) AtExit() {
	if t.is_counting {
		t.Stop()
	}
}
