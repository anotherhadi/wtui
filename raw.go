package wtui

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
	"strings"
)

func getch(in *bufio.Reader) string {
	state, _ := setRawMode()
	var input rune
	input, _, _ = in.ReadRune()
	str := strings.Builder{}
	if input == '\033' {
		str.WriteRune(input)
		for in.Buffered() > 0 {
			input, _, _ = in.ReadRune()
			str.WriteRune(input)
		}
	} else {
		str.WriteRune(input)
	}
	_ = restoreMode(state)

	return str.String()
}

func getState() (termios *unix.Termios, err error) {
	fd := int(os.Stdin.Fd())
	termios, err = unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}
	return termios, nil
}

func setRawMode() (termios *unix.Termios, err error) {
	fd := int(os.Stdin.Fd())
	termios, err = unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	rawTermios := termios
	rawTermios.Lflag &^= unix.ICANON | unix.ECHO | unix.ISIG
	rawTermios.Cc[unix.VMIN] = 1
	rawTermios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, rawTermios); err != nil {
		return nil, err
	}

	return termios, nil
}

func restoreMode(termios *unix.Termios) error {
	fd := int(os.Stdin.Fd())
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return err
	}

	return nil
}
