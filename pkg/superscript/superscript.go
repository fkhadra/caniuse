package superscript

import (
	"fmt"
	"strconv"
	"strings"
)

var table = []string{
	"⁰",
	"¹",
	"²",
	"³",
	"⁴",
	"⁵",
	"⁶",
	"⁷",
	"⁸",
	"⁹",
}

func Itoa(input int) string {
	if input > 9 {
		var s strings.Builder
		for _, c := range strconv.Itoa(input) {
			i, _ := strconv.Atoi(fmt.Sprintf("%c", c))

			s.WriteString(table[i])
		}
		return s.String()
	}
	
	return table[input]
}
