package counter

import (
	"fmt"
	"strconv"
)

func Count(state int, total int, message string) {
	var count int
	if state > total {
		count = total
	} else {
		count = state
	}

	chars := []string{
		"⠿",
		"⠾",
		"⠽",
		"⠻",
		"⠟",
		"⠯",
	}
	fmt.Print("\r\x1b[1m\x1b[34m" + chars[int(float64((count/20)%6))] + "\x1b[0m  " + message + ": \x1b[0m" + strconv.Itoa(count) + "/" + strconv.Itoa(total))
}
