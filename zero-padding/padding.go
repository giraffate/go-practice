package main

import "fmt"

func main() {
	fmt.Printf("zero-padding: %s\n", zeroPadding(2018, 8))
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
