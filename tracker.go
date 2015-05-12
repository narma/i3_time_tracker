package main

import (
	//	set "github.com/deckarep/golang-set"
	"log"
	//"strconv"
	"sync"
	"time"
)

type (
	entity struct {
		sync.RWMutex
		entries map[string]time.Duration
		started time.Time
		current string
	}
	Tracker struct {
		workspaces *entity
		windows    *entity
	}
)

func newTrackerEntity(name string) *entity {
	e := &entity{
		started: time.Now(),
		current: name,
		entries: make(map[string]time.Duration),
	}

	// sync every 1 sec
	c := time.Tick(1 * time.Second)
	go func() {
		for _ = range c {
			e.RLock()
			if time.Since(e.started) > time.Second {
				e.RUnlock()
				e.syncWL()
			} else {
				e.RUnlock()
			}

		}
	}()
	return e
}

func NewTracker(workspace string, window string) *Tracker {
	t := &Tracker{}
	t.windows = newTrackerEntity(window)
	t.workspaces = newTrackerEntity(workspace)

	c := time.Tick(60 * time.Second)
	go func() {
		for _ = range c {
			t.windows.RLock()
			log.Println("Windows stat", t.windows.entries)
			t.windows.RUnlock()

			t.workspaces.RLock()
			log.Println("Workspace stat", t.workspaces.entries)
			t.workspaces.RUnlock()
		}
	}()

	return t
}

// exclusive or
func ex_or(a bool, b bool) bool {
	return (a || b) && !(a && b)
}

func (t *Tracker) OnWorkspaceChange(wspace string) {
	if wspace == t.workspaces.current {
		return
	}
	t.workspaces.changed(wspace)
}

func (t *Tracker) OnWindowChange(class string) {
	if class == t.windows.current {
		return
	}
	t.windows.changed(class)
}

func (e *entity) changed(name string) {
	if len(e.current) == 0 {
		return
	}
	e.Lock()
	defer e.Unlock()
	e.sync()
	e.current = name
}

func (t *Tracker) Flush() {

}

func (e *entity) sync() {
	elapsed := time.Since(e.started)
	e.entries[e.current] += elapsed
	e.started = time.Now()
}

func (e *entity) syncWL() {
	e.Lock()
	defer e.Unlock()
	e.sync()
}

func (t *Tracker) Sync() {
	t.workspaces.syncWL()
	t.windows.syncWL()
	if false {
		log.Println("false")
	}
}
