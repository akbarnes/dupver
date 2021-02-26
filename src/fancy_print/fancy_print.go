package fancy_print

import (
	"fmt"
	"os"
	"log"
)

const ColorReset string = "\033[0m"
const ColorRed string = "\033[31m"
const ColorGreen string = "\033[32m"
const ColorYellow string = "\033[33m"
const ColorBlue string = "\033[34m"
const ColorPurple string = "\033[35m"
const ColorCyan string = "\033[36m"
const ColorWhite string = "\033[37m"

type VerbosityLevel int

const (
	CriticalLevel VerbosityLevel = iota + 1
	ErrorLevel
	WarningLevel
	NoticeLevel
	InfoLevel
	DebugLevel
)

var Verbosity = NoticeLevel
var ColorOutput = true
var logger *log.Logger

// Setup everything
func Setup(debug bool, verbose bool, quiet bool, monochrome bool) {
	InitPrinting()
	SetVerbosityLevel(debug, verbose, quiet)
	SetColoredOutput(monochrome)
}

// Initialize Go logger (used to print to stderr)
func InitPrinting() {
	logger = log.New(os.Stderr, "", 0)
}

// Set the verbosity level given command-line flags
func SetVerbosityLevel(debug bool, verbose bool, quiet bool) {
	if quiet {
		Verbosity = WarningLevel
		ColorOutput = false
	} else if debug {
		Verbosity = DebugLevel
	} else if verbose {
		Verbosity = InfoLevel
	} else {
		Verbosity = NoticeLevel
	}	
}

// Disable colored output if monochrome flag is true
func SetColoredOutput(monochrome bool) {
	if monochrome {
		ColorOutput = false
	}
}

func SetColor(color string)	{
	if ColorOutput {
		fmt.Printf("%s", color)
	}
}

func ResetColor() {
	if ColorOutput {
		fmt.Printf("%s", ColorReset)
	}
}

func Notice(msg interface{}) {
	if Verbosity >= NoticeLevel {
		logger.Println(msg)
	}
}

func Noticef(msg string, a ...interface{}) {
	if Verbosity >= InfoLevel {
		logger.Printf(msg, a...)
	}
}

func Info(msg interface{}) {
	if Verbosity >= InfoLevel {
		logger.Println(msg)
	}
}

func Infof(msg string, a ...interface{}) {
	if Verbosity >= InfoLevel {
		logger.Printf(msg, a...)
	}
}

func Debug(msg interface{}) {
	if Verbosity >= DebugLevel {
		logger.Println(msg)
	}
}

func Debugf(msg string, a ...interface{}) {
	if Verbosity >= DebugLevel {
		logger.Printf(msg, a...)
	}
}