package wtui

import (
	"fmt"
	"sync"
)

type Window struct {
	CursorX    int
	CursorY    int
	MouseX     int
	MouseY     int
	Width      int
	Height     int
	Focus      bool
	Key        string
	Button     string
	ButtonUp   bool
	ButtonDown bool
	Mod        string
	IsClosed   bool

	EventChannel chan EventType
	quit         chan bool
	wait         *sync.WaitGroup
	hexcode      *map[string]string
}

func Init() Window {
	var window Window
	window.quit = make(chan bool)
	window.EventChannel = make(chan EventType)
	var wait sync.WaitGroup
	window.wait = &wait
	window.Focus = true
	window.Width = 100
	window.Height = 100
	hexCode := make(map[string]string)
	hexCode["\033"] = "Escape"
	hexCode["\u000A"] = "Enter"
	hexCode["\u007f"] = "Backspace"
	hexCode["\u0000"] = "Ctrl @"
	hexCode["\u0001"] = "Ctrl A"
	hexCode["\u0002"] = "Ctrl B"
	hexCode["\u0003"] = "Ctrl C"
	hexCode["\u0004"] = "Ctrl D"
	hexCode["\u0005"] = "Ctrl E"
	hexCode["\u0006"] = "Ctrl F"
	hexCode["\u0007"] = "Ctrl G"
	hexCode["\u0008"] = "Ctrl H"
	hexCode["\u0009"] = "Ctrl I"
	hexCode["\u000B"] = "Ctrl K"
	hexCode["\u000C"] = "Ctrl L"
	hexCode["\u000D"] = "Ctrl M"
	hexCode["\u000E"] = "Ctrl N"
	hexCode["\u000F"] = "Ctrl O"
	hexCode["\u0010"] = "Ctrl P"
	hexCode["\u0011"] = "Ctrl Q"
	hexCode["\u0012"] = "Ctrl R"
	hexCode["\u0013"] = "Ctrl S"
	hexCode["\u0014"] = "Ctrl T"
	hexCode["\u0015"] = "Ctrl U"
	hexCode["\u0016"] = "Ctrl V"
	hexCode["\u0017"] = "Ctrl W"
	hexCode["\u0018"] = "Ctrl X"
	hexCode["\u0019"] = "Ctrl Y"
	hexCode["\u001A"] = "Ctrl Z"
	hexCode["\u001B"] = "Ctrl ["
	hexCode["\u001C"] = "Ctrl \\"
	hexCode["\u001D"] = "Ctrl ]"
	hexCode["\u001E"] = "Ctrl ^"
	hexCode["\u001F"] = "Ctrl _"
	window.hexcode = &hexCode

	return window
}

func (w *Window) Quit() {
	w.quit <- true
	w.IsClosed = true
	w.wait.Wait()
}

func (w *Window) resetState() {
	w.Key = ""
	w.Mod = ""
	w.Button = ""
	w.ButtonDown = false
	w.ButtonUp = false
}

func (w *Window) ListenForEvents() {
	w.wait.Add(1)
	state, _ := getState()
	fmt.Print("\033[?1003h\033[?1015h\033[?1006h")
	fmt.Print("\033[?1004h")
	go w.getEvents()
	go w.getResizeEvents()
	go func(wg *sync.WaitGroup) {
		for range w.quit {
			_ = restoreMode(state)
			fmt.Print("\033[?1003l\033[?1015l\033[?1006l")
			fmt.Print("\033[?1004l")
			wg.Done()
		}
	}(w.wait)
	w.wait.Wait()
}

func (w *Window) WaitForEvents() EventType {
	for event := range w.EventChannel {
		return event
	}
	return -1
}
