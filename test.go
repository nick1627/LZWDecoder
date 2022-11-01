package main

import (
	c "LZWDecoder/compression"
)

func main() {
	c.Decompress("LzwInputData/01-hello.txt.z")
}

// func learning() {
// 	fmt.Println("hi")
// 	c.SayHi()

// 	// Learning stuff

// 	// const earthsGravity = 9.80665
// 	// var length uint16 = 3
// 	var b byte = 1
// 	var c byte = 255
// 	b += c
// 	fmt.Println(b)
// }
