package main

import (
	"runtime"
	"os"
	"github.com/dballard/markdown-bullet-journal/process"
	"log"
	"strconv"
	"strings"
)

type processHandler struct {
	File                  *os.File
	totalCount, doneCount int
	header                string
	headerPrinted         bool
}

func (ph *processHandler) Writeln(line string) {
	ph.File.WriteString(line + "\n")
}

func (ph *processHandler) NewFile() {
	ph.totalCount = 0
	ph.doneCount = 0
	ph.header = ""
	ph.headerPrinted = false

}

func (ph *processHandler) Eof() {
	ph.Writeln(strconv.Itoa(ph.doneCount) +  " / " + strconv.Itoa(ph.totalCount))
}

func (ph *processHandler) ProcessLine(line string, stack []string, todo bool, done bool) {
	if strings.Trim(line, " \t\n\r") == "" {
		return
	}
	if line[0] == '#' {
		ph.header = line[2:]
		ph.headerPrinted = false;
		return
	}

	if todo {
		ph.totalCount += 1
	}

	if done {
		if !ph.headerPrinted {
			ph.Writeln(" # " + ph.header)
			ph.headerPrinted = true
		}
		ph.doneCount += 1
		// TODO: Math for [x] numXnum
		ph.Writeln("  " + strings.Join(stack, " / "))
	}
}

func main() {
	ph := new(processHandler)

	if runtime.GOOS == "windows" {
		var err error
		ph.File, err = os.Open("summary.md")
		if err != nil {
			log.Fatal("Cannot open summary.md: ", err)
		}
		defer ph.File.Close()
	} else {
		ph.File = os.Stdout
	}

	files := process.GetFiles()
	for _, file := range files {
		ph.Writeln("")
		ph.Writeln(file)
		process.ProcessFile(ph, file)
	}
}
