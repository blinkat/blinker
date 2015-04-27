package storm

import (
	"fmt"
	"time"
)

var warring_color int
var error_color int
var info_color int

func init() {
	warring_color = RGBToInt(240, 173, 78)
	error_color = RGBToInt(217, 83, 79)
	info_color = RGBToInt(91, 192, 222)
}

func Error(msg string) {
	write(msg, "Error", error_color)
}

func Warring(msg string) {
	write(msg, "Warring", warring_color)
}

func Info(msg string) {
	write(msg, "Info:", info_color)
}

func write(msg string, head string, col int) {
	now := time.Now()
	fmt.Println(TextColor(col, fmt.Sprint("[", now, "]", head+":", msg)))
}
