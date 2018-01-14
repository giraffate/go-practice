package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	n := 20
	fmt.Print("\x1b[31m")
	for i := 0; i < n; i++ {
		fmt.Print("\r")
		fmt.Print("[")
		fmt.Printf("%s", strings.Repeat("=", i))
		if i != n-1 {
			fmt.Print(">")
			fmt.Printf("%s", strings.Repeat(" ", n-2-i))
		}
		fmt.Print("] ")
		fmt.Printf("%s/%s", zeroPadding(i, 2), zeroPadding(n, 2))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Print("\x1b[0m")
	fmt.Print("\n")
}

func zeroPadding(n, wid int) []byte {
	b := make([]byte, wid, wid)
	for wid > 0 {
		i := n / 10
		b[wid-1] = byte('0' + n - i*10)
		n = i
		wid--
	}
	return b
}
