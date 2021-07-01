package utils

import (
	"fmt"
	"os"
)

type Logger struct {
	debug bool
}

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

func (l Logger) PrintConfig(config Config) {
	if l.debug {
		fmt.Println("Config\n===============================")
		fmt.Println(config.Username)
		fmt.Println(config.Password)
		fmt.Println(config.Host)
		fmt.Println(config.Port)
		fmt.Println(config.Database)
		fmt.Println(config.Type)
		fmt.Println(config.FilePath)
		fmt.Println(config.Environment)
		fmt.Println("===============================")
	}
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
