package main

import (
	i3lib "./i3"
	"flag"
	"log"
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
			return
		}
	}
	return
}

func find_focused(root i3lib.TreeNode) (w *i3lib.TreeNode) {
	for _, w := range root.Nodes {
		if w.Focused {
			return &w
		}
	}

	for _, w := range root.Nodes {
		finded := find_focused(w)
		if finded != nil {
			return finded
		}
	}
	return
}

func currentWindowClass(i3 *i3lib.Conn) string {
	root, err := i3.Tree()
	p(err)
	//var w i3lib.TreeNode
	w := find_focused(root)
	if w == nil {
		return ""
	}

	return w.Properties.Class
}

func main() {
	flag.Parse()
	i3, err := i3lib.Attach()
	p(err)

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

	tracker := NewTracker(currentWorkspaceName(i3), currentWindowClass(i3))

	install_subscribes(tracker, i3)

	defer tracker.Flush()
	defer tracker.Sync()
	<-exitChan
	log.Println("bye")
}

func install_subscribes(t *Tracker, i3 *i3lib.Conn) {
	i3.Subscribe("workspace", "window")
	i3.Event.Window = func(ev i3lib.WindowEvent) {
		what := ev.Container.Properties.Class
		if len(what) == 0 {
			what = ev.Container.Properties.Instance
		}
		t.OnWindowChange(what)
	}
	i3.Event.Workspace.Focus = func(current *i3lib.TreeNode, old *i3lib.TreeNode) {
		t.OnWorkspaceChange(current.Name)
	}
}
