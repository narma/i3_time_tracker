package main

import (
	i3lib "./i3"
	"flag"
	//"log"
	"os"
	"os/signal"
	"syscall"
)

func p(e error) {
	if e != nil {
		panic(e)
	}
}

func currentWorkspaceName(i3 *i3lib.Conn) (current_workspace string) {
	ws, err := i3.Workspaces()
	p(err)

	for _, v := range ws {
		if v.Focused {
			current_workspace = v.Name
		}
	}
	return
}

func main() {
	flag.Parse()
	i3, err := i3lib.Attach()
	if err != nil {
		panic(err)
	}

	go func() {
		panic(i3.Listen())
	}()

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := NewMainOptions()

	args := flag.Args()
	if len(args) > 0 {
		options.WorkWorkspaces.Clear()
		for _, v := range flag.Args() {
			options.WorkWorkspaces.Add(v)
		}
	}

	tracker := NewTracker(options)

	tracker.OnWorkspaceChange(currentWorkspaceName(i3))
	install_subscribes(tracker, i3)
	defer tracker.AtExit()
	<-exitChan
}

func install_subscribes(t *Tracker, i3 *i3lib.Conn) {
	i3.Subscribe("workspace")
	i3.Event.Workspace.Focus = func(current *i3lib.TreeNode, old *i3lib.TreeNode) {
		t.OnWorkspaceChange(current.Name)
	}
}
