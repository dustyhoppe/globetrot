package utils

import (
	"fmt"
	"os"
)

type Logger struct {
	debug bool
}

const GLOBETROT_HEADER = `
************************************************************************
*       _______. __        ___.           __                 __        *
*      /  _____/|  |   ____\_ |__   _____/  |________  _____/  |_      *
*     /   \  ___|  |  /  _ \| __ \_/ __ \   __\_  __ \/  _ \   __\     *
*     \    \_\  \  |_(  <_> ) \_\ \  ___/|  |  |  | \(  <_> )  |       *
*      \______  /____/\____/|___  /\___  >__|  |__|   \____/|__|       *
*             \/                \/     \/                              *
************************************************************************`

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

func (l Logger) OutputHeader() {
	fmt.Printf("%s\n\n", GLOBETROT_HEADER)
}

func (l Logger) Fatal(log string) {
	l.printLineColor(fmt.Sprintf("FATAL: %s", log), colorRed)
	os.Exit(1)
}

func (l Logger) Error(log string) {
	l.printLineColor(log, colorRed)
}

func (l Logger) Warn(log string) {
	l.printLineColor(log, colorYellow)
}

func (l Logger) Success(log string) {
	l.printLineColor(log, colorGreen)
}

func (l Logger) Details(log string) {
	l.printLineColor(log, colorCyan)
}

func (l Logger) printLineColor(log string, color string) {
	if l.debug {
		fmt.Println(string(color), log, string(colorReset))
	}
}
