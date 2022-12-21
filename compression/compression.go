package compression

import (
	"errors"
	"fmt"
	"io"

	// "math"
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

func (dict *Dictionary) Initialise() {
	for i := 0; i < 256; i++ {
		dict.entries[i] = string(byte(i))
	}
	dict.length = 256
}

func (dict *Dictionary) AddEntry(newElement string) {
	if dict.length == 4096 {
		dict.Clear()
	}
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
		This function takes three consecutive bytes in
		byteList.  It splits the second byte in half and
		concatenates left and right such that there are now two lots
		of 12 bits, treated as numbers.  These numbers
		are returned as two integers in an array.
	*/

	firstSection := uint16(clearNibble(byteList[1], false))
	secondSection := uint16(clearNibble(byteList[1], true))

	firstSection = firstSection >> 4
	secondSection = secondSection << 8

	firstCode := uint16(byteList[0])
	firstCode = firstCode << 4
	firstCode += firstSection

	secondCode := secondSection
	secondCode += uint16(byteList[2])

	return [2]uint16{firstCode, secondCode}
}

func clearNibble(x byte, firstNibble bool) byte {
	// clears either the first or second nibble of a byte
	// firstNibble true: clears first nibble
	// firstNibble false: clears second nibble
	if firstNibble {
		x = clearBit(x, 4)
		x = clearBit(x, 5)
		x = clearBit(x, 6)
		x = clearBit(x, 7)
	} else {
		x = clearBit(x, 0)
		x = clearBit(x, 1)
		x = clearBit(x, 2)
		x = clearBit(x, 3)
	}
	return x
}

func clearBit(x byte, n byte) byte {
	// clear nth bit of given byte x
	mask := ^(byte(1) << n)
	x &= mask
	return x
}

func Decompress(encodedFile string) {
	// encodedFile:		The path of the file to be decoded

	newFilePath := encodedFile[:len(encodedFile)-2]

	// Open the encoded file
	file, errMsg := os.Open(encodedFile)
	if errMsg != nil {
		fmt.Println(errMsg)
	}
	defer file.Close()

	newFile, errMsg := os.Create(newFilePath)
	if errMsg != nil {
		fmt.Println(errMsg)
	}
	defer newFile.Close()

	// Will now read out 3 bytes at a time into the bufferArray.
	buffer := make([]byte, 3)

	// Create and fill the dictionary using previously defined struct
	var dictionary Dictionary
	dictionary.Initialise()

	// Need to keep track of what was emitted previously
	var lastEmitted string
	var endReached bool = false

	// Loop until end of file
	for {
		// Fill buffer with new bytes from file
		_, errMsg := file.Read(buffer)
		if errMsg != nil {
			if errMsg != io.EOF {
				fmt.Println(errMsg)
				break
			} else {
				endReached = true
			}
			// Error or end of file, so break loop
			break
		}

		// Deal with EOF edge cases
		var startPoint int = 0
		if endReached {
			// Either all bytes in the buffer are zero,
			// or only the last one is zero, because
			// the file must have an even number of
			// bytes.
			if buffer[0] == 0 {
				// have reached end of file
				break
			} else {
				// Third byte is zero.  Special treatment
				// required to allow the code to be
				// extracted from the first two bytes
				// properly.
				if buffer[1] == 0 && buffer[2] == 0 {
					fmt.Println("Something's broken")
				}
				buffer[2] = buffer[1]
				buffer[1] = buffer[0]
				buffer[0] = 0

				// First code contains no data, so set startPoint
				// to 1 so that it is skipped over
				startPoint = 1
			}
		}

		// Extract the two codes from these three bytes
		codes := getCodes(buffer)
		for i := startPoint; i < 2; i++ {
			// For every code, we apply the rules of the LZW decoding algorithm.
			// See https://en.wikipedia.org/wiki/Lempel–Ziv–Welch
			if codes[i] >= dictionary.GetLength() {
				// The code is not in the dictionary
				v := lastEmitted + lastEmitted[0:1]
				dictionary.AddEntry(v)

				_, errMsg = newFile.WriteString(v)
				if errMsg != nil {
					fmt.Println(errMsg)
				}

				lastEmitted = v
			} else {
				// The code is in the dictionary
				w, _ := dictionary.GetEntry(codes[i])

				_, errMsg = newFile.WriteString(w)
				if errMsg != nil {
					fmt.Println(errMsg)
				}

				if lastEmitted != "" {
					newEntry := lastEmitted + w[0:1]
					dictionary.AddEntry(newEntry)
				}

				lastEmitted = w
			}
		}

	}
}
