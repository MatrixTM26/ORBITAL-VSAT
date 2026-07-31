package ansi

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Clear() {
	Cmd := exec.Command("clear")
	Cmd.Stdout = os.Stdout
	Cmd.Run()
}

func Typewriter(Text string, Delay time.Duration) {
	for _, Ch := range Text {
		fmt.Printf("%c", Ch)
		time.Sleep(Delay)
	}
	fmt.Println()
}

func TypewriterFast(Text string) {
	Typewriter(Text, time.Millisecond)
}
