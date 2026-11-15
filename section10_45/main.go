package main

import (
	"fmt"
	"unicode"
)

func main() {
	// k := "絵文字"
	// s := "emoji"

	// for _, c := range k {
	// 	fmt.Print(string(c) + " ")
	// }

	// fmt.Println()

	// for _, c := range s {
	// 	fmt.Print(string(c) + " ")
	// }

	// fmt.Println()

	zipEn := "123456"
	zipJp := "一二三四五六"

	for _, c := range zipEn {
		fmt.Printf("%q (IsDigit=%t)", c, unicode.IsDigit(c))
		fmt.Println()
	}

	fmt.Println()

	for _, c := range zipJp {
		fmt.Printf("%q (IsDigit=%t)", c, unicode.IsDigit(c))
		fmt.Println()
	}

	fmt.Println()
}
