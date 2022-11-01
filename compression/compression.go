package compression

import (
	"fmt"
	"io"
	"math"
	"os"
	// "bytes"
)

func getCodes(byteList []byte) (uint16, uint16) {
	/*
		This function takes three consecutive bytes.  It
		splits the second byte in half and concatenates
		left and right such that there are now two lots
		of 12 bits, treated as numbers.  These numbers
		are returned as two integers in an array.

		TODO: Is there a better way using splitting
		and concatenation of the bits?
	*/

	var total uint32
	total += uint32(byteList[2])
	total += (uint32(byteList[1]) * uint32(math.Pow(2, 8)))
	total += (uint32(byteList[0]) * uint32(math.Pow(2, 16)))

	total2 := total % uint32(math.Pow(2, 12))
	total -= total2
	total /= uint32(math.Pow(2, 12))

	a := uint16(total)
	b := uint16(total2)

	return a, b
}

func Decompress(encodedFile string) {
	// encodedFile:		The name of the file to be decoded
	// plaintextFile:	The resulting file will be saved with this name

	fmt.Println("starting")

	// Open the encoded file
	file, errMsg := os.Open(encodedFile)
	if errMsg != nil {
		fmt.Println("Error opening file:")
		fmt.Println(errMsg)
	}
	defer file.Close()

	/*
		Will now read out 3 bytes at a time into the bufferArray.
		Reading in 3 bytes ensures we don't run over the end of the file
		halfway through the buffer.
	*/

	buffer := make([]byte, 3)

	// Create and fill the dictionary
	dictionary := make([]string, 256, 4096)
	for i := 0; i < 256; i++ {
		dictionary[i] = string(byte(i))
	}

	for {
		numRead, errMsg := file.Read(buffer)
		if errMsg != nil {
			if errMsg != io.EOF {
				fmt.Println(errMsg)
				fmt.Println("Number of bytes read: ", numRead)
			}
			// Error or end of file, so break loop
			break
		}
		fmt.Println(string(buffer))

		// Extract the two codes from these three bytes
		//TODO put this into the function
		firstCode, secondCode := getCodes(buffer)
		codes := [2]uint16{firstCode, secondCode}

		for i := 0; i < 2; i++ {
			if codes[i] >= uint16(len(dictionary)) {
				// The code is not in the dictionary
			} else {
				// The code is in the dictionary
			}
		}

	}

}
