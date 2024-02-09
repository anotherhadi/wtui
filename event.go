package wtui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

type EventType int8

const (
	KeyEvent           EventType = 0
	MouseMovementEvent EventType = 1
	ButtonEvent        EventType = 2
	ResizeEvent        EventType = 3
	FocusEvent         EventType = 4
	CursorEvent        EventType = 5
	RefreshEvent       EventType = 6
)

func getSize(fd int) (int, int, error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}

func (w *Window) UpdateCursorPosition() {
	fmt.Print("\033[6n")
}

func (w *Window) getResizeEvents() {
	fd := int(os.Stdin.Fd())
	w.Width, w.Height, _ = getSize(fd)
	w.EventChannel <- ResizeEvent
	ticker := time.Tick(time.Millisecond * 300)
	for range ticker {
		wi, he, _ := getSize(fd)
		if wi != w.Width || he != w.Height {
			w.Width = wi
			w.Height = he
			w.EventChannel <- ResizeEvent
		}
	}
}

func (w *Window) getEvents() {
	in := bufio.NewReader(os.Stdin)

	for {
		key := getch(in)
		w.resetState()
		if strings.HasPrefix(key, "\033[") {
			eventType := w.parseEvent(key)
			if eventType != -1 {
				w.EventChannel <- eventType
			}
		} else {
			w.EventChannel <- w.parseKey(key)
		}
	}
}

func (w *Window) parseKey(input string) EventType {

	key, exist := (*w.hexcode)[input]
	if exist {
		w.Key = key
	} else {
		w.Key = input
	}
	return KeyEvent
}

func (w *Window) parseButtonOrMouseEvent(input string) EventType {
	list := strings.Split(input, ";")
	ev := list[0]
	x := list[1]
	y := list[2][:len(list[2])-1]
	kind := list[2][len(list[2])-1]

	w.MouseX, _ = strconv.Atoi(x)
	w.MouseY, _ = strconv.Atoi(y)

	var eventType EventType

	switch ev {
	case "35":
		return MouseMovementEvent
	case "0", "8", "16", "32", "40", "48", "24", "56":
		eventType = ButtonEvent
		w.Button = "Button"
	case "1", "9", "17", "33", "41", "42", "25", "57":
		eventType = ButtonEvent
		w.Button = "Button3"
	case "2", "10", "18", "34", "49", "50", "26", "58":
		eventType = ButtonEvent
		w.Button = "Button2"
	}

	switch ev {
	case "0", "1", "2":
		w.Mod = ""
	case "8", "9", "10":
		w.Mod = "Alt"
	case "16", "17", "18":
		w.Mod = "Ctrl"
	case "24", "25", "26":
		w.Mod = "Alt Ctrl"
	}

	switch ev {
	case "32", "33", "34":
		w.Mod = "Drag"
	case "40", "41", "42":
		w.Mod = "Drag Alt"
	case "48", "49", "50":
		w.Mod = "Drag Ctrl"
	case "56", "57", "58":
		w.Mod = "Drag Alt Ctrl"
	}

	switch kind {
	case 'm':
		w.ButtonUp = true
	case 'M':
		w.ButtonDown = true
	}
	return eventType
}

func (w *Window) parseEvent(input string) EventType {
	input = regexp.MustCompile(`^[^a-zA-Z]*[a-zA-Z]`).FindString(input)
	input = strings.TrimPrefix(input, "\033[")

	if strings.HasPrefix(input, "<") { // Button Or Mouse Movement
		return w.parseButtonOrMouseEvent(strings.TrimPrefix(input, "<"))
	} else if strings.HasPrefix(input, "I") {
		w.Focus = true
		return FocusEvent
	} else if strings.HasPrefix(input, "O") {
		w.Focus = false
		return FocusEvent
	} else if strings.HasPrefix(input, "A") {
		w.Key = "Up"
		return KeyEvent
	} else if strings.HasPrefix(input, "D") {
		w.Key = "Left"
		return KeyEvent
	} else if strings.HasPrefix(input, "B") {
		w.Key = "Down"
		return KeyEvent
	} else if strings.HasPrefix(input, "C") {
		w.Key = "Right"
		return KeyEvent
	} else if strings.HasSuffix(input, "R") {
		input = strings.TrimSuffix(input, "R")
		list := strings.Split(input, ";")
		x, _ := strconv.Atoi(list[0])
		y, _ := strconv.Atoi(list[1])
		if w.CursorX != x || w.CursorY != y {
			w.CursorX = x
			w.CursorY = y
			return CursorEvent
		}
	}

	return -1
}
