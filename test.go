package main

import (
	c "LZWDecoder/compression"
)

func main() {
	// c.Decompress("LzwInputData/01-hello.txt.z")
	c.Decompress("LzwInputData/02-poem.txt.z")
	// c.Decompress("LzwInputData/03-book.txt.z")
	// c.Decompress("LzwInputData/04-icon.txt.z")
}
