package generate

import (
	"fmt"
	"io"
	"runtime"

	"os"

	"github.com/fatih/color"
)

func printBanner() {
	var writer io.Writer

	writer = os.Stdout
	if runtime.GOOS == "windows" {
		writer = color.Output
	}
	w1 := (color.New(color.BgWhite).SprintFunc())(" ")
	w2 := (color.New(color.BgHiWhite).SprintFunc())(" ")
	bl := (color.New(color.BgBlack).SprintFunc())(" ")
	c1 := (color.New(color.BgCyan).SprintFunc())(" ")
	c2 := (color.New(color.BgHiCyan).SprintFunc())(" ")
	y1 := (color.New(color.BgYellow).SprintFunc())(" ")
	y2 := (color.New(color.BgHiYellow).SprintFunc())(" ")
	r1 := (color.New(color.BgRed).SprintFunc())(" ")
	r2 := (color.New(color.BgHiRed).SprintFunc())(" ")
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c2+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+c2+c2+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+c2+c2+y1+y2+c1+c1+bl+bl+c2+c2+c2+c2+c2+c2+c2+c2+c2+c2+c2+c2+bl+bl+c2+c2+y1+y2+c1+c1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+w2+w2+w2+w2+w2+w2+c1+c1+c1+c1+c1+c1+c1+c1+w2+w2+w2+w2+w2+w2+bl+bl+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+w2+w2+w2+w2+w2+w2+w2+w2+w1+w1+c1+c1+c1+c1+w2+w2+w2+w2+w2+w2+w2+w2+w1+w1+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+w2+w2+w2+w2+bl+bl+w2+w2+w1+w1+c1+c1+c1+c1+w2+w2+w2+w2+bl+bl+w2+w2+w1+w1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+bl+bl+r2+r2+r2+r2+r2+r1+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+w1+w1+w2+w2+w2+w2+w2+w2+w1+w1+c1+c1+c1+c1+w1+w1+w2+w2+w2+w2+w2+w2+w1+w1+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r2+r2+r1+bl+bl+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+w1+w1+w1+w1+w1+w1+c1+c1+bl+bl+bl+bl+c1+c1+w1+w1+w1+w1+w1+w1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+y2+y2+y2+y2+y2+y2+y2+y1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+y2+y2+y2+y2+y2+y2+y2+y1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+w2+w2+w2+w1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+w2+w2+w2+w1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+y2+y1+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+y2+y1+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+y2+y2+y2+y1+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+y2+y2+y2+y1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r1+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r1+bl+"\n")
	fmt.Fprintf(writer, bl+bl+y2+y2+y2+y1+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+y2+y2+y2+y1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r2+r1+bl+bl+r2+r2+r2+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+y2+y2+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+y2+y1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+r2+r2+r2+r1+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r1+r1+bl+bl+bl+bl+bl+bl+r2+r1+bl+bl+r2+r2+r1+r1+bl+bl+r2+r2+r1+r1+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r2+r2+r1+r1+bl+bl+r2+r2+r2+r2+r2+r2+r2+r2+r1+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+c2+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+y2+y2+y2+y1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+c1+y2+y2+y2+y1+bl+bl+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, bl+bl+bl+bl+bl+bl+bl+bl+y2+y1+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+bl+y2+y1+bl+bl+bl+bl+bl+bl+bl+bl+"\n")
	fmt.Fprintf(writer, "\n")
}
