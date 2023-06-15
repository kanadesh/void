package logger

import (
	"fmt"
	"strconv"
)

const ColorBlue int = 34
const ColorGreen int = 32

func Rich(colorCode int, name string, message string) {
	fmt.Println("\x1b[" + strconv.Itoa(colorCode) + "m" + name + "\x1b[0m\x1b[1m: " + message + "\x1b[0m")
}
