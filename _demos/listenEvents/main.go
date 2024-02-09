package main

import (
	"fmt"

	"github.com/anotherhadi/wtui"
	"github.com/anotherhadi/wtui-components/ansi"
	"github.com/anotherhadi/wtui-components/asciitext"
	"github.com/anotherhadi/wtui-components/shortcuts"
	"github.com/anotherhadi/wtui-components/wtopts"
)

func cleanUp() {
	fmt.Print(ansi.AlternativeBufferDisable)
	fmt.Print(ansi.CursorVisible)
}

func main() {
	w := wtui.Init()
	defer w.Quit()
	go w.ListenForEvents()

	fmt.Print(ansi.AlternativeBufferEnable)
	fmt.Print(ansi.CursorInvisible)
	defer cleanUp()

	opts := wtopts.DefaultOpts()
	for range w.EventChannel {
		if w.Key == "q" || w.Key == "Escape" || w.Key == "Ctrl C" {
			break
		}
		if w.Key == "u" {
			w.UpdateCursorPosition()
		}
		fmt.Print(ansi.CursorHome)
		fmt.Print(ansi.ScreenClear)
		opts.MaxCols = w.Width - opts.LeftPadding
		opts.MaxRows = w.Height

		asciitext.Asciitext("Events", opts)

		fmt.Print("\n\n")

		fmt.Println("Focus:", w.Focus)
		fmt.Println("Key:", w.Key)
		fmt.Println("Mod:", w.Mod)
		fmt.Println("Button:", w.Button)
		fmt.Println("ButtonUp:", w.ButtonUp)
		fmt.Println("ButtonDown:", w.ButtonDown)
		fmt.Println("MouseX:", w.MouseX)
		fmt.Println("MouseY:", w.MouseY)
		fmt.Println("CursorX:", w.CursorX)
		fmt.Println("CursorY:", w.CursorY)
		fmt.Println("Width:", w.Width)
		fmt.Println("Height:", w.Height)

		fmt.Print("\n\n")

		shortcuts.Shortcuts([][2]string{
			{"q,ctrl+c,escape", "Quit"},
			{"u", "Update cursor position"},
		})
	}
}
