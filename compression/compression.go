package compression

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	// "bytes"
)

type Dictionary struct {
	length  uint16
	entries [4096]string
}

func (dict *Dictionary) GetLength() uint16 {
	return dict.length
}

func (dict *Dictionary) Clear() {
	// Just have to reset the length
	dict.length = 256
}

func (dict *Dictionary) initialise() {
	for i := 0; i < 256; i++ {
		dict.entries[i] = string(byte(i))
	}
}

func (dict *Dictionary) AddEntry(newElement string) {
	dict.entries[dict.length] = newElement
	dict.length += 1
}

func (dict *Dictionary) GetEntry(index uint16) (string, error) {
	if index >= dict.length {
		return "", errors.New("index out of range")
	} else {
		return dict.entries[index], nil
	}
}

func getCodes(byteList []byte) [2]uint16 {
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

	// firstCode, secondCode := getCodes(buffer)
	codes := [2]uint16{uint16(total), uint16(total2)}

	return codes
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
	var dictionary Dictionary
	dictionary.initialise()

	// Need to keep track of what was emitted previously
	var lastEmitted string

	for {
		_, errMsg := file.Read(buffer)
		if errMsg != nil {
			if errMsg != io.EOF {
				fmt.Println(errMsg)
			}
			// Error or end of file, so break loop
			break
		}
		// fmt.Println(string(buffer))

		// Extract the two codes from these three bytes
		codes := getCodes(buffer)
		for i := 0; i < 2; i++ {
			// For every code, we apply the rules of the LZW decoding algorithm.
			// See https://en.wikipedia.org/wiki/Lempel–Ziv–Welch
			if codes[i] >= dictionary.GetLength() {
				// The code is not in the dictionaryß
				v := lastEmitted + lastEmitted[0:1]
				dictionary.AddEntry(v)
				fmt.Print(v)
			} else {
				// The code is in the dictionary
				w, _ := dictionary.GetEntry(codes[i])
				fmt.Print(w)
				newEntry := lastEmitted + w[0:1]
				dictionary.AddEntry(newEntry)
				lastEmitted = w
			}
		}
	}
}
