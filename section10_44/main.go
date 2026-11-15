package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	b := strings.Builder{}

	for n := range 10 {
		if n > 0 {
			b.WriteString(",")
		}
		
		b.WriteString(strconv.Itoa(n + 1));
		
	}

	fmt.Println(b.String())
}
